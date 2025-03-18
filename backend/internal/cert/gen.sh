rm *.pem

# 1. Generate CA's (certificate authority) private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=KE/ST=Nairobi/L=Nairobi/O=Kokomed-Fin/OU=Fintech/CN=OpenTelemetry CA/emailAddress=emiliocliff@gmail.com"

echo "CA's self-signed certificate"
openssl x509 -in ./ca-cert.pem -noout -text

# 2. Generate alloy server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -keyout alloy-key.pem -out alloy-req.pem -subj "/C=KE/ST=Nairobi/L=Nairobi/O=Alloy-Collector/OU=Fintech/CN=Alloy-Collector/emailAddress=emiliocliff@gmail.com"

# 3. Use CA's private key to sign alloys server's CSR and get back the signed certificate
openssl x509 -req -in ./alloy-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out alloy-cert.pem -extfile ./server-ext.cnf

echo "Alloy Server's signed certificate"
openssl x509 -in ./alloy-cert.pem -noout -text

# 4. Verify that alloy certificate is verified
# openssl verify -CAfile ./ca-cert.pem ./alloy-cert.pem

# subjectAltName=DNS:*.pcbook.com,DNS:*.pcbook.org,IP:0.0.0.0
