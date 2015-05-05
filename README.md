tutum/utils:container-events
============================
   
    docker run \
      -d \
      -v /var/run:/var/run:rw \
      -v /usr/lib/tutum/docker:/docker:r \
      -e TUTUM_HOST="https://dashboard.tutum.co/" \
      -e DOCKER_HOST="unix:///var/run/docker.sock" \
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
