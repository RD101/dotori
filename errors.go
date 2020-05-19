package main

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// IsDup 함수는 인자로 받은 에러가 mongoDB의 duplicate key error인지 판단한다.
func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 { // 11000은 errorDuplicateKey에러를 의미한다. https://github.com/mongodb/mongo-go-driver/blob/c814cfb676dc29a6ad6171e138510ee88c51983f/mongo/integration/collection_test.go#L27
				return true
			}
		}
	}
	return false
}
