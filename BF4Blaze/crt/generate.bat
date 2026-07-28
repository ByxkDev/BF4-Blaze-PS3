@echo off
setlocal

set cdir=C:\Users\brams\OneDrive\Bureaublad\ssl-bf4

set ROOT=OTG3
set GOS2011=GOS2011
set GOS2013=GOS2013
set LEAF=gosredirector
set MOD=gosredirector_mod


cd /d "C:\Program Files\OpenSSL-Win64\bin"


echo ==========================
echo Creating OTG3 Root CA
echo ==========================

openssl genrsa ^
-out %cdir%\%ROOT%.key.pem ^
1024


openssl req ^
-new ^
-x509 ^
-md5 ^
-days 28124 ^
-key %cdir%\%ROOT%.key.pem ^
-out %cdir%\%ROOT%.crt ^
-subj "/C=US/ST=California/L=Redwood City/O=Electronic Arts, Inc./OU=Online Technology Group/CN=OTG3 Certificate Authority"



echo ==========================
echo Creating GOS 2011 CA
echo ==========================


openssl genrsa ^
-out %cdir%\%GOS2011%.key.pem ^
1024


openssl req ^
-new ^
-x509 ^
-md5 ^
-days 20000 ^
-key %cdir%\%GOS2011%.key.pem ^
-out %cdir%\%GOS2011%.crt ^
-subj "/C=US/ST=California/L=Redwood City/O=Electronic Arts, Inc./OU=Global Online Studio/CN=GOS 2011 Certificate Authority"



echo ==========================
echo Creating GOS 2013 CA
echo ==========================


openssl genrsa ^
-out %cdir%\%GOS2013%.key.pem ^
1024


openssl req ^
-new ^
-x509 ^
-md5 ^
-days 20000 ^
-key %cdir%\%GOS2013%.key.pem ^
-out %cdir%\%GOS2013%.crt ^
-subj "/C=US/ST=California/L=Redwood City/O=Electronic Arts, Inc./OU=Global Online Studio/CN=GOS 2013 Certificate Authority"



echo ==========================
echo Creating gosredirector key
echo ==========================


openssl genrsa ^
-out %cdir%\%LEAF%.key.pem ^
1024



echo ==========================
echo Creating CSR
echo ==========================


openssl req ^
-new ^
-key %cdir%\%LEAF%.key.pem ^
-out %cdir%\%LEAF%.csr ^
-subj "/C=US/ST=California/O=Electronic Arts, Inc./OU=Global Online Studio/CN=gosredirector.ea.com"



echo ==========================
echo Signing leaf
echo ==========================


openssl x509 ^
-req ^
-md5 ^
-days 10000 ^
-in %cdir%\%LEAF%.csr ^
-CA %cdir%\%GOS2013%.crt ^
-CAkey %cdir%\%GOS2013%.key.pem ^
-CAcreateserial ^
-out %cdir%\%LEAF%.crt



echo ==========================
echo Export DER
echo ==========================


openssl x509 ^
-in %cdir%\%LEAF%.crt ^
-outform der ^
-out %cdir%\%LEAF%.der



echo.
echo ===================================
echo APPLY PROTOSSL PATCH
echo ===================================

echo Open:
echo %cdir%\%LEAF%.der

echo.
echo Search:
echo 2a864886f70d010104

echo.
echo Replace ONLY second occurrence:

echo OLD:
echo 2a864886f70d010104

echo NEW:
echo 2a864886f70d010101

echo.
echo Save:
echo %cdir%\%MOD%.der


pause



echo ==========================
echo Convert patched DER
echo ==========================


openssl x509 ^
-inform der ^
-in %cdir%\%MOD%.der ^
-out %cdir%\%MOD%.crt



echo ==========================
echo Verify
echo ==========================


openssl x509 ^
-in %cdir%\%MOD%.crt ^
-text ^
-noout



echo ==========================
echo Create PFX
echo ==========================


openssl pkcs12 ^
-export ^
-out %cdir%\%MOD%.pfx ^
-inkey %cdir%\%LEAF%.key.pem ^
-in %cdir%\%MOD%.crt ^
-certfile %cdir%\%GOS2013%.crt ^
-passout pass:12345



echo ==========================
echo COMPLETE
echo ==========================

pause