FROM alpine
MAINTAINER Feng Honglin <hfeng@tutum.co>

ADD events /events
ADD https://files.tutum.co/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV REPORT_INTERVAL=30 TUTUM_AUTH=**None** TUTUM_URL=**None**

ENTRYPOINT ["/events"]
