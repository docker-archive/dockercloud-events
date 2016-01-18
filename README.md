tutum/events-daemon
===================

System container that forwards docker events to Docker Cloud's API. System containers are launched, configured and managed automatically on every node by Docker Cloud.


## Usage

    docker run \
      -d \
      -v /usr/lib/dockercloud/docker:/usr/bin/docker
      -v /var/run/docker.sock:/var/run/docker.sock:rw \
      -e EVENTS_API_URL=xxxxxxxx \
      -e DOCKERCLOUD_AUTH=xxxxxxxx \
      [-e REPORT_INTERVAL=30] \
      tutum/events-daemon


## Arguments

Key | Description
----|------------
DOCKERCLOUD_AUTH | Docker Cloud's API role `Authorization` header
EVENTS_API_URL  | The URL that docker events are POSTed to
REPORT_INTERVAL | (optional) Interval in seconds to report autorestarted container events
