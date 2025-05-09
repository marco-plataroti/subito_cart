FROM golang:1.22-alpine

WORKDIR /mnt

EXPOSE 9090

ENTRYPOINT ["/bin/sh"]