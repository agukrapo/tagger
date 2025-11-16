FROM golang:1 AS builder
WORKDIR /usr/src/app
COPY go.mod .
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o tagger ./cmd

FROM alpine:3
RUN apk add --no-cache git
COPY --from=builder /usr/src/app/tagger /usr/local/bin/tagger
ENTRYPOINT ["tagger"]
