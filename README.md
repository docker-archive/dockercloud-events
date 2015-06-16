tutum/utils:container-events
============================
   
    docker run \
      -d \
      -v /var/run:/var/run:rw \
      -e TUTUM_HOST="https://dashboard.tutum.co/" \
      -e DOCKER_HOST="unix:///var/run/docker.sock" \
      -e REPORT_INTERVAL=30
      -e TUTUM_AUTH=xxxxxxxxxx \
      -e NODE_UUID=xxxxxxxxx \
      -e SENTRY_DSN=xxxxxxxx \
      tutum/events


**Arguments**

    TUTUM_HOST          tutum host, "https://dashboard.tutum.co/" by default
    DOCKER_HOST         docker host, "unix:///var/run/docker.sock" by default
    TUTUM_AUTH          tutum auth
    NODE_UUID           node uuid
    SENTRY_DSN          sentry dsn
    REPORT_INTERVAL     interval to report autorestarted container events to tutum
