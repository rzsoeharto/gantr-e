package models

type QueueModel struct {
	RestaurantName     string         `firebase:"RestaurantName"`
	CurrentQueueNumber int64          `firebase:"CurrentQueueNumber"`
	QueueBucket        int64          `firebase:"QueueBucket"`
	QueueList          map[string]int `firebase:"QueueList"`
}

type CustomerQueue struct {
	CustomerQueueNumber int64
	QueueBucket         int64
}
