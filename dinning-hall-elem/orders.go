package dinning_hall_elem

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	maxFoods = 10
	minFoods = 1
	nrFoods  = 13
)

var OrdersChannel = make(chan int, len(Tables))
var AiOrder autoInc

// used to autoincrement order's id, tables will run go routines when generating order,
// that is why it should be locked
type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int
}

func (a *autoInc) SetId(id int) {
	a.id = id
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()
	id = a.id
	a.id++
	return
}

type Order struct {
	Id       int     `json:"id"`
	Items    []int   `json:"items"`
	Priority int     `json:"priority"`
	MaxWait  float64 `json:"max_wait"`
}

// Function to create new order and return it
func newOrder() Order {

	nrItems := rand.Intn(maxFoods) + minFoods
	foodList := make([]int, nrItems)
	maxWait := 0

	for i := 0; i < nrItems; i++ {
		foodId := rand.Intn(nrFoods) + 1
		//in the json the id starts from 1
		foodList[i] = foodId

		if prepTime := Foods[foodId-1].PreparationTime; prepTime > maxWait {
			maxWait = prepTime
		}

	}
	return Order{
		Id:       AiOrder.ID(),
		Items:    foodList,
		Priority: int(math.Round(float64(nrItems) / 2)),
		MaxWait:  float64(maxWait) * 1.3,
	}
}

// order data received from kitchen
type ReceivedOrd struct {
	OrderId        int              `json:"order_id"`
	TableId        int              `json:"table_id"`
	WaiterId       int              `json:"waiter_id"`
	Items          []int            `json:"items"`
	Priority       int              `json:"priority"`
	MaxWait        float64          `json:"max_wait"`
	PickUpTime     time.Time        `json:"pick_up_time"`
	CookingTime    time.Duration    `json:"cooking_time"`
	CookingDetails []kitchenFoodInf `json:"cooking_details"`
}

// order data sent to kitchen
type sentOrd struct {
	OrderId    int       `json:"order_id"`
	TableId    int       `json:"table_id"`
	WaiterId   int       `json:"waiter_id"`
	Items      []int     `json:"items"`
	Priority   int       `json:"priority"`
	MaxWait    float64   `json:"max_wait"`
	PickUpTime time.Time `json:"pick_up_time"`
}
