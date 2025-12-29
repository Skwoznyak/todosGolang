FROM golang:1.25

WORKDIR /app

COPY . .

ENV CONFIG_PATH=/app/config/local.yaml

RUN go mod tidy

CMD ["go", "run", "cmd/todolist/main.go"]