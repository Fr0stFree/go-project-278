FROM --platform=$BUILDPLATFORM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o ./bin/shortener ./cmd/shortener


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/bin/shortener ./shortener

EXPOSE 8080

ENTRYPOINT ["./shortener"]