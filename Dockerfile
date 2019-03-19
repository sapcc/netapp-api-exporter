FROM alpine:3.5

WORKDIR /app
COPY netapp-api-exporter /app/

EXPOSE 9108
ENTRYPOINT [ "./netapp-api-exporter" ]
