FROM postgres:9.6.18-alpine

LABEL MAINTAINER="Deepak Sharma <deepshar@redhat.com>"

COPY run-pgbouncer.sh /docker-entrypoint-initdb.d/
COPY health-check-probe.sh /usr/bin/

EXPOSE 5432