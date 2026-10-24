package main

import "fmt"

// Only rocket payment has discount, How to apply ??

type Paypal struct {
	accountId string
}

func (p Paypal) Pay(amount int) string {
	return fmt.Sprintf("%d amount Paypal is successful for %s", amount, p.accountId)
}

func (p Paypal) CashOut(amount int) {
	// Nothig is here
}

type Bkash struct {
	number string
}

func (b Bkash) Pay(amount int) string {
	return fmt.Sprintf("%d amount Bkash is successful for %s", amount, b.number)
}

func (b Bkash) CashOut(amount int) {
	fmt.Println("Cashout from", b.number)
}

type Rocket struct {
	email string
}

func (r Rocket) Pay(amount int) string {
	return fmt.Sprintf("%d amount Rocket is successful for %s", amount, r.email)
}

func (r Rocket) CashOut(amount int) {
	fmt.Println("Cashout from", r.email)
}

// Paypal doesn't have CashOut method, here is an error
// How to fix this ??
type PaymentMethod interface {
	Pay(int) string
	CashOut(int)
}

func Payment(p PaymentMethod, amount int) {
	if _, ok := p.(Rocket); ok {
		amount = discountCalculate(amount, 10)
	}
	fmt.Println(p.Pay(amount))
}


func discountCalculate(amount int, discount int) int {
	return amount - (amount * discount / 100)
} 

func PaymentMangement() {
	
	p := Paypal{accountId: "labib"}
	Payment(p, 123)

	b := Bkash{number: "qw-34-qfe"}
	Payment(b, 10023)

	r := Rocket{email: "labibfaisal@gmail.com"}
	Payment(r, 1000)
}