package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

//func LoadMetricFromFile(fstg string) {
//	var res []byte
//
//	results, err := errorRetriable.ErrorRetriable(os.ReadFile, fstg)
//	if err != nil {
//		logger.Log.Infof("Ошибка вызова функции для повторного вызова функции: %s", err)
//	}
//	for _, result := range results {
//		switch v := result.(type) {
//		case []byte:
//			res = v
//		case error:
//			err = v
//		}
//	}
//
//	if err != nil {
//		logger.Log.Infof("Ошибка чтения файла %s: %s", fstg, err)
//	}
//
//	if err := json.Unmarshal(res, &StorageMetric); err != nil {
//		logger.Log.Infof("Ошибка Unmarshal: %s", err)
//	}
//}

func LoadMetricFromFile(fstg string) error {

	res, err := os.ReadFile(fstg)
	if err != nil {
		return fmt.Errorf("Ошибка чтения файла: %s", err)
	}

	if err := json.Unmarshal(res, &StorageMetric); err != nil {
		return fmt.Errorf("Ошибка Unmarshal: %s", err)
	}
	return nil
}
