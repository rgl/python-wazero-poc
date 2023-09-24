FROM golang:1.21-bookworm AS build
RUN <<EOF
set -eux
apt-get update
apt-get install -y \
    unzip
rm -rf /var/lib/apt/lists/*
EOF
WORKDIR /build
COPY Makefile ./
RUN make python.wasm lib
COPY go.* .
RUN go mod download
COPY *.go *.py ./
RUN CGO_ENABLED=0 go build -ldflags="-s"

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /build/lib ./lib/
COPY --from=build /build/example ./
# NB run it once to compile the wasm to native code and cache it.
# TODO try to embed the native code into the example binary instead.
RUN ./example
ENTRYPOINT ["./example"]
