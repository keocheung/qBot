FROM golang:alpine as builder
RUN apk add git
ADD . /go/src/qbot
WORKDIR /go/src/qbot
RUN go build -ldflags="-s -w -X qbot/internal/meta.Version=$(git describe --tags --always)"

FROM alpine
COPY --from=builder /go/src/qbot/qbot /app/qbot
ENTRYPOINT ["/app/qbot"]