FROM golang:1.22.3-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o gym_service

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/gym_service .

COPY .env .

EXPOSE 50001

CMD ["./gym_service"]
