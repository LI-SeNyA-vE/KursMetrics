package filemetric

import (
	"encoding/json"
	"fmt"
	"os"
)

func (s *FileStorage) LoadMetric() (err error) {
	if !s.cfg.FlagRestore {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := os.ReadFile(s.cfg.FlagFileStoragePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %s", err)
	}
	if err = json.Unmarshal(res, &s.data); err != nil {
		return fmt.Errorf("ошибка Unmarshal: %s", err)
	}
	return err
}
