# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.5.1-wheezy

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/chrismckenzie/dropship

RUN echo "deb http://ftp.debian.org/debian/ sid main" >> /etc/apt/sources.list && \
	apt-get update -qqy && \
	apt-get install --no-install-recommends -qqy \
	git libgit2-dev pkg-config apt-utils && \ 
	apt-get autoclean -qqy && \
  apt-get autoremove -qqy && \
  rm -rf /var/lib/apt/lists/*

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN cd /go/src/github.com/chrismckenzie/dropship && \
    go get ./... && \
    ln -s ./ui /go/bin && \
    go install

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/dropship

# Document that the service listens on port 8080.
EXPOSE 3000
