package http3

import (
	"context"
	"os"

	"github.com/ayushanand18/as-http3lib/internal/config"
)

func checkIfTlsCertificateIsMissing(ctx context.Context) bool {
	var keyFile, certFile string

	keyFile = config.GetString(ctx, "service.tls.key.path", "key.pem")
	certFile = config.GetString(ctx, "service.tls.certificate.path", "cert.pem")

	_, certErr := os.Stat(certFile)
	_, keyErr := os.Stat(keyFile)

	if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
		return true
	}

	return false
}

func injecteConstantHeaders() map[string]string {
	defaultHeaders := make(map[string]string)
	defaultHeaders["X-Server"] = "ashttp3lib"

	return defaultHeaders
}
