mkdir -p packaging/output

set -x
GOARCH=$1 GOOS=$2 go build -o packaging/root/usr/local/bin/dropship main.go

fpm -s dir -t deb -n dropship -v "$1" -p packaging/output/dropship.$2-$1.deb \
  --deb-priority optional --category admin \
  --force \
  --after-install packaging/scripts/postinstall.deb \
  --before-remove packaging/scripts/prerm.deb \
  --after-remove packaging/scripts/postrm.deb \
  --deb-compression bzip2 \
  --url https://github.com/chrismckenzie/dropship \
  --description "Dropship automatically keeps you software up to date" \
  -m "Chris McKenzie <chris@chrismckenzie.io>" \
  -a $1 \
  --config-files /etc/dropship.d/dropship.hcl \
  packaging/root/=/

cp packaging/output/dropship.$2-$1.deb .
