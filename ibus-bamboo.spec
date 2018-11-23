#
# Bamboo - A Vietnamese Input method editor
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

#release info -----------------------------------------------------------------


%define engine_name  bamboo
%define package_name ibus-%{engine_name}
%define version      0.2.0


#install directories ----------------------------------------------------------
%define engine_dir   /usr/share/%{package_name}
%define ibus_dir     /usr/share/ibus
%define ibus_cpn_dir /usr/share/ibus/component
%define usr_lib_dir  /usr/lib


#package info -----------------------------------------------------------------
Name:           ibus-%{engine_name}
Version:        %{version}
Release:        1
Summary:        A Vietnamese IME for IBus
License:        GPL-3.0
Group:          System/Localization
URL:            https://github.com/BambooEngine/ibus-bamboo
Packager:       Luong Thanh Lam <ltlam93@gmail.com>
BuildRequires:  go, libX11-devel
Requires:       ibus
Provides:       locale(ibus:vi)
Source0:        %{package_name}-%{version}.tar.gz

%description
A Vietnamese IME for IBus using BambooEngine
Bộ gõ tiếng Việt cho IBus sử dụng BambooEngine

%global debug_package %{nil}
%prep
%setup


%build
make build


%install
make DESTDIR=%{buildroot} install


%files
%defattr(-,root,root)
%doc README.md LICENSE MAINTAINERS
%dir %{ibus_dir}
%dir %{ibus_cpn_dir}
%dir %{engine_dir}
%{engine_dir}/*
%{ibus_dir}/component/%{engine_name}.xml
%{usr_lib_dir}/ibus-engine-%{engine_name}


%clean
cd ..
rm -rf %{package_name}-%{version}
rm -rf %{buildroot}