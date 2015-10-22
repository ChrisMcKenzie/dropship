GOOS=linux
GOARCH=amd64

mkdir -p packaging/output

set -x
go build -o packaging/root/usr/local/bin/dropship main.go

fpm -s dir -t deb -n "dropship" -a amd64 -v $1 \
  -p packaging/output/dropship.deb \
  -m "Chris McKenzie <chris@chrismckenzie.io>" \
  --description "Dropship code deployment agent"
  --force \
  --deb-compression bzip2 \
  packaging/root/=/

cp packaging/output/dropship.deb .
