package http3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ayushanand18/as-http3lib/internal/config"
	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/internal/tls"
	"github.com/ayushanand18/as-http3lib/internal/utils"
	"github.com/ayushanand18/as-http3lib/pkg/types"
	"github.com/gorilla/mux"
	"github.com/quic-go/quic-go"
	qchttp3 "github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
)

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range injectConstantHeaders() {
		w.Header().Set(k, v)
	}

	rec := &responseRecorder{ResponseWriter: w, status: 0}
	h.mux.ServeHTTP(rec, r)

	if !rec.wrote {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("404 page not found\n"))
	}
}

type server struct {
	qchttp3.Server
	mux            *mux.Router
	routeMatchMap  map[string]map[constants.HttpMethodTypes]types.HandlerFunc
	http1ServerTLS http.Server
	http1Server    http.Server
}

type Server interface {
	Initialize(context.Context) error
	ListenAndServe(context.Context) error

	// HTTP Methods
	GET(string) Method
	POST(string) Method
	PUT(string) Method
	PATCH(string) Method
	DELETE(string) Method
	HEAD(string) Method
	OPTIONS(string) Method
	CONNECT(string) Method
	TRACE(string) Method
}

func NewServer(ctx context.Context) Server {
	quicConfig := &quic.Config{
		Tracer:          qlog.DefaultConnectionTracer,
		Allow0RTT:       true,
		EnableDatagrams: true,
	}
	return &server{
		Server: qchttp3.Server{
			Addr:            utils.GetListeningAddress(ctx),
			Handler:         nil,
			EnableDatagrams: true,
			QUICConfig:      quicConfig,
		},
		http1Server: http.Server{
			Addr: utils.GetHttp1ListeningAddress(ctx),
		},
		http1ServerTLS: http.Server{
			Addr: utils.GetHttp1TLSListeningAddress(ctx),
		},
		mux:           mux.NewRouter(),
		routeMatchMap: make(map[string]map[constants.HttpMethodTypes]types.HandlerFunc),
	}
}

func (s *server) Initialize(ctx context.Context) error {
	if config.GetBool(ctx, "service.tls.generate_if_missing", true) && checkIfTlsCertificateIsMissing(ctx) {
		if err := tls.GenerateSelfSignedCert(ctx); err != nil {
			return fmt.Errorf("failed to generate self-signed certificate: %v", err)
		}
	}

	tlsConfig := tls.GenerateTLSConfig(ctx)
	root := &rootHandler{mux: s.mux}

	s.Handler = root
	s.TLSConfig = tlsConfig
	s.TLSConfig.NextProtos = []string{"h3"}

	h1Handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if on H/1 advertise H/3
		w.Header().Set("Alt-Svc", fmt.Sprintf(`h3=":%s"; ma=2592000,h3-29=":%s"; ma=2592000`, s.Addr[strings.LastIndex(s.Addr, ":")+1:], s.Addr[strings.LastIndex(s.Addr, ":")+1:]))
		root.ServeHTTP(w, r)
	})
	s.http1Server.Handler = h1Handler
	s.http1ServerTLS.Handler = h1Handler
	s.http1ServerTLS.TLSConfig = tlsConfig

	return nil
}

func (s *server) ListenAndServe(ctx context.Context) error {
	utils.PrintStartBanner()
	errChan := make(chan error, 2)

	go func() {
		log.Println("Starting HTTP/3 server on", s.Server.Addr)
		errChan <- s.Server.ListenAndServe()
	}()

	go func() {
		log.Println("Starting HTTP/1.1 + Alt-Svc server on", s.http1Server.Addr)
		// server over http
		errChan <- s.http1Server.ListenAndServe()
	}()

	go func() {
		log.Println("Starting HTTP/1.1 + Alt-Svc server (HTTPS) on", s.http1ServerTLS.Addr)
		// server over https
		errChan <- s.http1ServerTLS.ListenAndServeTLS("", "")
	}()

	return <-errChan
}

func (s *server) GET(url string) Method {
	return NewMethod(constants.HTTP_METHOD_GET, url, s)
}

func (s *server) POST(url string) Method {
	return NewMethod(constants.HTTP_METHOD_POST, url, s)

}
func (s *server) PUT(url string) Method {
	return NewMethod(constants.HTTP_METHOD_PUT, url, s)
}

func (s *server) PATCH(url string) Method {
	return NewMethod(constants.HTTP_METHOD_PATCH, url, s)
}

func (s *server) DELETE(url string) Method {
	return NewMethod(constants.HTTP_METHOD_DELETE, url, s)
}

func (s *server) HEAD(url string) Method {
	return NewMethod(constants.HTTP_METHOD_HEAD, url, s)
}

func (s *server) OPTIONS(url string) Method {
	return NewMethod(constants.HTTP_METHOD_OPTIONS, url, s)
}

func (s *server) CONNECT(url string) Method {
	return NewMethod(constants.HTTP_METHOD_CONNECT, url, s)
}

func (s *server) TRACE(url string) Method {
	return NewMethod(constants.HTTP_METHOD_TRACE, url, s)
}
