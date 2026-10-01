package main

import "fmt"

type MakePaymentRequest struct {
	PaymentAddress   string
	AppID            string
	IdempotencyToken string
}

type PaymentProcessor interface {
	MakePayment(MakePaymentRequest) error
}

type Razorpay struct {
	//Internal Deps
}

type MakeUPIPaymentRequest struct {
	PaymentAddress   string
	AppID            string
	IdempotencyToken string
}

func (r *Razorpay) MakeUPIPayment(req MakeUPIPaymentRequest) error {
	fmt.Println("making UPI Payment", req.PaymentAddress, req.AppID, req.IdempotencyToken)
	return nil
}

// Problem is that we can't use Razorpay for interface PaymentProcessor.
// Due to difference in method.
// Hence we have a adapter.

type RazorpayAdapter struct {
	Razorpay Razorpay
}

func (r *RazorpayAdapter) MakePayment(req MakePaymentRequest) error {
	return r.Razorpay.MakeUPIPayment(MakeUPIPaymentRequest{
		PaymentAddress:   req.PaymentAddress,
		AppID:            req.AppID,
		IdempotencyToken: req.IdempotencyToken,
	})
}

func main() {
	var processor PaymentProcessor = &RazorpayAdapter{
		Razorpay: Razorpay{},
	}

	err := processor.MakePayment(MakePaymentRequest{
		PaymentAddress:   "test@upi",
		AppID:            "lld-app",
		IdempotencyToken: "tok-123",
	})
	if err != nil {
		fmt.Println("payment failed:", err)
		return
	}

	fmt.Println("payment done via PaymentProcessor (adapter)")
}
