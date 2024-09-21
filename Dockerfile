FROM golang:1.18 as builder

ARG GOARCH=amd64
ARG GOOS=linux

COPY . /src
WORKDIR /src
RUN GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build main.go

FROM alpine:3.12.1
COPY --from=builder /src/main /airport-subscription-exporter
EXPOSE 9233
ENTRYPOINT ["/airport-subscription-exporter"]
