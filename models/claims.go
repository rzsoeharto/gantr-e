package models

import "github.com/golang-jwt/jwt/v5"

type QueueClaimsStruct struct {
	QueueNumber int64 `json:"qno"`
	// QueueBucket int64 `json:"bkt"`
	jwt.RegisteredClaims
}

func (t QueueClaimsStruct) GetQueueNumber() int64 {
	return t.QueueNumber
}
