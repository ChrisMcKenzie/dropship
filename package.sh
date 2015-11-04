mkdir -p packaging/output

arch=$1
os=$2
version=$3

set -x
GOARCH=$arch GOOS=$os go build -o packaging/root/usr/local/bin/dropship main.go

fpm -s dir -t deb -n dropship -v "$version" -p packaging/output/dropship-$version.$os-$arch.deb \
  --deb-priority optional --category admin \
  --force \
  --after-install packaging/scripts/postinstall.deb \
  --before-remove packaging/scripts/prerm.deb \
  --after-remove packaging/scripts/postrm.deb \
  --deb-compression bzip2 \
  --url https://github.com/chrismckenzie/dropship \
  --description "Dropship automatically keeps you software up to date" \
  -m "Chris McKenzie <chris@chrismckenzie.io>" \
  -a $arch \
  --config-files /etc/dropship.d/dropship.hcl \
  packaging/root/=/

cp packaging/output/dropship-$version.$os-$arch.deb .
