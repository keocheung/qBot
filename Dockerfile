FROM golang:1.20.2-alpine3.17
RUN go build
CMD qb-monitor