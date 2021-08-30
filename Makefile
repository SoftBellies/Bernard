run: install
	discordbot

install:
	goimports -w .
	go mod tidy
	go install
