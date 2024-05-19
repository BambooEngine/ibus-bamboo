package ui

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>

extern int openGUI(guint64 flags, int mode, guint32 *s, int size, char *mtext, char *cfgtext);
*/
import "C"
import (
	"encoding/json"
	"ibus-bamboo/config"
	"os"
	"unsafe"
)

var engineName string

//export fixFB
func fixFB(active int) {
	var (
		cfg = config.LoadConfig(engineName)
	)
	config.SaveConfig(cfg, engineName)
	if active == 1 {
		cfg.IBflags |= config.IBworkaroundForFBMessenger
	} else {
		cfg.IBflags &= ^config.IBworkaroundForFBMessenger
	}
	config.SaveConfig(cfg, engineName)
}

//export saveConfigText
func saveConfigText(text *C.char) {
	var (
		cfgText = C.GoString(text)
		cfgFn   = config.GetConfigPath(engineName)
	)
	err := os.WriteFile(cfgFn, []byte(cfgText), 0644)
	if err != nil {
		panic(err)
	}
}

//export saveMacroText
func saveMacroText(text *C.char) {
	var (
		macroText = C.GoString(text)
		macroFP   = config.GetMacroPath(engineName)
	)
	err := os.WriteFile(macroFP, []byte(macroText), 0644)
	if err != nil {
		panic(err)
	}
}

//export saveInputMode
func saveInputMode(mode int) {
	var (
		cfg = config.LoadConfig(engineName)
	)
	cfg.DefaultInputMode = mode
	config.SaveConfig(cfg, engineName)
}

//export saveShortcuts
func saveShortcuts(ptr *C.int, length int) {
	var (
		cfg = config.LoadConfig(engineName)
	)
	codes := makeSliceFromPtr(ptr, length)
	cfg.Shortcuts = codes
	config.SaveConfig(cfg, engineName)
}

func makeSliceFromPtr(ptr *C.int, length int) [10]uint32 {
	slice := unsafe.Slice(ptr, length)
	var ret [10]uint32
	for i, elem := range slice {
		ret[i] = uint32(elem)
	}
	return ret
}

func OpenGUI(engName string) {
	engineName = engName
	var (
		cfg           = config.LoadConfig(engineName)
		shortcuts     = cfg.Shortcuts[:]
		s             = (*C.guint32)(&shortcuts[0])
		size          = len(shortcuts)
		macroFilePath = config.GetMacroPath(engineName)
	)
	mText, err := os.ReadFile(macroFilePath)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		panic(err)
	}
	os.Setenv("GTK_IM_MODULE", "gtk-im-context-simple")
	C.openGUI(C.guint64(cfg.IBflags), C.int(cfg.DefaultInputMode), s, C.int(size), C.CString(string(mText)), C.CString(string(data)))
}
