FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD intelite-mqtt /app/intelite-mqtt
VOLUME /data
CMD ["/app/intelite-mqtt", "/data/config.yml"]
