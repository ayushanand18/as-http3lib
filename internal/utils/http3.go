package utils

import (
	"fmt"

	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func ValidateOptionsBeforeRequest(options types.ServeOptions) error {
	if options.URL == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	if options.Handler == nil {
		return fmt.Errorf("handler cannot be nil for URL %s", options.URL)
	}
	if options.Method == "" {
		return fmt.Errorf("HTTP method cannot be empty for URL %s", options.URL)
	}

	return nil
}
