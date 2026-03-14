FROM golang:1.25-alpine AS builder
ARG VERSION=dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o /nimsforestwebviewer ./cmd/nimsforestwebviewer

FROM alpine:3.21
COPY --from=builder /nimsforestwebviewer /usr/local/bin/nimsforestwebviewer
ENTRYPOINT ["nimsforestwebviewer"]
