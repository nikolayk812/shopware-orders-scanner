package domain

type OrderResult struct {
	OrderID      string
	OrderNumber  string
	ChannelID    string
	TrackingCode string
	CreatedDate  string
	Errors       map[string]error
}
