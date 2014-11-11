tutum/utils:container-events
============================
   
```
    docker run \
      -d \
      -v /var/run:/var/run:rw \
      -v /etc/tutum:/etc/tutum:r \
      -e TUTUM_HOST="https://dashboard.tutum.co/" \
      -e DOCKER_HOST="unix:///var/run/docker.sock" \
      tutum/utils:container-events
```

**Arguments**

```
    Environment variable: TUTUM_HOST      tutum host, "https://dashboard.tutum.co/" by default
    Environment varialbe: DOCKER_HOST     docker host, "unix:///var/run/docker.sock" by default
```
