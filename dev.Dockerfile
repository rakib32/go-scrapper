# DO NOT use in production!
# Dockerfile for local development

FROM golang:1.16

# install the watcher
RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /app
COPY ./ /app

ENTRYPOINT ["CompileDaemon", "-build=go build -mod=vendor", "-log-prefix=false", "-exclude-dir=.git"]
