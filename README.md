IBus Bamboo - Bộ gõ tiếng Việt cho Linux
===================================
[![GitHub release](https://img.shields.io/github/release/BambooEngine/ibus-bamboo.svg)](https://github.com/BambooEngine/ibus-bamboo/releases/latest)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://opensource.org/licenses/GPL-3.0)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/BambooEngine/ibus-bamboo)

## Mục lục

- [Sơ lược tính năng](#sơ-lược-tính-năng)
- [Hướng dẫn cài đặt](#hướng-dẫn-cài-đặt)
	- [Dành cho Ubuntu, Debian và các distro tương tự](#ubuntu-debian-và-các-distro-tương-tự)
	- [Dành cho Arch Linux và các distro tương tự](#arch-linux-và-các-distro-tương-tự)
	- [Cài đặt từ mã nguồn](https://github.com/BambooEngine/ibus-bamboo/wiki/H%C6%B0%E1%BB%9Bng-d%E1%BA%ABn-c%C3%A0i-%C4%91%E1%BA%B7t-t%E1%BB%AB-m%C3%A3-ngu%E1%BB%93n)
- [Hướng dẫn sử dụng](#hướng-dẫn-sử-dụng)
- [Báo lỗi](#báo-lỗi)
- [Giấy phép](#giấy-phép)

## Sơ lược tính năng
* Hỗ trợ tất cả các bảng mã phổ biến:
  * Unicode, TCVN (ABC)
  * VIQR, VNI, VPS, VISCII, BK HCM1, BK HCM2,…
  * Unicode UTF-8, Unicode NCR - for Web editors.
* Nhiều kiểu gõ:
  * Telex, Telex 2, Telex 3, Telex + VNI + VIQR
  * VNI, VIQR, Microsoft layout
* Nhiều tính năng hữu ích, dễ dàng tùy chỉnh:
  * Kiểm tra chính tả (sử dụng từ điển/luật ghép vần)
  * Dấu thanh chuẩn và dấu thanh kiểu mới
  * Bỏ dấu tự do, Gõ tắt, Emoji,...
* Sử dụng phím tắt `<Shift>`+`~` để loại trừ ứng dụng không dùng bộ gõ, chuyển qua lại giữa các chế độ gõ:
  	* Pre-edit (default)
  	* Surrounding text, IBus ForwardKeyEvent, XTestFakeKeyEvent,...
   ![ibus-bamboo](https://github.com/BambooEngine/ibus-bamboo/raw/gh-resources/demo.gif)

## Hướng dẫn cài đặt
### Ubuntu, Debian và các distro tương tự

```sh
sudo add-apt-repository ppa:bamboo-engine/ibus-bamboo
sudo apt-get update
sudo apt-get install ibus-bamboo -y
ibus restart
```

### Arch Linux và các distro tương tự
Với Arch Linux, bạn có thể cài đặt bằng cách chạy lệnh sau:
```sh
wget https://raw.githubusercontent.com/BambooEngine/ibus-bamboo/master/archlinux/install.sh
chmod +x install.sh
./install.sh
```

Sau đó script sẽ cài đặt `ibus-bamboo` cho bạn.

## Hướng dẫn sử dụng
Điểm khác biệt giữa `ibus-bamboo` và các bộ gõ khác là `ibus-bamboo` cung cấp nhiều chế độ gõ khác nhau (1 chế độ gõ có gạch chân và 5 chế độ gõ không gạch chân; tránh nhầm lẫn **chế độ gõ** với **kiểu gõ**, các kiểu gõ bao gồm `telex`, `vni`, ...). Để chuyển đổi giữa các chế độ gõ, chỉ cần nhấn vào một khung nhập liệu (một cái hộp để nhập văn bản) nào đó, sau đó nhấn tổ hợp `Shift + ~`, một bảng với những chế độ gõ hiện có sẽ xuất hiện, bạn chỉ cần nhấn phím số tương ứng để lựa chọn.

**Một số lưu ý:**
- Các chế độ gõ được lưu riêng biệt cho mỗi phần mềm (`firefox` có thể đang dùng chế độ 3, trong khi `libreoffice` thì lại dùng chế độ 2).
- Nếu một phần mềm chưa được đặt chế độ gõ thì nó sẽ dùng chế độ gõ mặc định.
- Bạn có thể dùng chế độ `thêm vào danh sách loại trừ` để không gõ tiếng Việt trong một chương trình nào đó.
- Để gõ ký tự `~` hãy nhấn tổ hợp `Shift + ~` 2 lần.

## Báo lỗi
Trước khi báo lỗi vui lòng đọc [những vấn đề thường gặp](https://github.com/BambooEngine/ibus-bamboo/wiki/C%C3%A1c-v%E1%BA%A5n-%C4%91%E1%BB%81-th%C6%B0%E1%BB%9Dng-g%E1%BA%B7p) và tìm vấn đề của mình ở trong đó.

Nếu trang phía trên không giải quyết vấn đề của bạn, vui lòng [báo lỗi tại đây](https://github.com/BambooEngine/ibus-bamboo/issues)

## Giấy phép
ibus-bamboo là bộ gõ được fork từ dự án [ibus-teni](https://github.com/teni-ime/ibus-teni), sử dụng Bamboo Engine để xử lý tiếng Việt thay cho Teni Engine.

ibus-bamboo là phần mềm tự do nguồn mở. Toàn bộ mã nguồn của ibus-bamboo được phát hành dưới các quy định ghi trong Giấy phép Công cộng GNU (GNU General Public License v3.0).
