FROM golang:1.18rc1-alpine as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /workspace

# Build
RUN go build -a -o traefik ./cmd/traefik

FROM alpine
WORKDIR /
COPY --from=builder /workspace/traefik .
USER 65532:65532

ENTRYPOINT ["/traefik"]
