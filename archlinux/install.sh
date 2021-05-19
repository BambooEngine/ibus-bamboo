#!/bin/bash
set_default() {
	ibus restart
	echo -n "Bạn có muốn đặt ibus-bamboo làm bộ gõ Tiếng Việt mặc định? [y/n]: "
	read choice
	case $choice in
		"y")
			env DCONF_PROFILE=ibus dconf write /desktop/ibus/general/preload-engines "['xkb:us::eng', 'Bamboo']" && gsettings set org.gnome.desktop.input-sources sources "[('xkb', 'us'), ('ibus', 'Bamboo')]"
			exit -1;;
		*) exit -1;;
	esac
}
echo "Chọn phiên bản muốn cài:"
echo "1. Bản release, cài đặt từ AUR (yay)"
echo "2. Bản release, cài đặt từ AUR (pamac)"
echo "3. Bản release, cài đặt từ Open Build Service (pacman)"
echo "4. Bản release, cài đặt từ mã nguồn"
echo "5. Bản git, cài đặt từ mã nguồn mới nhất lấy từ github"
echo "6. Thoát"
echo -n "Lựa chọn (1/2/3/4/5/6): "
read choice
case $choice in
	"1")
		if yay -S ibus-bamboo; then
			set_default
		fi
		exit -1;;
	"2")
		if pamac build ibus-bamboo; then
			set_default
		fi
		exit -1;;
	"3")
		echo -n Password:
		read -s szPassword
		echo $szPassword | sudo -S sh -c 'echo "[home_lamlng_Arch]" >> /etc/pacman.conf'
		echo $szPassword | sudo -S sh -c 'echo "Server = https://download.opensuse.org/repositories/home:/lamlng/Arch/\$arch" >> /etc/pacman.conf'
		key=$(curl -fsSL https://download.opensuse.org/repositories/home:lamlng/Arch/$(uname -m)/home_lamlng_Arch.key)
		fingerprint=$(gpg --quiet --with-colons --import-options show-only --import --fingerprint <<< "${key}" | awk -F: '$1 == "fpr" { print $10 }')
		echo $szPassword | sudo -S pacman-key --init
		echo $szPassword | sudo -S pacman-key --add - <<< "${key}"
		echo $szPassword | sudo -S pacman-key --lsign-key "${fingerprint}"
		if sudo -S pacman -Sy home_lamlng_Arch/ibus-bamboo; then
			set_default
		fi
		exit -1;;
	"4") VER="release";;
	"5") VER="git";;
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

set_default
