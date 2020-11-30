# Output files
# ca.key: CA private key file (not shared)
# ca.crt: CA trust cert (shared with users)
# server.key: Server private key, password protected
# server.csr: Server cert signing request
# server.crt: Server certificate signed by CA - keep on server
# server.pem: Conversion of the server.key into gRPC readable format

# Summary
# Private files: ca.key, server.key, server.pem, server.crt
# "Share" files: ca.crt(needed by client), server.csr (needed by CA)

# Changes these CN(cert name) to match host in your environment
SERVER_CN=localhost

# Step 1: Generate CA + Trust Cert (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# Step 2: Generate the Server Private key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

# Step 4: Sign the certificate with CA (self-signing) - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# Step 5: Convert the server cert to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem