#!/bin/bash
echo "Chọn phiên bản muốn cài:"
echo "1. Bản release, cài đặt từ Open Build Service"
echo "2. Bản release, cài đặt từ mã nguồn"
echo "3. Bản git, cài đặt từ mã nguồn mới nhất lấy từ github"
echo "4. Thoát"
echo -n "Lựa chọn (1/2/3/4): "
read choice
case $choice in
	"1")
		echo -n Password:
		read -s szPassword
		echo $szPassword | sudo -S sh -c 'echo "[home_lamlng_Arch]" >> /etc/pacman.conf'
		echo $szPassword | sudo -S sh -c 'echo "Server = https://download.opensuse.org/repositories/home:/lamlng/Arch/\$arch" >> /etc/pacman.conf'
		key=$(curl -fsSL https://download.opensuse.org/repositories/home:lamlng/Arch/$(uname -m)/home_lamlng_Arch.key)
		fingerprint=$(gpg --quiet --with-colons --import-options show-only --import --fingerprint <<< "${key}" | awk -F: '$1 == "fpr" { print $10 }')
		echo $szPassword | sudo -S pacman-key --init
		echo $szPassword | sudo -S pacman-key --add - <<< "${key}"
		echo $szPassword | sudo -S pacman-key --lsign-key "${fingerprint}"
		echo $szPassword | sudo -S pacman -Sy home_lamlng_Arch/ibus-bamboo
		exit -1;;
	"2") VER="release";;
	"3") VER="git";;
	*) exit -1;;
esac

if [ -d ibus-bamboo ]; then
	echo "Tìm thấy thư mục tên ibus-bamboo, đổi tên thành ibus-bamboo-bak"
        mv ibus-bamboo ibus-bamboo-bak
fi

if [ -f ibus-bamboo ]; then
	echo "Tìm thấy file tên ibus-bamboo, đổi tên thành ibus-bamboo~"
        mv ibus-bamboo ibus-bamboo~
fi

mkdir ibus-bamboo
cd ibus-bamboo
wget "https://raw.githubusercontent.com/BambooEngine/ibus-bamboo/master/archlinux/PKGBUILD-$VER" -O PKGBUILD
makepkg -si

cd ..
rm ibus-bamboo -rf
rm install.sh
