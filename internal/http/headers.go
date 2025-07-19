package http

import "context"

func PopulateDefaultServerHeaders(ctx context.Context, headers map[string][]string) map[string][]string {
	if headers == nil {
		headers = make(map[string][]string)
	}

	headers["X-Server"] = []string{"as-http3lib"}

	return headers
}
