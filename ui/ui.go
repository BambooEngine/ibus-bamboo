package ui

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>

extern int openGUI(guint flags, int mode, guint32 *s, int size, char *mtext, char *cfgtext);
*/
import "C"
import (
	"encoding/json"
	"ibus-bamboo/config"
	"io/ioutil"
	"os"
	"unsafe"
)

var engineName string

//export saveFlags
func saveFlags(flags C.guint) {
	var (
		cfg = config.LoadConfig(engineName)
	)
	config.SaveConfig(cfg, engineName)
	cfg.IBflags = uint(flags)
	config.SaveConfig(cfg, engineName)
}

//export saveConfigText
func saveConfigText(text *C.char) {
	var (
		cfgText = C.GoString(text)
		cfgFn   = config.GetConfigPath(engineName)
	)
	err := ioutil.WriteFile(cfgFn, []byte(cfgText), 0644)
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
	err := ioutil.WriteFile(macroFP, []byte(macroText), 0644)
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
func saveShortcuts(ptr *C.guint32, length int) {
	var (
		cfg = config.LoadConfig(engineName)
	)
	codes := makeSliceFromPtr(ptr, length)
	cfg.Shortcuts = codes
	config.SaveConfig(cfg, engineName)
}

func makeSliceFromPtr(ptr *C.guint32, size int) [10]uint32 {
	var out [10]uint32
	slice := (*[1 << 28]C.guint32)(unsafe.Pointer(ptr))[:size:size]
	for i, elem := range slice[:size] {
		out[i] = uint32(elem)
	}
	return out
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
	mText, err := ioutil.ReadFile(macroFilePath)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		panic(err)
	}
	os.Setenv("GTK_IM_MODULE", "gtk-im-context-simple")
	C.openGUI(C.guint(cfg.IBflags), C.int(cfg.DefaultInputMode), s, C.int(size), C.CString(string(mText)), C.CString(string(data)))
}
