ARG GO_VERSION=1.14
ARG ALPINE_VERSION=3.10.3

FROM golang:${GO_VERSION}-alpine AS builder

ARG SEABOLT_VERSION=v1.7.4

RUN apk add --update --no-cache ca-certificates cmake make g++ openssl-dev openssl-libs-static git curl pkgconfig libcap
RUN git clone -b ${SEABOLT_VERSION} https://github.com/neo4j-drivers/seabolt.git /seabolt
RUN update-ca-certificates 2>/dev/null || true

WORKDIR /seabolt/build

RUN cmake -D CMAKE_BUILD_TYPE=Release -D CMAKE_INSTALL_LIBDIR=lib .. && cmake --build . --target install

RUN mkdir -p /go/src/github.com/charlysan/goneo4jgql 
RUN mkdir /build
ADD . /go/src/github.com/charlysan/goneo4jgql/
WORKDIR /go/src/github.com/charlysan/goneo4jgql 

RUN go get github.com/vektah/gorunpkg
RUN go generate ./...

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags seabolt_static -o /app cmd/main.go

# Create alpine runtime image
FROM alpine:${ALPINE_VERSION} as app

# Environment variables
ENV API_PORT '8080'
ENV LOGGER_FORMATTER 'console'
ENV LOGGER_LEVEL 'debug'
ENV NEO4J_HOST 'localhost'
ENV NEO4J_PORT '7687'
ENV NEO4J_USER 'neo4j'
ENV NEO4J_PASS 'test'
ENV NEO4J_PROTO 'bolt'

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

USER 1000

EXPOSE 80

ENTRYPOINT ["/app"]