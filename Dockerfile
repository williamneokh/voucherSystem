FROM goland:1.18-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o main voucherAPI/cmd/main.go

EXPOSE 3000
CMD [ "/app/voucherAPI/cmd/main" ]



