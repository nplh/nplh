# Maintainer: Caleb Eby <caleb.eby01@gmail.com>

pkgname=nplh
pkgver=0.0.0
pkgrel=1
pkgdesc="A quick dotfile linker"
arch=('x86_64' 'i686')
url="http://SERVER/$pkgname/"
license=('GPL3')
makedepends=('go')
options=('!strip' '!emptydirs')
source=("./nplh.go")
sha256sums=('7704418f141a35e829bdb78034da31dc9896655324cfab5e3af118d6bfd72fba')

build() {
  go build -o "$pkgname"
}

package() {
  install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
}

# vim:set ts=2 sw=2 et:
