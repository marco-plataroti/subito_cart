FROM golang:1.24-alpine

WORKDIR /mnt

EXPOSE 9090

ENTRYPOINT ["/bin/sh"]