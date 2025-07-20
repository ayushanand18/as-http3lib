#!/bin/sh
set -euo pipefail

CERT_DIR="$(dirname "$0")"
CERT_NAME="localhost"

cd "$CERT_DIR"

echo "Generating self-signed certificate..."

openssl req -x509 -newkey rsa:4096 -sha256 -days 365 \
  -nodes -keyout key.pem -out cert.pem \
  -subj "/CN=${CERT_NAME}" \
  -addext "subjectAltName=DNS:${CERT_NAME}"

echo "Certificate generated: cert.pem and key.pem"
