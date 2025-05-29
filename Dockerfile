FROM golang:1.23-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./     

RUN go mod download      

COPY . .       

RUN go build -o main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /build/main .
COPY --from=build /build/.env .env
COPY --from=build /build/migrations ./migrations

CMD ["./main"]