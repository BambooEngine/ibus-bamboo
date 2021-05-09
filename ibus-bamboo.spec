%define engine_name bamboo
%define ibus_dir     /usr/share/ibus
%define engine_dir   /usr/share/ibus-%{engine_name}
%define ibus_comp_dir /usr/share/ibus/component

Name: ibus-bamboo
Version: 0.7.1
Release: 1%{?dist}
Summary: A Vietnamese input method for IBus

License: GPLv3+
URL: https://github.com/BambooEngine/ibus-bamboo
Source0: %{name}-%{version}.tar.gz

BuildRequires: go, libX11-devel, libXtst-devel, gtk3-devel
Requires: ibus

%description
A Vietnamese IME for IBus using Bamboo Engine.
Bộ gõ tiếng Việt mã nguồn mở hỗ trợ hầu hết các bảng mã thông dụng, các kiểu gõ tiếng Việt phổ biến, bỏ dấu thông minh, kiểm tra chính tả, gõ tắt,...

%global debug_package %{nil}
%prep
%setup

%build
make build

%install
make DESTDIR=%{buildroot} install

%files
%defattr(-,root,root)
%doc README.md
%license LICENSE
%dir %{ibus_dir}
%dir %{ibus_comp_dir}
%dir %{engine_dir}
%{engine_dir}/*
%{ibus_comp_dir}/%{engine_name}.xml
/usr/lib/ibus-engine-%{engine_name}

%clean
cd ..
rm -rf ibus-%{engine_name}-%{version}
rm -rf %{buildroot}

%changelog
* Wed Aug 14 2019 LuongThanhLam <ltlam93@gmail.com> 0.5.3
- Initial RPM release
