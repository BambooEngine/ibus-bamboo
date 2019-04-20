#!/bin/bash
if [ -d ibus-bamboo ]; then
	echo "Tìm thấy thư mục tên ibus-bamboo, đổi tên thành ibus-bamboo-bak"
        mv ibus-bamboo ibus-bamboo-bak
fi

if [ -f ibus-bamboo ]; then
	echo "Tìm thấy file tên ibus-bamboo, đổi tên thành ibus-bamboo~"
        mv ibus-bamboo ibus-bamboo~
fi

echo "Chọn phiên bản muốn cài:"
echo "1. Bản git, mã nguồn mới nhất lấy từ github"
echo "2. Bản release"
echo "3. Thoát"
echo "Lựa chọn (1/2/3): "
read choice
case $choice in
	"1") VER="git";;
	"2") VER="release";;
	*) exit -1;;
esac

mkdir ibus-bamboo
cd ibus-bamboo
wget "https://raw.githubusercontent.com/BambooEngine/ibus-bamboo/master/archlinux/PKGBUILD-$VER" -O PKGBUILD
makepkg -si

cd ..
rm ibus-bamboo -rf
rm install.sh
