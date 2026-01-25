package handler

import (
	"net"
	"net/http"
)

// TrustedSubnetOnly the client's IP address transmitted in the X-Real-IP request header is part of a trusted subnet
func (h *Handler) TrustedSubnetOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// если trusted_subnet не задан — доступ запрещён
		if h.cfg.TrustedSubnet == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		_, subnet, err := net.ParseCIDR(h.cfg.TrustedSubnet)
		if err != nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		ipStr := r.Header.Get("X-Real-IP")
		if ipStr == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		ip := net.ParseIP(ipStr)
		if ip == nil || !subnet.Contains(ip) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
