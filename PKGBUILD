#
# Bamboo - A Vietnamese Input method editor
# Copyright (C) 2018 Nguyen Cong Hoang <hoangnc.jp@gmail.com>
# Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
# Maintainer: Luong Thanh Lam <ltlam93@gmail.com>
pkgname=ibus-bamboo
pkgver=0.1.8
pkgrel=1
pkgdesc='A Vietnamese IME for IBus'
arch=(any)
license=(GPL3)
url="https://github.com/BambooEngine/ibus-bamboo"
depends=(ibus)
makedepends=('go' 'libx11')
source=($pkgname-$pkgver.tar.gz)
md5sums=('SKIP')
options=('!strip')

build() {
  cd "$pkgname-$pkgver"

  make
}


package() {
  cd "$pkgname-$pkgver"

  make DESTDIR="$pkgdir/" install
}
