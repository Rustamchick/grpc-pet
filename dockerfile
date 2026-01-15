FROM golang:1.24.2 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main ./cmd/main

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/configs/config.yaml . 
COPY --from=builder /app/configs/dbConfig.yaml . 
COPY --from=builder /app/configs/local_tests.yaml . 

EXPOSE 9090

CMD ["./main"]