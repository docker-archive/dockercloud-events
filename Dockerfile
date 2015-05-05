FROM alpine
MAINTAINER Feng Honglin <hfeng@tutum.co>

ADD run.sh /run.sh
RUN chmod +x /run.sh
ADD container-events /container-events
ADD https://files.tutum.co/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV TUTUM_HOST https://dashboard.tutum.co/
ENV DOCKER_HOST unix:///var/run/docker.sock
ENV RESTART_INTERVAL 3600
ENV TUTUM_AUTH **None**
ENV NODE_UUID **None**

CMD ["/run.sh"]
