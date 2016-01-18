FROM alpine
MAINTAINER Feng Honglin <hfeng@tutum.co>

ADD events /events
ADD https://files.tutum.co/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV REPORT_INTERVAL=30 DOCKERCLOUD_AUTH=**None** EVENTS_API_URL=**None**

ENTRYPOINT ["/events"]
