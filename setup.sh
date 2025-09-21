#!/usr/bin/env bash
set -euo pipefail

CERT_NAME="crazyhttp-local"
CA_KEY="ca.key.pem"
CA_CERT="ca.cert.pem"
SERVER_KEY="key.pem"
SERVER_CSR="server.csr.pem"
SERVER_CERT="cert.pem"
EXAMPLES_DIR="./examples"

echo "=== Step 1: Configuring UDP buffer size (Linux only) ==="
if [[ "$(uname)" == "Linux" ]]; then
    sudo sysctl -w net.core.rmem_max=7500000
    sudo sysctl -w net.core.wmem_max=7500000
else
    echo "Skipping UDP buffer config on non-Linux system."
fi

echo "=== Step 2: Generating CA (if not exists) ==="
if [[ ! -f "${CA_CERT}" || ! -f "${CA_KEY}" ]]; then
    openssl req -x509 -newkey rsa:4096 -days 3650 -nodes \
        -keyout "${CA_KEY}" -out "${CA_CERT}" \
        -subj "/CN=${CERT_NAME} Root CA" \
        -addext "basicConstraints=critical,CA:TRUE" \
        -addext "keyUsage=keyCertSign,cRLSign"
    echo "Generated CA certificate."
else
    echo "CA already exists, skipping generation."
fi

echo "=== Step 3: Generating server private key and CSR ==="
openssl req -newkey rsa:2048 -nodes -keyout "${SERVER_KEY}" -out "${SERVER_CSR}" \
    -subj "/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

echo "=== Step 4: Signing server certificate with CA ==="
openssl x509 -req -in "${SERVER_CSR}" -CA "${CA_CERT}" -CAkey "${CA_KEY}" \
    -CAcreateserial -out "${SERVER_CERT}" -days 365 \
    -extfile <(printf "subjectAltName=DNS:localhost,IP:127.0.0.1")

echo "=== Step 5: Copying cert/key to each examples/*/ folder ==="
for dir in "${EXAMPLES_DIR}"/*/; do
    cp "${SERVER_CERT}" "${dir}"
    cp "${SERVER_KEY}" "${dir}"
done

echo "=== Step 6: Trusting CA certificate ==="
if [[ "$(uname)" == "Darwin" ]]; then
    sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "${CA_CERT}"
    echo "Trusted the CA on macOS."
elif [[ "$(uname)" == "Linux" ]]; then
    if [[ -d /usr/local/share/ca-certificates ]]; then
        sudo cp "${CA_CERT}" "/usr/local/share/ca-certificates/${CERT_NAME}.crt"
        sudo update-ca-certificates
        echo "Trusted the CA on Linux."
    else
        echo "WARNING: CA directory not found, skipping trust setup. Please add ${CA_CERT} manually."
    fi
else
    echo "WARNING: Unknown platform, please trust ${CA_CERT} manually if needed."
fi

echo "=== Setup complete. Use '${SERVER_CERT}' and '${SERVER_KEY}' for your server. ==="
