FROM golang:1 AS builder
WORKDIR /usr/src/app
COPY go.mod .
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o tagger ./cmd

FROM aplpine:3
WORKDIR /tagger
COPY --from=builder /usr/src/app/tagger .
ENTRYPOINT ["./tagger"]
