FROM alpine
MAINTAINER Feng Honglin <hfeng@tutum.co>

ADD dockercloud-events /events
RUN apk update && apk add ca-certificates
ENV REPORT_INTERVAL=30 DOCKERCLOUD_AUTH=**None** EVENTS_API_URL=**None**

ENTRYPOINT ["/events"]
