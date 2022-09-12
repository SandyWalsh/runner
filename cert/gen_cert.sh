#!/bin/bash

# server
openssl genrsa -out ca.key 4096
openssl req -new -x509 -key ca.key -sha256 -subj "/C=CA/ST=NS/O=Dark Secret Software, Inc." -days 365 -out ca.cert
openssl genrsa -out server.key 4096
openssl req -new -key server.key -out server.csr -config server-cert.conf
openssl x509 -req -in server.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out server.cert -days 365 -sha256 -extfile server-cert.conf -extensions req_ext

# client
openssl req -newkey rsa:4096 -nodes -keyout client.key -out client.csr -subj "/C=CA/ST=NS/O=Dark Secret Software, Inc./UID=cert-A"
openssl x509 -req -in client.csr -days 60 -CA ca.cert -CAkey ca.key -CAcreateserial -out client.cert -extfile client-cert.conf
#openssl x509 -in client-cert.pem -noout -text

