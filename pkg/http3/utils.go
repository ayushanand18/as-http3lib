package http3

import (
	"context"
	"os"

	"github.com/ayushanand18/as-http3lib/internal/config"
)

func checkIfTlsCertificateIsMissing(ctx context.Context) bool {
	keyRawBytes := config.GetBytes(ctx, "service.tls.key.raw")
	if len(keyRawBytes) > 0 {
		return false
	}

	certRawBytes := config.GetBytes(ctx, "service.tls.certificate.raw")
	if len(certRawBytes) > 0 {
		return false
	}

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

func injectConstantHeaders() map[string]string {
	defaultHeaders := make(map[string]string)
	defaultHeaders["X-Server"] = "ashttp3lib"

	return defaultHeaders
}
