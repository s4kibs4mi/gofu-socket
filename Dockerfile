FROM golang:alpine
RUN apk add --no-cache --update git
ENV GOPATH=/go
COPY gofu-socket /usr/local/bin/
ENTRYPOINT ["gofu-socket"]
