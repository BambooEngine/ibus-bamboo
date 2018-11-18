IBus Bamboo - Bộ gõ tiếng Việt cho Linux
===================================
[![Build Status](https://travis-ci.com/BambooEngine/ibus-bamboo.svg?branch=master)](https://travis-ci.com/BambooEngine/ibus-bamboo)
[![GitHub release](https://img.shields.io/github/release/BambooEngine/ibus-bamboo.svg)](https://github.com/BambooEngine/ibus-bamboo/releases/latest)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://opensource.org/licenses/GPL-3.0)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/BambooEngine/ibus-bamboo)

### Sơ lược tính năng
* Hỗ trợ tất cả các bảng mã phổ biến:
  * Unicode, TCVN (ABC)
  * VIQR, VNI, VPS, VISCII, BK HCM1, BK HCM2,…
  * Unicode UTF-8, Unicode NCR - for Web editors.
* Nhiều kiểu gõ:
  * Simple Telex, Telex 2, Telex 3, Telex + VNI + VIQR
  * VNI, VIQR, Microsoft layout
* Nhiều chế độ gõ:
  * Kiểm tra chính tả (tự động khôi phục tiếng anh với từ gõ sai)
  * Dấu thanh chuẩn và dấu thanh kiểu mới
  * Bỏ dấu tự do

Cài đặt và cấu hình
------------------

### Cài đặt (Ubuntu)

```sh
sudo add-apt-repository ppa:bamboo-engine/ibus-bamboo
sudo apt-get update
sudo apt-get install ibus-bamboo
ibus restart
```

*Hướng dẫn cài đặt từ mã nguồn: [wiki](https://github.com/BambooEngine/ibus-bamboo/wiki/H%C6%B0%E1%BB%9Bng-d%E1%BA%ABn-build-t%E1%BB%AB-source)*

### Gỡ bỏ
```
sudo apt remove ibus-bamboo
ibus restart
```

Góp ý và báo lỗi
--------------
https://github.com/BambooEngine/ibus-bamboo/issues

Giấy phép
-------
ibus-bamboo là bộ gõ được fork từ dự án [ibus-teni](https://github.com/teni-ime/ibus-teni) của tác giả Nguyễn Công Hoàng. Xem tệp AUTHORS để biết thêm chi tiết.

ibus-bamboo là phần mềm tự do nguồn mở. Toàn bộ mã nguồn của ibus-bamboo và bamboo-core đều được phát hành dưới các quy định ghi trong Giấy phép Công cộng GNU (GNU General Public License v3.0).
