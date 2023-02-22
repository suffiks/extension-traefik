FROM golang:1.20-alpine as builder

RUN apk add --no-cache git

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /workspace

# Build
RUN CGO_ENABLED=0 go build -a -o traefik ./

FROM alpine
WORKDIR /
COPY --from=builder /workspace/traefik .
USER 65532:65532

ENTRYPOINT ["/traefik"]
