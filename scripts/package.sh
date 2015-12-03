mkdir -p packaging/output

API=https://api.bintray.com
NOT_FOUND=404
SUCCESS=200
CREATED=201
PACKAGE_DESCRIPTOR=bintray-package.json

pkg=dropship

arch=$1
os=$2
version=$3
user=$4
key=$5
org=chrismckenzie
repo=deb
file=dropship-$version.$os-$arch.deb

function main() {
  set -x
  GOARCH=$arch GOOS=$os go build \
    -ldflags="-X github.com/ChrisMcKenzie/dropship/commands.version=${version}
    -X github.com/ChrisMcKenzie/dropship/commands.buildDate=`date -u +.%Y%m%d.%H%M%S`
    -X github.com/ChrisMcKenzie/dropship/commands.commitHash=`git rev-parse --short HEAD 2>/dev/null`" \
    -o packaging/root/usr/local/bin/dropship main.go

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

  set +x
  cp packaging/output/$file .

  git tag v${version}
  git push --tags

  init_curl

  release_package

  upload_content
}

function init_curl() {
  CURL="curl -u${user}:${key} -H Content-Type:application/json -H Accept:application/json"
}

function release_package() {
  echo "[DEBUG] Releasing package ${pkg}...\n"
  data="{
    \"name\": \"v${version}\",
    \"vcs_tag\": \"v${version}\",
  }"
  ${CURL} -X POST  -d  "${data}" ${API}/packages/${org}/${repo}/${pkg}/versions
  echo "\n"
}

function upload_content() {
  echo "[DEBUG] Uploading ${file}...\n"
  [ $(${CURL} -X PUT --write-out %{http_code} --silent --output /dev/null -T ${file} -H "X-Bintray-Package:${pkg}" -H "X-Bintray-Version: v${version}" -H "X-Bintray-Debian-Distribution: trusty" -H "X-Bintray-Debian-Component: main" -H "X-Bintray-Debian-Architecture: amd64" ${API}/content/${org}/${repo}/v${version}/pool/main/d/dropship/${file}) -eq ${CREATED} ]
  uploaded=$?
  echo "[DEBUG] DEB ${file} uploaded? y:1/N:0 ${uploaded}"
  ${CURL} -X POST --silent --output /dev/null  -d  "${data}" ${API}/content/${org}/${repo}/${pkg}/v${version}/publish
  return ${uploaded}
}

main "$@"
