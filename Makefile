SHELL=bash

run: install
	. .env; discordbot

install:
	goimports -w .
	go mod tidy
	go install

makecert:
	openssl pkcs12 -in cert.p12 -clcerts -nokeys -out usercert.pem
	openssl pkcs12 -in cert.p12 -nocerts -out userkey.pem -nodes
