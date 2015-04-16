# Maintainer:  <s-kostyaev@ngs>
pkgname=webtop-container-git
pkgver=0.1.0
pkgrel=1
pkgdesc="container part of web-based top for cgroup"
arch=('i686' 'x86_64')
url="https://github.com/s-kostyaev/webtop-container"
license=('unknown')
depends=('git')
makedepends=('go')
branch='dev'
source=("${pkgname}::git+https://github.com/s-kostyaev/webtop-container#branch=${branch}")
md5sums=('SKIP')
build(){
  cd ${srcdir}/${pkgname}
  deps=`go list -f '{{join .Deps "\n"}}' |  xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'`
  for dep in $deps; do go get $dep; done
  go build -o webtop-container
}
package(){
  install -D -m 755 ${srcdir}/${pkgname}/webtop-container ${pkgdir}/usr/bin/webtop-container
  install -D -m 644 ${srcdir}/${pkgname}/webtop-container.service ${pkgdir}/usr/lib/systemd/system/webtop-container.service
  install -d -m 644 ${srcdir}/${pkgname}/static ${pkgdir}/usr/share/webtop
}

