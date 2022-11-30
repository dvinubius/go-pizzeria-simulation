package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// supposing there is only 1 type of pizza
type PizzaOrder struct {
	customerNumber int
	attemptNumber  int
	orderNumber    int
	message        string
	success        bool
	rejected       bool
}

const numberCustomers = 5

var activeCustomers = 0

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// create a producer and run in the background
	pizzaMaker := createMaker()
	go pizzaMaker.doBusiness()

	// create customers and run them in the background
	customersEmo := strings.Repeat(" 🧑 ", numberCustomers)
	fmt.Printf("\n| - | - | - CUSTOMERS%v - | - | - |\n\n", customersEmo)
	customers := make([]*Customer, numberCustomers)
	var wgCust sync.WaitGroup
	wgCust.Add(numberCustomers)
	for i := 0; i < numberCustomers; i++ {
		activeCustomers++
		customers[i] = createCustomer(i + 1)
		go customers[i].consume()

		// collect orders
		go collectOrders(pizzaMaker, customers[i].ordersOut)
		go listenToQuit(pizzaMaker, customers[i].quit, &wgCust)
	}

	// handle order responses from pizza maker
	// (rejected / accepted but failed / accepted & fulfilled)
	go dispatchResponses(pizzaMaker, customers)

	// wait for all customers to quit
	wgCust.Wait()

	Review(pizzaMaker)
}

// customers -> pizzeria
func collectOrders(p *PizzaMaker, custOrdersOut chan *PizzaOrder) {
	for order := range custOrdersOut {
		p.ordersIn <- order
	}
	// customer has made all their orders
	activeCustomers--
	if activeCustomers == 0 {
		close(p.ordersIn)
		fmt.Println("\n= = = Customers are done ordering = = =\n")
	}
}

// pizzeria -> customers
func dispatchResponses(p *PizzaMaker, customers [](*Customer)) {
	for order := range p.ordersOut {
		for _, c := range customers {
			if order.customerNumber == c.number {
				orderPart := ""
				if !order.rejected {
					orderPart = fmt.Sprintf(" (order #%v)", order.orderNumber)
				}
				fmt.Printf(
					"    -->> Response to attempt #%v%v: %v\n",
					order.attemptNumber,
					orderPart,
					color.CyanString(order.message),
				)
				c.ordersIn <- order
			}
		}
	}
	// once no more responses can come from producer, close customer IN channels
	time.Sleep(time.Second)
	for _, c := range customers {
		close(c.ordersIn)
	}
}

func listenToQuit(p *PizzaMaker, quitChan chan bool, wg *sync.WaitGroup) {
	<-quitChan
	wg.Done()
}

func Review(p *PizzaMaker) {
	color.Cyan(
		"* * * PIZZERIA REVIEW * * *\n\nMade %d pizzas, failed to make %d and rejected %d orders.\n\n",
		p.pizzasMade,
		p.pizzasFailed,
		p.ordersRejected,
	)
}
