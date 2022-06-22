# base go image
FROM golang:1.18.3-alpine3.16

WORKDIR /app

COPY . .

RUN go build -o main main.go

EXPOSE 3000

CMD [ "/app/voucherapi/main" ]

# build a tiny docker image


