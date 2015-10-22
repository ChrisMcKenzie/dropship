
mkdir -p packaging/output

set -x
GOARCH=amd64 GOOS=linux go build -o packaging/root/usr/local/bin/dropship main.go

fpm -s dir -t deb -n dropship -v "$1" -p packaging/output/dropship.deb \
  --deb-priority optional --category admin \
  --force \
  --after-install packaging/scripts/postinstall.deb \
  --before-remove packaging/scripts/prerm.deb \
  --after-remove packaging/scripts/postrm.deb \
  --deb-compression bzip2 \
  --url https://github.com/chrismckenzie/dropship \
  --description "Dropship automatically keeps you software up to date" \
  -m "Chris McKenzie <chris@chrismckenzie.io>" \
  -a amd64 \
  --config-files /etc/dropship.d/dropship.hcl \
  packaging/root/=/

cp packaging/output/dropship.deb .
