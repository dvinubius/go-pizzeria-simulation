# Pizzeria & Consumers

A simulation in Go to demonstrate channels.
 
>- A pizzeria that takes orders and responds with 3 possible outcomes
>  - reject if closed
>  - successfully make pizza
>  - fail at making pizza 
>- Pizzeria starts as open but closes after a while
>- Making a pizza fails based on random chance
>- Customers react to responses
>- Customers make limited number of orders with random delays. 
>- When all orders have been placed and responses received, simulation can conclude
>- As a conclusion, show a review of the pizzeria's business for the day

Run with

`go run main.go customer.go pizzamaker.go`

![Screenshot 2022-11-30 at 18 38 07](https://user-images.githubusercontent.com/32189942/204856585-26a2155c-cb21-4e2b-ab87-ad48fc08ecd1.png)

# Design considerations

Pizzeria, each customer and the main thread execute concurrently and communicate via go channels.

## Channels

```
--some_random_orders--> Consumer order OUT 
  -------collect--------> Pizzeria order IN
    ---------process------> Pizzeria order OUT
      --------dispatch------> Consumer order IN
```


## Main
1. create pizzeria and start business
2. create customers and start consuming (order & react to response)
3. wait for all customers to quit
4. business review


## Simulation Control
Each consumer orders a fixed number of times, then closes its order OUT channel

-> All Consumers close orders OUT 
  -> Main closes Pizzeria orders IN 
    -> Pizzeria finishes ranging over its orders IN 
      -> Pizzeria closes its orders OUT
        -> Main Dispatcher finishes ranging over Pizzeria orders OUT
          -> Main Dispatcher closes all customers' orders IN
            -> All Customers finish ranging over their orders IN
              -> All Customers quit
                -> Main proceeds to Business Review


For demo purposes, the params are chosen such that pizzeria closes before customers have placed all their orders => last orders will surely be rejected.

Tweak the params:
- customer delay times
- total number of order attepts per customer
- pizzeria's open time 
