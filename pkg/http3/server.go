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
		Server: qchttp3.Server{
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

	s.TLSConfig = tls.GenerateTLSConfig(ctx)
	s.Handler = &rootHandler{mux: s.mux}
	return nil
}

func (s *server) ListenAndServe() error {
	return s.Server.ListenAndServe()
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

			switch options.ResponseType {
			case constants.RESPONSE_TYPE_JSON_RESPONSE:
				jsonDefaultHandler(r.Context(), w, options, handler, r)
			case constants.RESPONSE_TYPE_STREAMING_RESPONSE:
				streamingDefaultHandler(r.Context(), w, options, handler, r)
			default:
				httpDefaultHandler(r.Context(), w, options, handler, r)
			}
		})
	}

	return nil
}
