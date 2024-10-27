package errorRetriable

import (
	"errors"
	"github.com/LI-SeNyA-vE/KursMetrics/internal/middleware/logger"
	"reflect"
	"time"
)

var timeDelay = [3]time.Duration{1, 3, 5}

//	func ErrorRetriable(inputFunc func()) (any, error) {
//		for _, i := range timeDelay {
//			result, err := inputFunc
//			if err == nil {
//				return result, nil
//			}
//			time.Sleep(i * time.Second)
//		}
//		return nil, err
//	}
func ErrorRetriable(inputFunc interface{}, args ...interface{}) ([]interface{}, error) {
	// Проверяем, что inputFunc — это функция
	funcValue := reflect.ValueOf(inputFunc)
	if funcValue.Kind() != reflect.Func {
		return nil, errors.New("передана была не функция")
	}

	// Преобразуем аргументы в reflect.Value
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// Пробуем вызвать функцию с перезапусками на случай ошибки
	var result []reflect.Value
	var err error
	for _, delay := range timeDelay {
		// Выполняем вызов функции и ловим возможные ошибки
		result = funcValue.Call(in)
		logger.Log.Infof("Выполнена функция %s", inputFunc)
		// Проверяем последний результат на наличие ошибки
		if len(result) > 0 {
			lastResult, ok := result[len(result)-1].Interface().(error)
			if ok && lastResult != nil {
				err = lastResult
				time.Sleep(delay * time.Second) // Ждём перед следующим повтором
				continue
			} else if lastResult == nil {
				err = nil //Если функция выполнилась без ошибок
			} else {
				err = errors.New("пока неизвестно что будет, но скорее всего последний элемент не типа errors")
			}
		}
		break
	}

	// Если осталась ошибка после всех попыток, возвращаем её
	if err != nil {
		return nil, err
	}

	// Формируем выходные данные
	output := make([]interface{}, len(result))
	for i, res := range result {
		output[i] = res.Interface()
	}

	return output, nil
}

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
