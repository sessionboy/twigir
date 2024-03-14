
openssl genrsa -out ./cert/server.key 2048
openssl req -new -x509 -key ./cert/server.key -out ./cert/server.pem -days 365