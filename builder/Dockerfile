FROM golang:stretch

RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
    apt-get install cross-gcc-dev -y && \
    dpkg --add-architecture armhf && \
    dpkg --add-architecture armel && \
    apt-get update && \
    apt-get install -y crossbuild-essential-armel crossbuild-essential-armhf

WORKDIR /workspace
ENTRYPOINT ["/usr/bin/make", "build/covid-alert-arm"]
