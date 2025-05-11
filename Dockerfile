FROM golang:1.24-alpine

WORKDIR /mnt

# Install necessary build tools
RUN apk add --no-cache make git

# Create bin directory
RUN mkdir -p /mnt/bin

EXPOSE 9090

ENTRYPOINT ["/bin/sh"]