package tools

import (
	"context"
	"net"
	"net/http"

	_ "net/http/pprof"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/config"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"
)

// SetupServers init servers
func SetupServers(ctx context.Context, cfg *config.Config, r *chi.Mux) (*http.Server, *http.Server) {
	var httpServer *http.Server
	if cfg.UseTLS {
		// конструируем менеджер TLS-сертификатов
		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
			// перечень доменов, для которых будут поддерживаться сертификаты
			HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
		}
		// конструируем сервер с поддержкой TLS
		httpServer = &http.Server{
			Addr:    cfg.RunAddr,
			Handler: r,
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: manager.TLSConfig(),
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		}
	} else {
		httpServer = &http.Server{
			Addr:    cfg.RunAddr,
			Handler: r,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		}
	}

	pprofServer := &http.Server{
		Addr: "localhost:6060",
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	return httpServer, pprofServer
}
