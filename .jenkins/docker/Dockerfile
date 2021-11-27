FROM debian:bullseye-slim

ENV GOPROXY "direct"

RUN apt-get update
RUN apt-get install -y \
    curl \
    git \
    golang

RUN useradd -u 1001 jenkins -m
RUN mkdir /home/jenkins/go

USER jenkins
