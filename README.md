tutum/events
============
   
    docker run \
      -d \
      -v /var/run:/var/run:rw \
      -e TUTUM_AUTH=xxxxxxxxxx \
      -e NODE_UUID=xxxxxxxxx \
      [-e SENTRY_DSN=xxxxxxxx] \
      [-e REPORT_INTERVAL=30] \
      [-e TUTUM_HOST="https://dashboard.tutum.co/"] \
      [-e DOCKER_HOST="unix:///var/run/docker.sock"] \
      tutum/events


## Arguments


Key | Description
----|------------
TUTUM_AUTH | Tutum's API role `Authorization` header
NODE_UUID | Tutum's node UUID
TUTUM_HOST | (optional) Tutum API host
DOCKER_HOST | (optional) Docker host
SENTRY_DSN | (optional) Sentry DSN for bug reporting
REPORT_INTERVAL | (optional) Interval in seconds to report autorestarted container events
