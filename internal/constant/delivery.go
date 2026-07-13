package constant

type DeliveryAgentStatus string

const (
	DeliveryAgentStatusPendingInvite DeliveryAgentStatus = "pending_invite"
	DeliveryAgentStatusActive        DeliveryAgentStatus = "active"
	DeliveryAgentStatusSuspended     DeliveryAgentStatus = "suspended"
)

func (s DeliveryAgentStatus) String() string {
	return string(s)
}

type DeliveryLocationType string

const (
	DeliveryLocationTypePickup   DeliveryLocationType = "pickup"
	DeliveryLocationTypeDelivery DeliveryLocationType = "delivery"
)

func (t DeliveryLocationType) String() string {
	return string(t)
}

const RoleDelivery = "delivery"
