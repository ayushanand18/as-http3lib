package types

import (
	"net/http"

	"github.com/ayushanand18/as-http3lib/internal/constants"
)

type ServeOptions struct {
	URL     string
	Handler func(http.ResponseWriter, *http.Request)
	Method  constants.HttpMethodTypes
}
