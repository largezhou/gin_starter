FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o ./tmp/main main.go

FROM alpine AS prod

WORKDIR /app

COPY --from=builder /app/ .

RUN apk update && apk add tzdata

CMD ["./tmp/main"]
