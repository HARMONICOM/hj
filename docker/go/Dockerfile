FROM golang:1.25.0-trixie

RUN apt-get update && apt-get install -y \
    jq \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app
