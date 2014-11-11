FROM tutum/curl:trusty
MAINTAINER Feng Honglin <hfeng@tutum.co>

ADD . /gopath/src/github.com/tutumcloud/tutum-docker-utils/container-events

RUN apt-get update -y && \
    apt-get install --no-install-recommends -y -q git && \
    mkdir /goroot && \
    curl -s https://storage.googleapis.com/golang/go1.3.linux-amd64.tar.gz | tar xzf - -C /goroot --strip-components=1 && \ 
    export GOROOT=/goroot && \
    export GOPATH=/gopath && \
    export PATH=$PATH:/goroot/bin && \
    go get github.com/tutumcloud/tutum-docker-utils/container-events && \
    cp /gopath/bin/* / && \
    rm -fr /goroot /gopath /var/lib/apt/lists && \
    apt-get autoremove -y git && \
    apt-get clean

ENV TUTUM_HOST https://dashboard.tutum.co/
ENV DOCKER_HOST unix:///var/run/docker.sock

ADD run.sh /run.sh
RUN chmod +x /run.sh

ENTRYPOINT ["/run.sh"]
