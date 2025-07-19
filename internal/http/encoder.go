package http

import (
	"context"
	"encoding/json"
	"net/http"
)

func DefaultHttpEncode(ctx context.Context, response interface{}) (headers map[string][]string, body []byte, err error) {
	headers = map[string][]string{
		"Content-Type": {"application/json; charset=utf-8"},
	}

	switch v := response.(type) {
	case string:
		body = []byte(v)
	case []byte:
		body = v
	default:
		body, err = json.Marshal(v)
		if err != nil {
			return headers, nil, err
		}
	}

	return headers, body, nil
}

func DefaultHttpDecode(ctx context.Context, r *http.Request) (outgoingRequest interface{}, err error) {
	if e := json.NewDecoder(r.Body).Decode(&outgoingRequest); e != nil {
		return outgoingRequest, err
	}

	return outgoingRequest, nil
}
