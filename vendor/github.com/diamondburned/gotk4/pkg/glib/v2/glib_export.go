// Code generated by girgen. DO NOT EDIT.

package glib

import (
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/core/gbox"
	"github.com/diamondburned/gotk4/pkg/core/gextras"
)

// #include <stdlib.h>
// #include <glib.h>
import "C"

//export _gotk4_glib2_CompareDataFunc
func _gotk4_glib2_CompareDataFunc(arg1 C.gconstpointer, arg2 C.gconstpointer, arg3 C.gpointer) (cret C.gint) {
	var fn CompareDataFunc
	{
		v := gbox.Get(uintptr(arg3))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(CompareDataFunc)
	}

	var _a unsafe.Pointer // out
	var _b unsafe.Pointer // out

	_a = (unsafe.Pointer)(unsafe.Pointer(arg1))
	_b = (unsafe.Pointer)(unsafe.Pointer(arg2))

	gint := fn(_a, _b)

	var _ int

	cret = C.gint(gint)

	return cret
}

//export _gotk4_glib2_EqualFuncFull
func _gotk4_glib2_EqualFuncFull(arg1 C.gconstpointer, arg2 C.gconstpointer, arg3 C.gpointer) (cret C.gboolean) {
	var fn EqualFuncFull
	{
		v := gbox.Get(uintptr(arg3))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(EqualFuncFull)
	}

	var _a unsafe.Pointer // out
	var _b unsafe.Pointer // out

	_a = (unsafe.Pointer)(unsafe.Pointer(arg1))
	_b = (unsafe.Pointer)(unsafe.Pointer(arg2))

	ok := fn(_a, _b)

	var _ bool

	if ok {
		cret = C.TRUE
	}

	return cret
}

//export _gotk4_glib2_Func
func _gotk4_glib2_Func(arg1 C.gpointer, arg2 C.gpointer) {
	var fn Func
	{
		v := gbox.Get(uintptr(arg2))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(Func)
	}

	var _data unsafe.Pointer // out

	_data = (unsafe.Pointer)(unsafe.Pointer(arg1))

	fn(_data)
}

//export _gotk4_glib2_HFunc
func _gotk4_glib2_HFunc(arg1 C.gpointer, arg2 C.gpointer, arg3 C.gpointer) {
	var fn HFunc
	{
		v := gbox.Get(uintptr(arg3))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(HFunc)
	}

	var _key unsafe.Pointer   // out
	var _value unsafe.Pointer // out

	_key = (unsafe.Pointer)(unsafe.Pointer(arg1))
	_value = (unsafe.Pointer)(unsafe.Pointer(arg2))

	fn(_key, _value)
}

//export _gotk4_glib2_HRFunc
func _gotk4_glib2_HRFunc(arg1 C.gpointer, arg2 C.gpointer, arg3 C.gpointer) (cret C.gboolean) {
	var fn HRFunc
	{
		v := gbox.Get(uintptr(arg3))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(HRFunc)
	}

	var _key unsafe.Pointer   // out
	var _value unsafe.Pointer // out

	_key = (unsafe.Pointer)(unsafe.Pointer(arg1))
	_value = (unsafe.Pointer)(unsafe.Pointer(arg2))

	ok := fn(_key, _value)

	var _ bool

	if ok {
		cret = C.TRUE
	}

	return cret
}

//export _gotk4_glib2_LogFunc
func _gotk4_glib2_LogFunc(arg1 *C.gchar, arg2 C.GLogLevelFlags, arg3 *C.gchar, arg4 C.gpointer) {
	var fn LogFunc
	{
		v := gbox.Get(uintptr(arg4))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(LogFunc)
	}

	var _logDomain string       // out
	var _logLevel LogLevelFlags // out
	var _message string         // out

	_logDomain = C.GoString((*C.gchar)(unsafe.Pointer(arg1)))
	_logLevel = LogLevelFlags(arg2)
	_message = C.GoString((*C.gchar)(unsafe.Pointer(arg3)))

	fn(_logDomain, _logLevel, _message)
}

//export _gotk4_glib2_LogWriterFunc
func _gotk4_glib2_LogWriterFunc(arg1 C.GLogLevelFlags, arg2 *C.GLogField, arg3 C.gsize, arg4 C.gpointer) (cret C.GLogWriterOutput) {
	var fn LogWriterFunc
	{
		v := gbox.Get(uintptr(arg4))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(LogWriterFunc)
	}

	var _logLevel LogLevelFlags // out
	var _fields []LogField      // out

	_logLevel = LogLevelFlags(arg1)
	{
		src := unsafe.Slice((*C.GLogField)(arg2), arg3)
		_fields = make([]LogField, arg3)
		for i := 0; i < int(arg3); i++ {
			_fields[i] = *(*LogField)(gextras.NewStructNative(unsafe.Pointer((&src[i]))))
		}
	}

	logWriterOutput := fn(_logLevel, _fields)

	var _ LogWriterOutput

	cret = C.GLogWriterOutput(logWriterOutput)

	return cret
}

//export _gotk4_glib2_SourceFunc
func _gotk4_glib2_SourceFunc(arg1 C.gpointer) (cret C.gboolean) {
	var fn SourceFunc
	{
		v := gbox.Get(uintptr(arg1))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(SourceFunc)
	}

	ok := fn()

	var _ bool

	if ok {
		cret = C.TRUE
	}

	return cret
}

//export _gotk4_glib2_SourceOnceFunc
func _gotk4_glib2_SourceOnceFunc(arg1 C.gpointer) {
	var fn SourceOnceFunc
	{
		v := gbox.Get(uintptr(arg1))
		if v == nil {
			panic(`callback not found`)
		}
		fn = v.(SourceOnceFunc)
	}

	fn()
}