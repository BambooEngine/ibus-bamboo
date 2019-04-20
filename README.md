IBus Bamboo - Bộ gõ tiếng Việt cho Linux
===================================
[![GitHub release](https://img.shields.io/github/release/BambooEngine/ibus-bamboo.svg)](https://github.com/BambooEngine/ibus-bamboo/releases/latest)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://opensource.org/licenses/GPL-3.0)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/BambooEngine/ibus-bamboo)

## Mục lục

- [Sơ lược tính năng](#sơ-lược-tính-năng)
- [Hướng dẫn cài đặt](#hướng-dẫn-cài-đặt)
	- [Dành cho Ubuntu, Debian và các distro tương tự](#ubuntu-debian-và-các-distro-tương-tự)
	- [Dành cho Archlinux và các distro tương tự](#archlinux-và-các-distro-tương-tự)
	- [Cài đặt từ mã nguồn](https://github.com/BambooEngine/ibus-bamboo/wiki/H%C6%B0%E1%BB%9Bng-d%E1%BA%ABn-c%C3%A0i-%C4%91%E1%BA%B7t-t%E1%BB%AB-m%C3%A3-ngu%E1%BB%93n)
- [Hướng dẫn sử dụng](#hướng-dẫn-sử-dụng)
	- [Cài đặt biến môi trường để sử dụng ibus (nên đọc)](#cài-đặt-biến-môi-trường)
	- [Chuyển đổi giữa các chế độ gõ](#chuyển-đổi-giữa-các-chế-độ-gõ)
- [Giấy phép](#giấy-phép)

## Sơ lược tính năng
* Hỗ trợ tất cả các bảng mã phổ biến:
  * Unicode, TCVN (ABC)
  * VIQR, VNI, VPS, VISCII, BK HCM1, BK HCM2,…
  * Unicode UTF-8, Unicode NCR - for Web editors.
* Nhiều kiểu gõ:
  * Telex, Telex 2, Telex 3, Telex + VNI + VIQR
  * VNI, VIQR, Microsoft layout
* Nhiều chế độ gõ:
  * Kiểm tra chính tả (sử dụng từ điển/luật ghép vần)
  * Dấu thanh chuẩn và dấu thanh kiểu mới
  * Bỏ dấu tự do, Gõ tắt (macro),...
* Sử dụng phím tắt `<Shift>`+`~` để loại trừ ứng dụng không dùng bộ gõ, chuyển qua lại giữa các chế độ gõ:
  	* Pre-edit (default)
  	* Surrounding text, IBus Forward key event, XTestFakeKeyEvent
   ![ibus-bamboo](https://github.com/BambooEngine/ibus-bamboo/raw/gh-resources/demo.gif)

## Hướng dẫn cài đặt
### Ubuntu, Debian và các distro tương tự

```sh
sudo add-apt-repository ppa:bamboo-engine/ibus-bamboo
sudo apt-get update
sudo apt-get install ibus-bamboo -y
ibus restart
```

### Archlinux và các distro tương tự
Với Archlinux, cách cài đặt giống như trên AUR. Đầu tiên các bạn tải file PKGBUILD tương ứng về máy. Có 2 phiên bản để cài đặt, bản git là bản dùng mã nguồn mới nhất từ master, bản còn lại là release:
```sh
mkdir ibus-bamboo
cd ibus-bamboo

# nếu muốn cài bản git
wget https://raw.githubusercontent.com/BambooEngine/ibus-bamboo/master/archlinux/PKGBUILD-git -O PKGBUILD
# nếu muốn cài bản release
wget https://raw.githubusercontent.com/BambooEngine/ibus-bamboo/master/archlinux/PKGBUILD-release -O PKGBUILD
```

Cuối cùng build gói và cài đặt
```sh
makepkg -si
```

## Hướng dẫn sử dụng
### Cài đặt biến môi trường
Việc cài đặt biến môi trường là để đảm bảo các phần mềm khác sẽ sử dụng ibus. Để cài đặt các bạn thêm những dòng sau vào trong file `~/.bashrc` và `~/.profile`

```sh
export GTK_IM_MODULE=ibus
export QT_IM_MODULE=ibus
export XMODIFIERS=@im=ibus
# Dành cho những phần mềm dựa trên qt4?
export QT4_IM_MODULE=ibus
# Dành cho những phần mềm dùng clutter (hình như chỉ có trên gnome)
export CLUTTER_IM_MODULE=ibus
```

Việc cài đặt trên chỉ có hiệu lực cho người dùng hiện tại, nếu muốn cài đặt cho toàn bộ hệ thống hãy để những dòng trên vào file `/etc/bash.bashrc` và `/etc/profile`.

**Lưu ý:** Nếu bạn dùng shell khác như `zsh` thì thay vì `.bashrc`, hãy thêm vào `.zshrc`. Tương tự với `fish` hay những shell khác.
*Tham khảo thêm tại [wiki ibus của Archlinux](https://wiki.archlinux.org/index.php/IBus#Initial_setup)*

### Chuyển đổi giữa các chế độ gõ
Tránh nhầm lẫn **chế độ gõ** với **kiểu gõ** (các kiểu gõ bao gồm `telex`, `vni`, ...). Để chuyển đổi giữa các chế độ gõ, chỉ cần nhấn vào một khung nhập liệu (một cái hộp để nhập văn bản) nào đó, sau đó nhấn tổ hợp `Shift + ~`, một bảng với những chế độ gõ hiện có sẽ xuất hiện, bạn chỉ cần nhấn phím số tương ứng để lựa chọn.

**Một số lưu ý:**
- Các chế độ gõ được lưu riêng biệt cho mỗi phần mềm (`firefox` có thể đang dùng chế độ 3, trong khi `libreoffice` thì lại dùng chế độ 2).
- Nếu một phần mềm chưa được đặt chế độ gõ thì nó sẽ dùng chế độ gõ mặc định.
- Bạn có thể dùng chế độ `thêm vào danh sách loại trừ` để không gõ tiếng Việt trong một chương trình nào đó.
- Để gõ ký tự `~` hãy nhấn tổ hợp `Shift + ~` 2 lần.

## Giấy phép
ibus-bamboo là bộ gõ được fork từ dự án [ibus-teni](https://github.com/teni-ime/ibus-teni), sử dụng Bamboo Engine để xử lý tiếng Việt thay cho Teni Engine.

ibus-bamboo là phần mềm tự do nguồn mở. Toàn bộ mã nguồn của ibus-bamboo được phát hành dưới các quy định ghi trong Giấy phép Công cộng GNU (GNU General Public License v3.0).
