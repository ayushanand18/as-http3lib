package tls

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ayushanand18/as-http3lib/internal/config"
)

func GenerateSelfSignedCert(ctx context.Context) error {
	var keyFile, certFile string

	keyFile = config.GetString(ctx, "service.tls.key.path", "key.pem")
	certFile = config.GetString(ctx, "service.tls.certificate.path", "cert.pem")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA private key: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,

		DNSNames: []string{"localhost"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %w", certFile, err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write data to %s: %w", certFile, err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("error closing %s: %w", certFile, err)
	}

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %w", keyFile, err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}); err != nil {
		return fmt.Errorf("failed to write data to %s: %w", keyFile, err)
	}
	if err := keyOut.Close(); err != nil {
		return fmt.Errorf("error closing %s: %w", keyFile, err)
	}

	return nil
}

func GenerateTLSConfig(ctx context.Context) *tls.Config {
	var err error
	keyBytesRaw := config.GetBytes(ctx, "service.tls.key.raw")
	if len(keyBytesRaw) <= 0 {
		keyFile := config.GetString(ctx, "service.tls.key.path", "key.pem")
		keyBytesRaw, err = os.ReadFile(keyFile)
		if err != nil {
			return &tls.Config{}
		}
	}
	certBytesRaw := config.GetBytes(ctx, "service.tls.certificate.raw")
	if len(certBytesRaw) <= 0 {
		certFile := config.GetString(ctx, "service.tls.certificate.path", "cert.pem")
		certBytesRaw, err = os.ReadFile(certFile)
		if err != nil {
			return &tls.Config{}
		}
	}

	cert, err := tls.X509KeyPair(certBytesRaw, keyBytesRaw)
	if err != nil {
		log.Fatalf("failed to load TLS key pair: %v", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}
