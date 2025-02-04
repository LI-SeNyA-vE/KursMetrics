/*
Package errorretriable предоставляет механизм повторных попыток (retry)
выполнения некоторой функции, пока она не вернёт nil в качестве ошибки
или не исчерпается количество попыток.
*/
package errorretriable

import (
	"time"
)

// timeDelay определяет интервалы (в секундах) между повторами.
// Здесь всего 3 повторные попытки, с задержкой в 1 секунду каждую.
var timeDelay = [3]time.Duration{1, 1, 1}

// ErrorRetriableHTTP принимает функцию inputFunc, которая возвращает (interface{}, error).
// Если при её выполнении возникает ошибка, функция повторно вызывается через
// указанные интервалы (см. timeDelay). По истечении всех попыток возвращается
// последняя ошибка, если она не была устранена.
//
// Возвращает:
//   - interface{} — результат успешного вызова inputFunc,
//   - error — если не удалось выполнить без ошибки за заданное число повторов.
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
