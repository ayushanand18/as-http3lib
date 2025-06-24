package http3

import (
	"context"
	"net/http"

	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func httpDefaultHandler(ctx context.Context, w http.ResponseWriter, options types.ServeOptions, handler types.HandlerFunc, r *http.Request) {
	resp := handler(ctx, r)

	response := resp.(*types.HttpResponse)
	if response.StatusCode == 0 {
		response.StatusCode = http.StatusOK
	}
	w.WriteHeader(response.StatusCode)

	for key, value := range options.DefaultHeaders {
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
}
