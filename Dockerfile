FROM golang:1.19-alpine3.15 AS builder

COPY . /telegramBot/
WORKDIR /telegramBot/

RUN go mod download
RUN go build -o ./bin/app cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /telegramBot/bin/app .
COPY --from=0 /telegramBot/configs configs/

EXPOSE 80

CMD ["./app"]