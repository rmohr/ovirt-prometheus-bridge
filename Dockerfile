FROM sdurrheimer/alpine-glibc

MAINTAINER Roman Mohr <rmohr@redhat.com>

COPY ovirt-prometheus-bridge / 

ENTRYPOINT ["/ovirt-prometheus-bridge"]
