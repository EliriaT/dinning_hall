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

## Running the App
To run the App, run in terminal the following command:<br />


`go run .`


## Running in Docker container
1. To run the app in a docker container, first build the image:<br />

`docker build -t dinning-hall .`

2. Then run the container using the created image:<br />

`docker run --name dinning-hall --network restaurant -it --rm  -p 8082:8082 dinning-hall`

For this you firstly need a created docker network. To create a docker network run:

`docker network create restaurant`

3. To stop the running container:

`docker stop {docker's id}`

## Combining with kitchen server

The kitchen server listens first for Post request coming from the dining-hall. To run the system correctly, the kitchen server must run first, and after it the dinning-hall
server should start running. These servers use HTTP Post request for communication.
