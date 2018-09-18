FROM golang:1.11-alpine as builder
WORKDIR /app
COPY . .
RUN apk add --no-cache gcc git ca-certificates musl-dev
RUN go mod download; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o eir

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/eir /eir
CMD ["/eir"]
