FROM golang:1.22-alpine3.18 AS builder

RUN apk upgrade --no-cache && apk add --no-cache libgcc gcc musl-dev bind-tools libffi libffi-dev

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /app/main main.go


FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main /app/.env /app

EXPOSE 8888

CMD ["/app/main"]
