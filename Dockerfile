FROM golang:1.20.2-alpine3.17 as builder
ADD . /go/src/qb-monitor
WORKDIR /go/src/qb-monitor
RUN go install -ldflags="-s -w"

FROM alpine:3.17
COPY --from=builder /go/bin/qb-monitor /app/qb-monitor
ENTRYPOINT ["/app/qb-monitor"]