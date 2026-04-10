FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -v -o /reconsiliation-api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -v -o /reconsiliation-consumer ./cmd/consumer/main.go

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=build /reconsiliation-api /reconsiliation-api
COPY --from=build /reconsiliation-consumer /reconsiliation-consumer

USER nonroot:nonroot

ENTRYPOINT ["/reconsiliation-api"]
