package main

import "fmt"

type Payment interface {
	Pay() error
}

type cardPayment struct {
	number string
	cvv    string
}

func NewCardPayment() Payment {
	return &cardPayment{
		number: "12345",
		cvv:    "67890",
	}
}

func (p *cardPayment) Pay() error {
	fmt.Println("Оплата картой")
	return nil
}

type qiwiPayment struct {
	number string
}

func NewQIWIPayment() Payment {
	return &qiwiPayment{
		number: "12345",
	}
}

func (p *qiwiPayment) Pay() error {
	fmt.Println("Оплата через QIWI")
	return nil
}

func processOrder(product string, payment Payment) {
	if err := payment.Pay(); err != nil {
		return
	}
}

func main() {
	product := "phone"
	payWay := 2

	var payment Payment
	switch payWay {
	case 1:
		payment = NewCardPayment()
	case 2:
		payment = NewQIWIPayment()
	}

	processOrder(product, payment)
}
