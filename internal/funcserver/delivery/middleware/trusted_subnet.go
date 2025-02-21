package middleware

import (
	"fmt"
	"github.com/LI-SeNyA-vE/KursMetrics/pkg/utils/ipandcidr"
	"net/http"
)

func (m *Middleware) TrustedSubnet(h http.Handler) http.Handler {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		ok, err := ipandcidr.IsIPInCIDR(ip, m.Server.FlagTrustedSubnet)
		if err != nil {
			http.Error(w, fmt.Sprintf("на этапе проверки ipandcidr произошла ошибка: %v", err), http.StatusInternalServerError)
			m.log.Errorf("Запрос не обработан, так как ipandcidr: %s передан невалидным", m.Server.FlagTrustedSubnet)
			return
		}
		if !ok {
			http.Error(w, "данный ipandcidr не входит в зону доверенных", http.StatusConflict)
			m.log.Errorf("Запрос не обработан, так как ipandcidr %s, не входит в доверенную зону: %s", ip, m.Server.FlagTrustedSubnet)
			return
		}
		// Передаём управление следующему хендлеру
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hashFn)
}
