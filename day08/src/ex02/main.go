package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "window/application.m"
#include "window/window.m"
#import "window/application.h"
#import "window/window.h"
*/
import "C"

func main() {
	C.InitApplication()
	wnd := C.Window_Create(0, 0, 300, 200, C.CString("School 21"))
	C.Window_MakeKeyAndOrderFront(wnd)
	C.RunApplication()
}
