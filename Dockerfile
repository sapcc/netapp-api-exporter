FROM golang:1.14 AS builder

WORKDIR /build
COPY . /build
RUN go get -d
RUN go build -o netapp-api-exporter

FROM alpine:3.14
LABEL source_repository="https://github.com/sapcc/netapp-api-exporter"

RUN apk add --no-cache libc6-compat

WORKDIR /app
COPY --from=builder /build/netapp-api-exporter /app/netapp-api-exporter

EXPOSE 9108
ENTRYPOINT [ "./netapp-api-exporter" ]
