FROM golang:alpine as builder
ADD . /go/src/qbot
WORKDIR /go/src/qbot
RUN go install -ldflags="-s -w -X qbot/internal/meta.Version=$(git describe --tags)"

FROM alpine
COPY --from=builder /go/bin/qbot /app/qbot
ENTRYPOINT ["/app/qbot"]