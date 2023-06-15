FROM golang:1.20.2-alpine3.17 as builder
ADD . /go/src/qbot
WORKDIR /go/src/qbot
RUN go install -ldflags="-s -w"

FROM alpine:3.17
COPY --from=builder /go/bin/qbot /app/qbot
ENTRYPOINT ["/app/qbot"]