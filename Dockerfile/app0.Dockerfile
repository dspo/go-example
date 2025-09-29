FROM ubuntu/go:1.24-25.04_edge AS builder

WORKDIR /go/src/app

COPY . .

RUN go build -o app0 cmd/app0/main.go

FROM ubuntu
COPY --from=builder /go/src/app/app0 /app0

EXPOSE 8080

ENTRYPOINT ["/app0"]
