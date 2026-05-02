package address

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressModule module.AddressModule
}

func NewAddressHandler(aModule module.AddressModule) *AddressHandler {
	return &AddressHandler{addressModule: aModule}
}

// CreateAddress creates a new shipping address for the user
// @Summary Create address
// @Description Add a new shipping address for the authenticated user
// @Tags Address
// @Accept json
// @Produce json
// @Param request body dto.CreateAddressRequest true "Address Data"
// @Success 201 {object} response.BaseResponse{data=dto.Address}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/addresses [post]
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Errorf("invalid request body: %v", err))
		return
	}

	address, err := h.addressModule.CreateAddress(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, address)
}

// GetAddresses retrieves all addresses for the authenticated user
// @Summary List addresses
// @Description Get a list of all shipping addresses saved by the user
// @Tags Address
// @Produce json
// @Success 200 {object} response.BaseResponse{data=[]dto.Address}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/addresses [get]
func (h *AddressHandler) GetAddresses(c *gin.Context) {
	userID := c.GetInt64("user_id")

	addresses, err := h.addressModule.GetAddressesByUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, addresses)
}

// UpdateAddress updates an existing address
// @Summary Update address
// @Description Update an existing shipping address
// @Tags Address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param request body dto.UpdateAddressRequest true "Address Data"
// @Success 200 {object} response.BaseResponse{data=dto.Address}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/addresses/:id [put]
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, fmt.Errorf("invalid address id"))
		return
	}

	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Errorf("invalid request body: %v", err))
		return
	}

	address, err := h.addressModule.UpdateAddress(c.Request.Context(), userID, id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, address)
}

// DeleteAddress deletes an address
// @Summary Delete address
// @Description Remove a shipping address
// @Tags Address
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/addresses/:id [delete]
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, fmt.Errorf("invalid address id"))
		return
	}

	if err := h.addressModule.DeleteAddress(c.Request.Context(), userID, id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "address deleted successfully"})
}

// SetDefaultAddress sets an address as the default for the user
// @Summary Set default address
// @Description Set a specific address as the default shipping address
// @Tags Address
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/addresses/:id/default [put]
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, fmt.Errorf("invalid address id"))
		return
	}

	if err := h.addressModule.SetDefaultAddress(c.Request.Context(), userID, id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "default address updated"})
}
