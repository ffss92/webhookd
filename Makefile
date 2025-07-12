target := bin/server

.PHONY: server/build
server/build:
	@go build -o $(target) cmd/server/*.go

.PHONY: server/run
server/run: server/build
	clear
	@$(target) -dev

.PHONY: server/watch
server/watch:
	reflex -r '\.go$$' -d none -t 15s -s make server/run
