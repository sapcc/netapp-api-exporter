FROM alpine:3.5
LABEL source_repository="https://github.com/sapcc/netapp-api-exporter"

WORKDIR /app
COPY bin/netapp-api-exporter_linux_amd64 /app/netapp-api-exporter

EXPOSE 9108
ENTRYPOINT [ "./netapp-api-exporter" ]
