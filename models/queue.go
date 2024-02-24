package models

type QueueModel struct {
	CurrentQueueNumber int64          `firebase:"CurrentQueueNumber"`
	QueueBucket        int64          `firebase:"QueueBucket"`
	QueueList          map[string]int `firebase:"QueueList"`
}
