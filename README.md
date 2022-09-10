# dining_hall

## About
This project simulates a dining hall of a restaurant. The Dining hall has a finite amount of tables that "clients" can occupy. Main work unit of the Dining hall are waiters. `Tables` (clients) generates orders based on restaurant menu. Menu consist of `foods` . A `food` is described by the following fields:

```golang
{
"id": 1,
"name": "pizza", 
"preparation-time": 20,
"complexity": 2 ,
"cooking-apparatus": "oven",
}
```

Tables generates orders and an order should is described by the following information:

```golang
{
"id": 1,
"items": [ 3, 4, 4, 2 ],
"priority": 3 ,
"max_wait": 45
}
```

`Waiters` are  an object instances which run their logic of serving tables on separate threads , one thread per waiter.
For `Waiters` which are running on separate threads , tables are shared resource. When waiter is picking up the order from a table , the table generates a random order with random foods and random number of foods, random `priority` and unique order `ID`. `Waiter` have to send order to
`kitchen` by performing HTTP (POST) request, with order details.
When order will be ready, `kitchen` will send a `HTTP` (POST) request back to `Dinning Hall` . The Dinning Hall server has
to handle that request and to notify waiter that order is ready to be served to the table which requested this order.
