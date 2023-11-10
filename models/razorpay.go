package models


// RazorPay Model
type RazorPay struct {
	UserID          uint    `JSON:"userid"`
	RazorPaymentID  string  `JSON:"razorpaymentid" gorm:"primaryKey"`
	RazorPayOrderID string  `JSON:"razorpayorderid"`
	Signature       string  `JSON:"signature"`
	AmountPaid      float64 `JSON:"amountpaid"`
}