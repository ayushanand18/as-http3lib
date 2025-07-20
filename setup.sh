#!/usr/bin/env bash
set -euo pipefail

CERT_NAME="as-http3lib-local"
CERT_FILE="cert.pem"
KEY_FILE="key.pem"
EXAMPLES_DIR="./examples"

echo "=== Step 1: Configuring UDP buffer size (Linux only) ==="
if [[ "$(uname)" == "Linux" ]]; then
    sudo sysctl -w net.core.rmem_max=7500000
    sudo sysctl -w net.core.wmem_max=7500000
else
    echo "Skipping UDP buffer config on non-Linux system."
fi

echo "=== Step 2: Generating TLS certificate ==="
openssl req -x509 -newkey rsa:2048 -nodes \
    -keyout "${KEY_FILE}" -out "${CERT_FILE}" \
    -days 365 -subj "/CN=${CERT_NAME}" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

echo "=== Step 3: Copying cert/key to each examples/*/ folder ==="
for dir in "${EXAMPLES_DIR}"/*/; do
    cp "${CERT_FILE}" "${dir}"
    cp "${KEY_FILE}" "${dir}"
done

echo "=== Step 4: Trusting certificate ==="
if [[ "$(uname)" == "Darwin" ]]; then
    # macOS trust store
    sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "${CERT_FILE}"
    echo "Trusted the certificate on macOS."
elif [[ "$(uname)" == "Linux" ]]; then
    # Linux (Debian-based)
    if [[ -d /usr/local/share/ca-certificates ]]; then
        sudo cp "${CERT_FILE}" "/usr/local/share/ca-certificates/${CERT_NAME}.crt"
        sudo update-ca-certificates
        echo "Trusted the certificate on Linux."
    else
        echo "WARNING: CA directory not found, skipping trust setup. Please add ${CERT_FILE} manually."
    fi
else
    echo "WARNING: Unknown platform, please trust ${CERT_FILE} manually if needed."
fi

echo "=== Setup complete. ==="
