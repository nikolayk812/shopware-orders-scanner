package shopware

import "time"

const MaxSearchLimit = 500

type tokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` //in seconds
}

type filter struct {
	Type       filterType        `json:"type"`
	Field      string            `json:"field"`
	Value      interface{}       `json:"value"`
	Parameters *filterParameters `json:"parameters,omitempty"`
}

type filterParameters struct {
	GTE string `json:"gte"`
	LTE string `json:"lte"`
}

// FilterType as described in https://docs.shopware.com/en/shopware-platform-dev-en/api/filter-search-limit#filter-1
type filterType string

const (
	filterTypeEquals    filterType = "equals"
	filterTypeEqualsAny filterType = "equalsAny"
	filterTypeRange     filterType = "range"
)

type Order struct {
	ID                string             `json:"id"`
	SalesChannelID    string             `json:"salesChannelId"`
	AutoIncrement     int                `json:"autoIncrement"`
	Number            string             `json:"orderNumber"`
	Deliveries        []OrderDelivery    `json:"deliveries"`
	Transactions      []OrderTransaction `json:"transactions"`
	Documents         []OrderDocument    `json:"documents"`
	LineItems         []LineItem         `json:"lineItems"`
	StateMachineState struct {
		Name OrderState `json:"name"`
	} `json:"stateMachineState"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type OrderDelivery struct {
	OrderID           string   `json:"orderId"`
	TrackingCodes     []string `json:"trackingCodes"`
	UpdatedAt         string   `json:"updatedAt"`
	StateMachineState struct {
		Name OrderDeliveryState `json:"name"`
	} `json:"stateMachineState"`
}

type OrderTransaction struct {
	OrderID           string `json:"orderId"`
	StateMachineState struct {
		Name OrderTransactionState `json:"name"`
	} `json:"stateMachineState"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderDocument struct {
	FileType string `json:"fileType"`
	Config   struct {
		Custom struct {
			FileName string `json:"fileName"`
		} `json:"custom"`
	} `json:"config"`
}

type LineItem struct {
	ProductID string `json:"productId"`
	Payload   struct {
		ProductNumber string `json:"productNumber"`
	} `json:"payload"`
}

type OrderState string

const (
	OrderStateOpen       OrderState = "Open"
	OrderStateInProgress OrderState = "In progress"
	OrderStateDone       OrderState = "Done"
	OrderStateCancelled  OrderState = "Cancelled"
)

type OrderTransactionState string

const (
	OrderTransactionStateOpen              OrderTransactionState = "Open"
	OrderTransactionStatePaid              OrderTransactionState = "Paid"
	OrderTransactionStateCancelled         OrderTransactionState = "Cancelled"
	OrderTransactionStateRefunded          OrderTransactionState = "Refunded"
	OrderTransactionStateRefundedPartially OrderTransactionState = "Refunded (partially)"
	OrderTransactionStateReminded          OrderTransactionState = "Reminded"
	OrderTransactionStateFailed            OrderTransactionState = "Failed"
	OrderTransactionStateInProgress        OrderTransactionState = "In Progress"
)

type OrderDeliveryState string

const (
	OrderDeliveryStateOpen              OrderDeliveryState = "Open"
	OrderDeliveryStateShipped           OrderDeliveryState = "Shipped"
	OrderDeliveryStateShippedPartially  OrderDeliveryState = "Shipped (partially)"
	OrderDeliveryStateCancelled         OrderDeliveryState = "Cancelled"
	OrderDeliveryStateReturned          OrderDeliveryState = "Returned"
	OrderDeliveryStateReturnedPartially OrderDeliveryState = "Returned (partially)"
)

type Product struct {
	Stock int `json:"stock"`
}
