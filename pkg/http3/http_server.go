package http3

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ayushanand18/as-http3lib/internal/config"
	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/internal/tls"
	"github.com/ayushanand18/as-http3lib/internal/utils"
)

func (s *server) Initialize(ctx context.Context) error {
	if config.GetBool(ctx, "service.tls.generate_if_missing", true) && checkIfTlsCertificateIsMissing(ctx) {
		if err := tls.GenerateSelfSignedCert(ctx); err != nil {
			return fmt.Errorf("failed to generate self-signed certificate: %v", err)
		}
	}

	tlsConfig := tls.GenerateTLSConfig(ctx)
	root := &rootHandler{mux: s.mux}

	s.h3server.Handler = root
	s.h3server.TLSConfig = tlsConfig
	s.h3server.TLSConfig.NextProtos = []string{"h3"}

	h1Handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if on H/1 advertise H/3
		w.Header().Set("Alt-Svc", fmt.Sprintf(`h3=":%s"; ma=2592000,h3-29=":%s"; ma=2592000`, s.h3server.Addr[strings.LastIndex(s.h3server.Addr, ":")+1:], s.h3server.Addr[strings.LastIndex(s.h3server.Addr, ":")+1:]))
		root.ServeHTTP(w, r)
	})
	s.http1Server.Handler = h1Handler
	s.http1ServerTLS.Handler = h1Handler
	s.http1ServerTLS.TLSConfig = tlsConfig

	s.mcpServer.Handler = h1Handler
	return nil
}

func (s *server) ListenAndServe(ctx context.Context) error {
	utils.PrintStartBanner()
	errChan := make(chan error, 3)

	go func() {
		// server over http/3
		slog.InfoContext(ctx, "Starting HTTP/3 server", "port", s.h3server.Addr)
		errChan <- s.h3server.ListenAndServe()
	}()

	go func() {
		slog.InfoContext(ctx, "Starting HTTP/1.1 + Alt-Svc server", "port", s.http1Server.Addr)
		// server over http
		errChan <- s.http1Server.ListenAndServe()
	}()

	go func() {
		slog.InfoContext(ctx, "Starting HTTP/1.1 + Alt-Svc server (HTTPS)", "port", s.http1ServerTLS.Addr)
		// server over https
		errChan <- s.http1ServerTLS.ListenAndServeTLS("", "")
	}()

	return <-errChan
}

func (s *server) GET(url string) Method {
	return NewMethod(constants.HttpMethodGet, url, s)
}

func (s *server) POST(url string) Method {
	return NewMethod(constants.HttpMethodPost, url, s)

}
func (s *server) PUT(url string) Method {
	return NewMethod(constants.HttpMethodPut, url, s)
}

func (s *server) PATCH(url string) Method {
	return NewMethod(constants.HttpMethodPatch, url, s)
}

func (s *server) DELETE(url string) Method {
	return NewMethod(constants.HttpMethodDelete, url, s)
}

func (s *server) HEAD(url string) Method {
	return NewMethod(constants.HttpMethodHead, url, s)
}

func (s *server) OPTIONS(url string) Method {
	return NewMethod(constants.HttpMethodOptions, url, s)
}

func (s *server) CONNECT(url string) Method {
	return NewMethod(constants.HttpMethodConnect, url, s)
}

func (s *server) TRACE(url string) Method {
	return NewMethod(constants.HttpMethodTrace, url, s)
}
