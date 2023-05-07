FROM golang:1.20-alpine as build

WORKDIR /goback

COPY . /goback

RUN go mod tidy

RUN go build cmd/main.go

CMD ["./main"]
