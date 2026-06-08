FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o mlb-transform-api .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/mlb-transform-api .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./mlb-transform-api"]
