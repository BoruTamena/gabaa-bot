package address

import (
	"context"
	"fmt"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type addressModule struct {
	addressStorage storage.AddressStorage
}

func NewAddressModule(storage storage.AddressStorage) module.AddressModule {
	return &addressModule{
		addressStorage: storage,
	}
}

func (m *addressModule) CreateAddress(ctx context.Context, userID int64, req dto.CreateAddressRequest) (*dto.Address, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.IsDefault {
		_ = m.addressStorage.ClearDefaultAddress(ctx, userID)
	}

	label := req.Label
	if label == "" {
		label = "home"
	}

	address := &db.Address{
		UserID:        userID,
		Label:         label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		Street:        req.Street,
		City:          req.City,
		Region:        req.Region,
		Country:       req.Country,
		IsDefault:     req.IsDefault,
	}

	if err := m.addressStorage.CreateAddress(ctx, address); err != nil {
		return nil, err
	}

	// If it's the first address, make it default
	if !req.IsDefault {
		addresses, _ := m.addressStorage.GetAddressesByUserID(ctx, userID)
		if len(addresses) == 1 {
			address.IsDefault = true
			_ = m.addressStorage.UpdateAddress(ctx, address)
		}
	}

	return m.mapToDTO(address), nil
}

func (m *addressModule) GetAddress(ctx context.Context, id int64) (*dto.Address, error) {
	address, err := m.addressStorage.GetAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.mapToDTO(address), nil
}

func (m *addressModule) GetAddressesByUser(ctx context.Context, userID int64) ([]dto.Address, error) {
	addresses, err := m.addressStorage.GetAddressesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var dtos []dto.Address
	for _, a := range addresses {
		dtos = append(dtos, *m.mapToDTO(&a))
	}
	return dtos, nil
}

func (m *addressModule) UpdateAddress(ctx context.Context, userID int64, id int64, req dto.UpdateAddressRequest) (*dto.Address, error) {
	address, err := m.addressStorage.GetAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if address.UserID != userID {
		return nil, fmt.Errorf("unauthorized to update this address")
	}

	if req.IsDefault && !address.IsDefault {
		_ = m.addressStorage.ClearDefaultAddress(ctx, userID)
	}

	if req.Label != "" {
		address.Label = req.Label
	}
	if req.RecipientName != "" {
		address.RecipientName = req.RecipientName
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.Street != "" {
		address.Street = req.Street
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.Region != "" {
		address.Region = req.Region
	}
	if req.Country != "" {
		address.Country = req.Country
	}
	address.IsDefault = req.IsDefault

	if err := m.addressStorage.UpdateAddress(ctx, address); err != nil {
		return nil, err
	}

	return m.mapToDTO(address), nil
}

func (m *addressModule) DeleteAddress(ctx context.Context, userID int64, id int64) error {
	address, err := m.addressStorage.GetAddressByID(ctx, id)
	if err != nil {
		return err
	}

	if address.UserID != userID {
		return fmt.Errorf("unauthorized to delete this address")
	}

	return m.addressStorage.DeleteAddress(ctx, id)
}

func (m *addressModule) SetDefaultAddress(ctx context.Context, userID int64, id int64) error {
	address, err := m.addressStorage.GetAddressByID(ctx, id)
	if err != nil {
		return err
	}

	if address.UserID != userID {
		return fmt.Errorf("unauthorized to update this address")
	}

	_ = m.addressStorage.ClearDefaultAddress(ctx, userID)

	address.IsDefault = true
	return m.addressStorage.UpdateAddress(ctx, address)
}

func (m *addressModule) mapToDTO(a *db.Address) *dto.Address {
	return &dto.Address{
		ID:            a.ID,
		UserID:        a.UserID,
		Label:         a.Label,
		RecipientName: a.RecipientName,
		Phone:         a.Phone,
		Street:        a.Street,
		City:          a.City,
		Region:        a.Region,
		Country:       a.Country,
		IsDefault:     a.IsDefault,
		CreatedAt:     a.CreatedAt,
	}
}
