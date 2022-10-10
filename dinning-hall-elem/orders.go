package dinning_hall_elem

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

var AiOrder autoInc

// used to autoincrement order's id
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

	nrItems := rand.Intn(MaxFoods) + minFoods

	for i := 0; i < 3; i++ {
		if nrItems > 5 {
			nrItems = rand.Intn(MaxFoods) + minFoods
		} else {
			break
		}
	}

	foodList := make([]int, nrItems)
	maxWait := 0

	for i := 0; i < nrItems; i++ {
		foodId := rand.Intn(NrFoods) + 1
		//in the json the id starts from 1
		foodList[i] = foodId

		if prepTime := Foods[foodId-1].PreparationTime; prepTime > maxWait {
			maxWait = prepTime
		}

	}
	//sort in descending order according foods complexity
	sort.Slice(foodList, func(i, j int) bool {
		return Foods[foodList[i]-1].Complexity > Foods[foodList[j]-1].Complexity

	})

	return Order{
		Id:    AiOrder.ID(),
		Items: foodList,
		//5 the lowest priority, 1 the most priority.
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
