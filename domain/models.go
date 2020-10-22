package domain

type OrderResult struct {
	OrderID      string
	OrderNumber  string
	ChannelID    string
	TrackingCode string
	Errors       map[string]error
}
