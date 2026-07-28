#!/bin/bash

set -e

echo " GOS 2013 Certificate Generator"

rm -rf crt
mkdir crt

cd crt


SHA="sha1"
DAYS=7305


CA_DN="/C=US/ST=California/L=Redwood City/O=Electronic Arts, Inc./OU=Global Online Studio/CN=GOS 2013 Certificate Authority/emailAddress=GOSDirtysockSupport@ea.com"

SERVER_DN="/C=US/ST=California/L=Redwood City/O=Electronic Arts, Inc./OU=Global Online Studio/CN=gosredirector.ea.com/emailAddress=GOSDirtysockSupport@ea.com"

echo "[+] Generate CA RSA 1024"

openssl genrsa \
    -out ca.key \
    1024



echo "[+] Create CA config"

cat > ca.conf <<EOF
[req]
distinguished_name=req_dn
x509_extensions=v3_ca

[req_dn]

[v3_ca]
basicConstraints=CA:true
keyUsage=keyCertSign,cRLSign
EOF

echo "[+] Generate CA certificate"

openssl req \
    -x509 \
    -new \
    -nodes \
    -key ca.key \
    -sha1 \
    -days $DAYS \
    -out ca.pem \
    -subj "$CA_DN" \
    -config ca.conf

echo "[+] Generate server RSA 1024"

openssl genrsa \
    -out gosredirector.ea.com-priv.pem \
    1024

echo "[+] Create server config"

cat > server.conf <<EOF
[req]
distinguished_name=req_dn
req_extensions=v3_req

[req_dn]

[v3_req]

keyUsage=digitalSignature,keyEncipherment
subjectAltName=@alt_names


[alt_names]

DNS.1=gosredirector.ea.com
DNS.2=*.ea.com
DNS.3=*.easports.com

EOF

echo "[+] Generate CSR"

openssl req \
    -new \
    -key gosredirector.ea.com-priv.pem \
    -out server.csr \
    -config server.conf \
    -subj "$SERVER_DN"

echo "[+] Sign server certificate"

openssl x509 \
    -req \
    -in server.csr \
    -CA ca.pem \
    -CAkey ca.key \
    -CAcreateserial \
    -out gosredirector.ea.com-cert.pem \
    -days $DAYS \
    -sha1 \
    -extfile server.conf \
    -extensions v3_req

echo "[+] Create chain"

cat gosredirector.ea.com-cert.pem ca.pem > fullchain.pem

echo "[+] Verify"

openssl verify \
    -CAfile ca.pem \
    gosredirector.ea.com-cert.pem

echo
echo " DONE"

echo
echo "Generated:"
echo "gosredirector.ea.com-cert.pem"
echo "gosredirector.ea.com-priv.pem"
echo "ca.pem"
echo "fullchain.pem"

openssl x509 \
    -in gosredirector.ea.com-cert.pem \
    -text \
    -noout
