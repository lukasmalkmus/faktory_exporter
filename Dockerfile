FROM  quay.io/prometheus/busybox:latest
LABEL maintainer "Lukas Malkmus <mail@lukasmalkmus.com>"

COPY faktory_exporter /bin/faktory_exporter

EXPOSE      9386
ENTRYPOINT  ["/bin/faktory_exporter"]