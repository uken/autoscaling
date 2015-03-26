FROM ubuntu:trusty
MAINTAINER Uken Ops "ops@uken.com"

ENV DEBIAN_FRONTEND noninteractive

ADD ./sources.list /etc/apt/sources.list

RUN rm -Rfv /var/lib/apt/lists/*

ADD ./packages.txt /tmp/packages.txt
RUN apt-get -y update && xargs apt-get install --force-yes -y < /tmp/packages.txt
RUN rm -f /tmp/packages.txt

RUN adduser --disabled-password --home=/app --gecos 'SheepIt' app

ADD https://godist.herokuapp.com/projects/ddollar/forego/releases/current/linux-amd64/forego /usr/bin/forego
RUN chmod 655 /usr/bin/forego

ENTRYPOINT ["/bin/bash", "-l", "-c"]
