# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.5.1-wheezy

ENV GITHUB_CLIENT_ID "<your-github-client-id>"
ENV GITHUB_CLIENT_SECRET "<your-github-client-secret>"
ENV APP_URL "http://localhost:3000"

WORKDIR /go/src/github.com/chrismckenzie/dropship 

COPY . /go/src/github.com/chrismckenzie/dropship

RUN echo "deb http://ftp.debian.org/debian/ sid main" >> /etc/apt/sources.list
RUN	apt-get update --no-upgrade --no-install-recommends && \
    apt-get install --no-upgrade --no-install-recommends \
    -y libgit2-dev pkg-config

RUN go get ./... && \
    go build -o bin/dropship

ENTRYPOINT ./bin/dropship

EXPOSE 3000
