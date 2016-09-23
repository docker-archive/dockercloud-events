FROM alpine:3.4
MAINTAINER Feng Honglin <hfeng@tutum.co>


ADD https://files.tutum.co/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD https://get.docker.com/builds/Linux/x86_64/docker-1.10.3 /usr/local/bin/docker
ADD dockercloud-events /events
RUN chmod +x /usr/local/bin/docker

ENV REPORT_INTERVAL=30 DOCKERCLOUD_AUTH=**None** EVENTS_API_URL=**None**

ENTRYPOINT ["/events"]
