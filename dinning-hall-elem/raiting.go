package dinning_hall_elem

import (
	"time"
)

func CalculateAverage(rating int) float64 {
	sum = sum + rating
	//for _, mark := range OrderMarks {
	//	sum += mark
	//}
	avg := float64(sum) / float64(MarkLength)
	return avg
}

func giveOrderStars(serveTime time.Duration, maxWait float64) int {
	//serveTimeMillisec := float64(serveTime)*1000 //time in milliseconds
	serveTimeNonUnit := float64(serveTime) / float64(TimeUnit)

	switch {
	case serveTimeNonUnit < maxWait:
		return 5
	case serveTimeNonUnit < maxWait*1.1:
		return 4
	case serveTimeNonUnit < maxWait*1.2:
		return 3
	case serveTimeNonUnit < maxWait*1.3:
		return 2
	case serveTimeNonUnit < maxWait*1.4:
		return 1
	default:
		return 0
	}
}
