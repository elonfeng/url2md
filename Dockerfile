FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /url2md ./cmd/url2md

FROM alpine:3.20

RUN apk add --no-cache ca-certificates chromium

ENV CHROME_BIN=/usr/bin/chromium-browser

COPY --from=builder /url2md /usr/local/bin/url2md

EXPOSE 8080

ENTRYPOINT ["url2md"]
CMD ["serve", "--port", "8080"]
