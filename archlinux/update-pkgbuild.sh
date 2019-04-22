#!/bin/bash

if [ ! -f bamboo.xml ]; then
	echo "Không tìm thấy file bamboo.xml, hãy chạy script này ở thư mục root của project"
	echo "./archlinux/update-pkgbuild.sh"
	exit -1;
fi

pkgver=`cat bamboo.xml |
	grep -o "<version>[0-9].[0-9].[0-9]</version>" |
	grep -o "[0-9].[0-9].[0-9]"`
echo $pkgver

echo "Sửa file PKGBUILD:"
echo "1. Chỉ sửa file PKGBUILD-git"
echo "2. Sửa cả file PKGBUILD-git và PKGBUILD-release (dành cho những lần release)"
echo "Chọn lựa (1/2)"
read opt

case $opt in
	"1")
		sed -i "s/pkgver=.*/pkgver=$pkgver/g" archlinux/PKGBUILD-git;;
	"2")
		sed -i "s/pkgver=.*/pkgver=$pkgver/g" archlinux/PKGBUILD-release
		sed -i "s/pkgver=.*/pkgver=$pkgver/g" archlinux/PKGBUILD-git;;
esac
