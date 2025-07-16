package http3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func jsonDefaultHandler(ctx context.Context, w http.ResponseWriter, options types.ServeOptions, handler types.HandlerFunc, r *http.Request) {
	resp, err := handler(ctx, r)

	marshalledBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %+v", err), http.StatusInternalServerError)
		return
	}

	response := &types.HttpResponse{
		Body: marshalledBody,
	}
	if response.StatusCode == 0 {
		response.StatusCode = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
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
