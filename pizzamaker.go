package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const MaxPizzasConc = 5
const MaxPizzasTotal = 20

type PizzaMaker struct {
	open           bool
	attemptsCount  int // attempts to order
	orderCount     int // accepted orders
	ordersOut      chan *PizzaOrder
	ordersIn       chan *PizzaOrder
	pizzasMade     int
	pizzasFailed   int
	ordersRejected int
}

func createMaker() *PizzaMaker {
	return &PizzaMaker{
		ordersOut: make(chan *PizzaOrder, MaxPizzasConc),
		ordersIn:  make(chan *PizzaOrder, MaxPizzasConc),
	}
}

func (p *PizzaMaker) doBusiness() {
	go p.handleOrders()

	p.open = true
	// print out a message
	color.Cyan("= = = Pizzeria is open for business! = = =\n\n")

	time.Sleep(15 * time.Second)
	p.open = false
	color.Cyan("\n= = = Pizzeria takes no more orders = = =\n\n")
}

func (p *PizzaMaker) handleOrders() {
	for o := range p.ordersIn {
		p.attemptsCount++
		o.attemptNumber = p.attemptsCount
		if !p.open {
			p.rejectOrder(o)
			p.ordersOut <- o
		} else {
			p.processOrder(o)
			p.ordersOut <- o // successful or failed
		}
	}
	close(p.ordersOut)
}

func (p *PizzaMaker) rejectOrder(order *PizzaOrder) {
	conf := fmt.Sprintf(
		"On attempt #%v REJECT order from CUST %v!\n",
		order.attemptNumber,
		order.customerNumber,
	)
	color.Cyan(conf)
	order.rejected = true
	order.message = "⛔️ We're closed, sorry."
	p.ordersRejected++
}

func (p *PizzaMaker) processOrder(order *PizzaOrder) {
	p.orderCount++
	order.orderNumber = p.orderCount

	delay := rand.Intn(5) + 1
	conf := fmt.Sprintf(
		"On attempt #%v ACCEPT order #%d from CUST %v!\n",
		order.attemptNumber,
		order.orderNumber,
		order.customerNumber,
	)
	color.Cyan(conf)

	rnd := rand.Intn(12) + 1
	message := ""
	success := false

	if rnd < 5 {
		p.pizzasFailed++
	} else {
		p.pizzasMade++
	}

	fmt.Printf(" -->> Making pizza #%d. It will take %d seconds....\n", p.orderCount, delay)
	// delay for a bit
	time.Sleep(time.Duration(delay) * time.Second)

	if rnd <= 2 {
		message = fmt.Sprintf("❌ We ran out of ingredients for pizza #%d!", p.orderCount)
	} else if rnd <= 4 {
		message = fmt.Sprintf("😣 The cook quit while making pizza #%d!", p.orderCount)
	} else {
		success = true
		message = fmt.Sprintf("✅ Enjoy pizza number #%d!", p.orderCount)
	}

	order.message = message
	order.success = success
}
