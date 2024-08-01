FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -p $(nproc) -a -installsuffix cgo -o /app/main ./cmd/rest/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
EXPOSE 8081

CMD ["./main"]
