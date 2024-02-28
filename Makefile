build:
	@go build -v -o go-monitor .

run:
	@echo "RUN go-monitor..."
	make build
	@./go-monitor