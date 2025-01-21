package errorretriable

import (
	"time"
)

var timeDelay = [3]time.Duration{1, 3, 5}

// Функция для повторного выполнения
func ErrorRetriableHTTP(inputFunc func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error

	for _, delay := range timeDelay {
		result, err = inputFunc()
		if err == nil {
			return result, nil
		}
		time.Sleep(delay * time.Second)
	}

	return nil, err
}
