package http3

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ayushanand18/as-http3lib/internal/config"
	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/internal/tls"
	"github.com/ayushanand18/as-http3lib/internal/utils"
	"github.com/ayushanand18/as-http3lib/pkg/types"
	qchttp3 "github.com/quic-go/quic-go/http3"
)

type server struct {
	h3Server      qchttp3.Server
	mux           *http.ServeMux
	routeMatchMap map[string]map[constants.HttpMethodTypes]types.HandlerFunc
}

type Server interface {
	Initialize(context.Context) error
	ListenAndServe() error
	AddServeMethod(context.Context, types.ServeOptions) error
}

func NewServer(ctx context.Context) Server {
	return &server{
		h3Server: qchttp3.Server{
			Addr:    utils.GetListeningAddress(ctx),
			Handler: nil,
		},
		mux:           http.NewServeMux(),
		routeMatchMap: make(map[string]map[constants.HttpMethodTypes]types.HandlerFunc),
	}
}

func (s *server) Initialize(ctx context.Context) error {
	if config.GetBool(ctx, "service.tls.generate_if_missing", true) && checkIfTlsCertificateIsMissing(ctx) {
		if err := tls.GenerateSelfSignedCert(ctx); err != nil {
			return fmt.Errorf("failed to generate self-signed certificate: %v", err)
		}
	}

	s.h3Server.TLSConfig = tls.GenerateTLSConfig(ctx)
	s.h3Server.Handler = s.mux
	return nil
}

func (s *server) ListenAndServe() error {
	return s.h3Server.ListenAndServe()
}

func (s *server) AddServeMethod(ctx context.Context, options types.ServeOptions) error {
	if err := utils.ValidateOptionsBeforeRequest(options); err != nil {
		return err
	}

	if _, ok := s.routeMatchMap[options.URL]; !ok {
		s.routeMatchMap[options.URL] = make(map[constants.HttpMethodTypes]types.HandlerFunc)
	}

	// if the combination exists, reassign it
	s.routeMatchMap[options.URL][options.Method] = options.Handler
	if len(s.routeMatchMap[options.URL]) == 1 {
		s.mux.HandleFunc(options.URL, func(w http.ResponseWriter, r *http.Request) {
			requestMethod := constants.HttpMethodTypes(strings.ToUpper(r.Method))

			methodHandlers, ok := s.routeMatchMap[options.URL]
			if !ok {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			handler, ok := methodHandlers[requestMethod]
			if !ok {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			response := handler(ctx, r)
			if response.StatusCode == 0 {
				response.StatusCode = http.StatusOK
			}
			w.WriteHeader(response.StatusCode)

			defaultHeaders := injecteConstantHeaders()
			for key, value := range defaultHeaders {
				w.Header().Set(key, value)
			}

			for key, value := range response.Headers {
				w.Header().Set(key, value)
			}

			if response.Body != nil {
				_, err := w.Write(response.Body)
				if err != nil {
					panic(err)
				}
			}
		})
	}

	return nil
}
