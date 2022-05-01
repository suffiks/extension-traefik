FROM golang:1.18.1-alpine as builder

RUN apk add --no-cache git

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /workspace

# Build
RUN CGO_ENABLED=0 go build -a -o traefik ./cmd/traefik

FROM alpine
WORKDIR /
COPY --from=builder /workspace/traefik .
COPY ./docs/ ./docs/
USER 65532:65532

ENTRYPOINT ["/traefik"]
