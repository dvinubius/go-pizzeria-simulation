package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type Customer struct {
	number    int
	ordersOut chan *PizzaOrder
	ordersIn  chan *PizzaOrder
	quit      chan bool
}

func createCustomer(num int) *Customer {
	makeCh := make(chan *PizzaOrder)
	respCh := make(chan *PizzaOrder)
	return &Customer{
		number:    num,
		ordersOut: makeCh,
		ordersIn:  respCh,
		quit:      make(chan bool),
	}
}

const orderAttempts = 3

// customer orders from time to time and processes the responses (delivered / failed / rejected)
func (c *Customer) consume() {
	go c.processResponses()

	// order a few pizzas with delays in between
	rand.Seed(time.Now().UnixNano())
	start := rand.Intn(5000)
	interval := rand.Intn(3000)

	time.Sleep(time.Duration(start) * time.Millisecond)
	for i := 0; i < orderAttempts; i++ {
		c.ordersOut <- &PizzaOrder{
			customerNumber: c.number,
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
	close(c.ordersOut)
}

func (c *Customer) processResponses() {
	for o := range c.ordersIn {
		if o.rejected {
			color.Yellow(fmt.Sprintf("      -->> CUST %v: Why are you closed?! 😭\n", c.number))
		} else if o.success {
			color.Green(fmt.Sprintf("      -->> CUST %v: Received order #%d 😋!", c.number, o.orderNumber))
		} else {
			color.Red(fmt.Sprintf("      -->> CUST %v: I wanted a pizza, what is this? 😡", c.number))
		}
	}
	// once orders IN is closed, customer can quit
	c.quit <- true
	close(c.quit)
}
