FROM golang:1.11-alpine as builder
RUN apk add --no-cache gcc git ca-certificates musl-dev make tzdata zip

WORKDIR /app
COPY . .
RUN make build-linux
ENV TZ=America/Los_Angeles
CMD ["/app/eir"]
