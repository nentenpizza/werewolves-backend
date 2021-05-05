FROM golang:1.13

WORKDIR /storage

COPY . .

RUN go get -u github.com/pressly/goose/cmd/goose

ENTRYPOINT ["goose", "up"]