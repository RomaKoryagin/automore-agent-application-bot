FROM golang:alpine AS builder
WORKDIR /var/www/app
COPY . .
RUN apk add --no-cache gcc musl-dev sqlite-dev \
    && go mod download \
    && go build -o app cmd/main.go

FROM alpine:latest AS app
WORKDIR /var/www/app
RUN mkdir migrations
COPY --from=builder /var/www/app/app .
COPY --from=builder /var/www/app/.env .
COPY --from=builder /var/www/app/migrations migrations
COPY --from=builder /var/www/app/config config
CMD [ "./app" ]