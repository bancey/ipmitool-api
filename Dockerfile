FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ipmitool-api .

FROM alpine:3.20

RUN apk add --no-cache ipmitool
COPY --from=builder /app/ipmitool-api /usr/local/bin/ipmitool-api

EXPOSE 8080
ENTRYPOINT ["ipmitool-api"]
CMD ["-config", "/etc/ipmitool-api/config.yaml"]
