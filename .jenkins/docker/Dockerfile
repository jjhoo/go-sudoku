FROM debian:bullseye-slim

ARG JENKINS_UID=1001
ENV DEBIAN_FRONTEND teletype
ENV GOPROXY "direct"

RUN apt-get update
RUN apt-get install -y \
    curl \
    git \
    golang

RUN useradd -u $JENKINS_UID jenkins -m

WORKDIR /home/jenkins
USER jenkins
RUN mkdir /home/jenkins/go
