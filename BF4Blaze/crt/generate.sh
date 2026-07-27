#!/bin/bash

set -e

ORG="Electronic Arts"
OU="EA Secure"
CA_NAME="EA Secure Certificate Authority"

DAYS=3650

echo "===================================="
echo " EA Multi Domain Certificate"
echo "===================================="


rm -rf crt
mkdir crt

cd crt


echo "[+] Generating Fake EA Root CA key"

openssl genrsa \
    -out ca.key \
    1024


echo "[+] Creating CA config"

cat > ca.conf <<EOF
[req]
distinguished_name=req_dn
x509_extensions=v3_ca

[req_dn]

[v3_ca]
basicConstraints=critical,CA:true
keyUsage=critical,keyCertSign,cRLSign
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
EOF


echo "[+] Generating CA certificate"

openssl req \
    -x509 \
    -new \
    -nodes \
    -key ca.key \
    -sha256 \
    -days $DAYS \
    -out ca.pem \
    -config ca.conf \
    -subj "/C=US/O=$ORG/OU=$OU/CN=$CA_NAME"



echo "[+] Generating server key"

openssl genrsa \
    -out privkey.pem \
    1024

echo "[+] Creating server config"

cat > server.conf <<EOF
[req]
distinguished_name=req_dn
req_extensions=req_ext

[req_dn]

[req_ext]
subjectAltName=@alt_names

[alt_names]
DNS.1=gosredirector.ea.com
DNS.2=bf4.gos.ea.com
EOF

echo "[+] Creating CSR"

openssl req \
    -new \
    -key privkey.pem \
    -out server.csr \
    -config server.conf \
    -subj "/C=US/O=$ORG/OU=$OU/CN=gosredirector.ea.com"

echo "[+] Signing certificate"

openssl x509 \
    -req \
    -in server.csr \
    -CA ca.pem \
    -CAkey ca.key \
    -CAcreateserial \
    -out server.pem \
    -days $DAYS \
    -sha256 \
    -extfile server.conf \
    -extensions req_ext

echo "[+] Creating chain"

cat server.pem ca.pem > fullchain.pem

echo "[+] Verify"

openssl verify \
    -CAfile ca.pem \
    server.pem

echo
echo "===================================="
echo " Generated"
echo "===================================="

echo
echo "CA:"
echo " crt/ca.pem"

echo
echo "Server:"
echo " crt/server.pem"

echo
echo "Key:"
echo " crt/privkey.pem"

echo
echo "Chain:"
echo " crt/fullchain.pem"

echo
echo "Check SAN:"
openssl x509 \
    -in server.pem \
    -text \
    -noout | grep DNS