package http

import "context"

func PopulateDefaultServerHeaders(ctx context.Context, headers map[string][]string) map[string][]string {
	if headers == nil {
		headers = make(map[string][]string)
	}

	headers["X-Server"] = []string{"crazyhttp"}
	headers["Access-Control-Allow-Origin"] = []string{"*"}
	headers["Access-Control-Allow-Methods"] = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	headers["Access-Control-Allow-Headers"] = []string{"Content-Type", "Authorization"}
	headers["Access-Control-Allow-Credentials"] = []string{"true"}
	headers["Access-Control-Max-Age"] = []string{"86400"}
	headers["Content-Type"] = []string{"text/event-stream"}

	return headers
}
