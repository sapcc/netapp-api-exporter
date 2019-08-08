FROM alpine:3.5

WORKDIR /app
COPY bin/netapp-api-exporter_linux_amd64 /app/netapp-api-exporter

EXPOSE 9108
ENTRYPOINT [ "./netapp-api-exporter" ]
