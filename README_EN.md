IBus Bamboo - An open source Vietnamese IME for IBus using Bamboo Engine
===================================
[![GitHub release](https://img.shields.io/github/release/BambooEngine/ibus-bamboo.svg)](https://github.com/BambooEngine/ibus-bamboo/releases/latest)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://opensource.org/licenses/GPL-3.0)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/BambooEngine/ibus-bamboo)

IBus Bamboo is a Vietnamese input method engine for IBus that translates key strokes into Vietnamese characters. For example, when you type `a`, an `a` will appear on the screen, but if one more `a` is typed, IBus Bamboo will replace the first `a` with the letter `â` according to [Telex](https://en.wikipedia.org/wiki/Telex_(input_method)) typing.

   ![ibus-bamboo](https://github.com/BambooEngine/ibus-bamboo/raw/gh-resources/demo.gif)

## Getting Started

- [Features](#features)
- [Installation](#installation)
	- [Ubuntu, Mint and derivatives](#ubuntu-and-derivatives)
	- [Arch Linux and derivatives](#arch-linux-and-derivatives)
	- [Void Linux](#void-linux)
	- [Install from OpenBuildService](#install-from-openbuildservice)
	- [Install from source](https://github.com/BambooEngine/ibus-bamboo/wiki/H%C6%B0%E1%BB%9Bng-d%E1%BA%ABn-c%C3%A0i-%C4%91%E1%BA%B7t-t%E1%BB%AB-m%C3%A3-ngu%E1%BB%93n)
- [Usage](#usage)
- [Bug reports](#bug-reports)
- [License](#license)

## Features
* Support many Vietnamese character sets/encodings:
  * Unicode, TCVN (ABC)
  * VIQR, VNI, VPS, VISCII, BK HCM1, BK HCM2,…
  * Unicode UTF-8, Unicode NCR - for Web editors.
* All popular typing methods:
  * Telex, Telex W, Telex 2, Telex + VNI + VIQR
  * VNI, VIQR, Microsoft layout
* Using shortcut <kbd>Shift</kbd>+<kbd>~</kbd> to switch between typing modes for an application or add it to the exclusion list:
  	* Pre-edit (default)
  	* Surrounding text, IBus ForwardKeyEvent,...
* Other useful futures, easy to use:
  * Spelling check (using dictionary/rules)
  * Use oà, uý (instead of òa, úy)
  * Free tone making, macro,...
  * 2666 emojis from [emojiOne](https://github.com/joypixels/emojione)

## Installation
### Ubuntu and derivatives

```sh
sudo add-apt-repository ppa:bamboo-engine/ibus-bamboo
sudo apt-get update
sudo apt-get install ibus-bamboo
ibus restart
# Make ibus-bamboo your default input method, this will remove other existing input layouts
env DCONF_PROFILE=ibus dconf write /desktop/ibus/general/preload-engines "['xkb:us::eng', 'Bamboo']" && gsettings set org.gnome.desktop.input-sources sources "[('xkb', 'us'), ('ibus', 'Bamboo')]"
```

### Arch Linux and derivatives
```
bash -c "$(curl -fsSL https://raw.githubusercontent.com/BambooEngine/ibus-bamboo/master/archlinux/install.sh)"
```

### NixOS
`ibus-bamboo` is available on the main Nixpkgs repo. Make sure your NixOS configuration must contain this code to install it.

```nix
{
 i18n.inputMethod = {
  enabled = "ibus";
  ibus.engines = with pkgs.ibus-engines; [
    bamboo
  ];
 };
}
```

### Void Linux
`ibus-bamboo` is available on the main Void Linux repo. You can install it directly.

```sh
sudo xbps-install -S ibus-bamboo
```

### Install from OpenBuildService
[![OpenBuildService](https://github.com/BambooEngine/ibus-bamboo/raw/gh-resources/obs.png)](https://software.opensuse.org//download.html?project=home%3Alamlng&package=ibus-bamboo)

## Usage
The difference between `ibus-bamboo` and other input methods is that `ibus-bamboo` provides different typing modes (1 underlined and 5 non-underlined typing modes - don't confuse **typing mode** with **typing method**, typing methods are `telex`, `vni`, ...).

To switch between typing modes, simply click on an input box (a box to enter text), then press the combination <kbd>Shift</kbd>+<kbd>~</kbd>, a table with the available typing modes will appear, you just need to press the corresponding number key to select.

**Note:**
 - An app may work well with one typing mode while not working well with another.
 - Typing modes are saved separately for each software (`firefox` is probably using mode 5, while `libreoffice` is using mode 2).
 - You can use `Add to the exclusion list` mode to not type Vietnamese in a certain program.
 - To type the character `~`, press the combination <kbd>Shift</kbd>+<kbd>~</kbd> twice.

## Bug reports
Before submitting a question or bug report, please ensure you have read through [these common issues](https://github.com/BambooEngine/ibus-bamboo/wiki/C%C3%A1c-v%E1%BA%A5n-%C4%91%E1%BB%81-th%C6%B0%E1%BB%9Dng-g%E1%BA%B7p) and see if you can resolve the problem on your own. If you still encounter issues after trying these steps, or you don't see something similar to your issue listed, please submit a bug report in the [Bamboo issue tracker](https://github.com/BambooEngine/ibus-bamboo/issues)

## License
IBus Bamboo is released under the GNU General Public License v3.0
