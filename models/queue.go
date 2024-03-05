package models

type QueueModel struct {
	CurrentQueueNumber int64          `firebase:"CurrentQueueNumber"`
	QueueBucket        int64          `firebase:"QueueBucket"`
	QueueList          map[string]int `firebase:"QueueList"`
}

type CustomerQueue struct {
	CustomerQueueNumber int64
	QueueBucket         int64
}
