FROM golang:1.15.5-alpine3.12 as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

COPY assets ./assets

RUN go mod download

COPY . .

RUN go build -o main ./cmd/apiserver

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

COPY --from=builder /app/assets ./assets

EXPOSE 7070

CMD ["./main"]