#!/bin/bash
# Script to generate self-signed certificates for development

CERT_DIR="./infrastructure/haproxy/certs"
CERT_NAME="selfsigned"
DAYS_VALID=365

echo "====================================="
echo "  SSL Certificate Generation Script"
echo "====================================="
echo ""

# Create certificate directory
if [ ! -d "$CERT_DIR" ]; then
    echo "[*] Creating certificate directory: $CERT_DIR"
    mkdir -p "$CERT_DIR"
    chmod 700 "$CERT_DIR"
fi

# Check if OpenSSL is installed
if ! command -v openssl &> /dev/null; then
    echo "[ERROR] OpenSSL is not installed."
    echo "Please install OpenSSL: https://www.openssl.org"
    exit 1
fi

echo "[*] OpenSSL found: $(openssl version)"
echo ""

# Generate private key
echo "[*] Generating private key..."
openssl genrsa -out "$CERT_DIR/$CERT_NAME.key" 2048 2>/dev/null

if [ $? -ne 0 ]; then
    echo "[ERROR] Failed to generate private key"
    exit 1
fi

echo "[✓] Private key generated: $CERT_DIR/$CERT_NAME.key"
echo ""

# Generate certificate signing request (CSR)
echo "[*] Generating Certificate Signing Request (CSR)..."
openssl req -new \
    -key "$CERT_DIR/$CERT_NAME.key" \
    -out "$CERT_DIR/$CERT_NAME.csr" \
    -subj "/C=ES/ST=Madrid/L=Madrid/O=SIG-Agro/CN=localhost" 2>/dev/null

if [ $? -ne 0 ]; then
    echo "[ERROR] Failed to generate CSR"
    exit 1
fi

echo "[✓] CSR generated: $CERT_DIR/$CERT_NAME.csr"
echo ""

# Self-sign the certificate
echo "[*] Signing certificate for $DAYS_VALID days..."
openssl x509 -req \
    -days $DAYS_VALID \
    -in "$CERT_DIR/$CERT_NAME.csr" \
    -signkey "$CERT_DIR/$CERT_NAME.key" \
    -out "$CERT_DIR/$CERT_NAME.crt" 2>/dev/null

if [ $? -ne 0 ]; then
    echo "[ERROR] Failed to sign certificate"
    exit 1
fi

echo "[✓] Certificate generated: $CERT_DIR/$CERT_NAME.crt"
echo ""

# Create combined PEM file for HAProxy
echo "[*] Creating combined PEM file for HAProxy..."
cat "$CERT_DIR/$CERT_NAME.crt" "$CERT_DIR/$CERT_NAME.key" \
    > "$CERT_DIR/$CERT_NAME.pem"

if [ $? -ne 0 ]; then
    echo "[ERROR] Failed to create PEM file"
    exit 1
fi

chmod 600 "$CERT_DIR/$CERT_NAME.pem"
echo "[✓] PEM file created: $CERT_DIR/$CERT_NAME.pem"
echo ""

# Display certificate information
echo "[*] Certificate Information:"
openssl x509 -in "$CERT_DIR/$CERT_NAME.crt" -text -noout | grep -A 2 "Subject:\|Issuer:\|Not Before\|Not After"

echo ""
echo "[✓] Certificate generation complete!"
echo ""
echo "Files created:"
echo "  - Private key:  $CERT_DIR/$CERT_NAME.key"
echo "  - Certificate:  $CERT_DIR/$CERT_NAME.crt"
echo "  - CSR:          $CERT_DIR/$CERT_NAME.csr"
echo "  - PEM:          $CERT_DIR/$CERT_NAME.pem (for HAProxy)"
echo ""
echo "Note: These are self-signed certificates for DEVELOPMENT only."
echo "For production, obtain proper certificates from a trusted CA."
echo ""
