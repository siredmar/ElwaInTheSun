FROM golang:1.23.1-alpine AS builder
RUN apk add --no-cache make
COPY . /go
WORKDIR /go
RUN make build

FROM alpine:3.9
COPY --from=builder /go/bin/controller /controller
CMD ["/controller"]
