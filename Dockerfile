FROM golang:alpine3.20 AS builder

RUN apk add --no-cache git

# Save shared libraries that is required for git to work
RUN mkdir -p /libs \
    && cp /lib/ld-musl-x86_64.so.1 /libs/ \
    && cp /lib/libc.musl-x86_64.so.1 /libs/ \
    && cp /lib/libz.so.1 /libs/ \
    && cp /usr/lib/libpcre2-8.so.0 /libs/

WORKDIR /app
COPY . /app

# Build gitmate without debug symbols
RUN go build -ldflags "-w" -o /usr/local/bin/gitmate main.go

FROM busybox:1.36.0-uclibc

COPY --from=builder /usr/local/bin/gitmate /usr/bin/
COPY --from=builder /usr/bin/git /usr/bin/

# Copy shared libraries from builder
COPY --from=builder /libs/ /lib/

WORKDIR /app

# Configure Git to trust the /app directory to avoid "dubious ownership" errors
RUN git config --global --add safe.directory /app

CMD ["gitmate"]

