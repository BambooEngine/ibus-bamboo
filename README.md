IBus Bamboo - Bộ gõ tiếng Việt cho Linux
===================================
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
  * Fast committing

Cài đặt và cấu hình
------------------

### Cài đặt (Ubuntu)

```sh
sudo add-apt-repository ppa:bamboo-engine/ibus-bamboo
sudo apt-get update
sudo apt-get install ibus-bamboo
ibus restart
```

**Lệnh bên dưới cho phép đọc event chuột, không bắt buộc nhưng cần để ibus-bamboo hoạt động tốt**
```sh
sudo usermod -a -G input $USER
```

### Gỡ bỏ
```
sudo apt remove ibus-bamboo
ibus restart
```

Sử dụng
-------------
* Dùng phím tắt mặc định của IBus (thường là ⊞Win+Space) để chuyển đổi giữa các bộ gõ
* ibus-bamboo dùng pre-edit để xử lý phím gõ, mặc định sẽ có gạch chân chữ khi đang gõ
* **Khi pre-edit chưa kết thúc mà nhấn chuột vào chỗ khác thì có 3 khả năng xảy ra tùy chương trình**
    * **Chữ đang gõ bị mất**
    * **Chữ đang gõ được commit vào vị trí mới con trỏ**
    * **Chữ đang gõ được commit vào vị trí cũ**
* **Vì vậy đừng quên commit: khi gõ chỉ một chữ, hoặc chữ cuối câu, hoặc sửa chữ giữa câu: nhấn phím *Ctrl* để commit.**


Góp ý và báo lỗi
--------------
https://github.com/BambooEngine/ibus-bamboo/issues

Giấy phép
-------
ibus-bamboo là phần mềm tự do nguồn mở.

Toàn bộ mã nguồn của ibus-bamboo và bamboo-core đều được phát hành dưới các quy định ghi trong Giấy phép Công cộng GNU (GNU General Public License v3.0).