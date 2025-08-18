package http3

import (
	"context"
	"net/http"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	ashttp "github.com/ayushanand18/as-http3lib/internal/http"
	"github.com/ayushanand18/as-http3lib/pkg/types"
	"github.com/gorilla/mux"
)

func httpDefaultHandler(
	ctx context.Context,
	w http.ResponseWriter,
	handler types.HandlerFunc,
	decoder types.HttpDecoder,
	encoder types.HttpEncoder,
	r *http.Request) {

	var request interface{}
	var err error

	ctx, err = defaultMiddleware(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if decoder != nil {
		request, err = decoder(ctx, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		request, err = ashttp.DefaultHttpDecode(ctx, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp, err := handler(ctx, request)
	if err != nil {
		return
	}

	headers := make(map[string][]string)
	var body []byte
	if encoder != nil {
		headers, body, err = encoder(ctx, resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		headers, body, err = ashttp.DefaultHttpEncode(ctx, resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	headers = ashttp.PopulateDefaultServerHeaders(ctx, headers)

	for key, value := range headers {
		w.Header().Del(key)
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	if body != nil {
		_, err := w.Write(body)
		if err != nil {
			panic(err)
		}
	}
}

func defaultMiddleware(ctx context.Context, r *http.Request) (outgoingContext context.Context, err error) {
	ctx = context.WithValue(ctx, constants.HttpRequestHeaders, r.Header)

	params := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	ctx = context.WithValue(ctx, constants.HttpRequestURLParams, params)

	ctx = context.WithValue(ctx, constants.HttpRequestPathValues, mux.Vars(r))

	return ctx, nil
}
