# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/chrismckenzie/dropship

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN cd /go/src/github.com/chrismckenzie/dropship && \
    go get ./... && \
    go install

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/dropship

# Document that the service listens on port 8080.
EXPOSE 3000
