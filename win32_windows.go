// this file was generated from win32_windows.ft using github.com/akovaski/templateParser
package gkdl

import (
	"syscall"
	"unsafe"
)

func boolToUintptr(a bool) uintptr {
	if a {
		return 1
	}
	return 0
}

func uintptrToBool(a uintptr) bool {
	return a != 0
}

func stringToUintptr(s string) uintptr {
	if s == "" {
		return 0
	}
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

// load DLLS
var (
	user32 = syscall.MustLoadDLL("user32.dll")

	kernel32 = syscall.MustLoadDLL("kernel32.dll")

	gdi32 = syscall.MustLoadDLL("gdi32.dll")
)

var user32CreateWindowExW = user32.MustFindProc("CreateWindowExW")

func createWindowEx(exStyle uint32, className, windowName *uint16, style uint32, x, y, width, height int32, parent _HWND, menu _HMENU, instance _HINSTANCE, param unsafe.Pointer) _HWND {
	_ret, _, _ := user32CreateWindowExW.Call(
		uintptr((exStyle)),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		uintptr((style)),
		uintptr((x)),
		uintptr((y)),
		uintptr((width)),
		uintptr((height)),
		uintptr((parent)),
		uintptr((menu)),
		uintptr((instance)),
		uintptr((param)),
	)
	return _HWND((_ret))
}

var user32GetDC = user32.MustFindProc("GetDC")

func getDC(hWnd _HWND) _HDC {
	_ret, _, _ := user32GetDC.Call(
		uintptr((hWnd)),
	)
	return _HDC((_ret))
}

var user32AdjustWindowRectEx = user32.MustFindProc("AdjustWindowRectEx")

func adjustWindowRectEx(_rect *rect, style uint32, menu bool, exStyle uint32) bool {
	_ret, _, _ := user32AdjustWindowRectEx.Call(
		uintptr(unsafe.Pointer(_rect)),
		uintptr((style)),
		uintptr(boolToUintptr(menu)),
		uintptr((exStyle)),
	)
	return bool(uintptrToBool(_ret))
}

var user32ReleaseDC = user32.MustFindProc("ReleaseDC")

func releaseDC(hwnd _HWND, _HDC _HDC) bool {
	_ret, _, _ := user32ReleaseDC.Call(
		uintptr((hwnd)),
		uintptr((_HDC)),
	)
	return bool(uintptrToBool(_ret))
}

var user32DestroyWindow = user32.MustFindProc("DestroyWindow")

func destroyWindow(hwnd _HWND) bool {
	_ret, _, _ := user32DestroyWindow.Call(
		uintptr((hwnd)),
	)
	return bool(uintptrToBool(_ret))
}

var user32LoadIconW = user32.MustFindProc("LoadIconW")

func loadIcon(hInstance _HINSTANCE, iconName *uint16) _HICON {
	_ret, _, _ := user32LoadIconW.Call(
		uintptr((hInstance)),
		uintptr(unsafe.Pointer(iconName)),
	)
	return _HICON((_ret))
}

var user32LoadCursorW = user32.MustFindProc("LoadCursorW")

func loadCursor(hinstance _HINSTANCE, cursorName *uint16) _HCURSOR {
	_ret, _, _ := user32LoadCursorW.Call(
		uintptr((hinstance)),
		uintptr(unsafe.Pointer(cursorName)),
	)
	return _HCURSOR((_ret))
}

var user32RegisterClassExW = user32.MustFindProc("RegisterClassExW")

func registerClassEx(wcx *_WNDCLASSEX) uint16 {
	_ret, _, _ := user32RegisterClassExW.Call(
		uintptr(unsafe.Pointer(wcx)),
	)
	return uint16((_ret))
}

var user32ShowWindow = user32.MustFindProc("ShowWindow")

func showWindow(hwnd _HWND, nCmdShow int32) bool {
	_ret, _, _ := user32ShowWindow.Call(
		uintptr((hwnd)),
		uintptr((nCmdShow)),
	)
	return bool(uintptrToBool(_ret))
}

var user32SetFocus = user32.MustFindProc("SetFocus")

func setFocus(hwnd _HWND) _HWND {
	_ret, _, _ := user32SetFocus.Call(
		uintptr((hwnd)),
	)
	return _HWND((_ret))
}

var user32DefWindowProcW = user32.MustFindProc("DefWindowProcW")

func defWindowProc(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	_ret, _, _ := user32DefWindowProcW.Call(
		uintptr((hwnd)),
		uintptr((msg)),
		uintptr((wparam)),
		uintptr((lparam)),
	)
	return uintptr((_ret))
}

var user32GetMessageW = user32.MustFindProc("GetMessageW")

func getMessage(lpMsg *_MSG, hwnd uintptr, wMsgFilterMin, wMsgFilterMax uint32) bool {
	_ret, _, _ := user32GetMessageW.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr((hwnd)),
		uintptr((wMsgFilterMin)),
		uintptr((wMsgFilterMax)),
	)
	return bool(uintptrToBool(_ret))
}

var user32PeekMessageW = user32.MustFindProc("PeekMessageW")

func peekMessage(lpMsg *_MSG, hwnd uintptr, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) bool {
	_ret, _, _ := user32PeekMessageW.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr((hwnd)),
		uintptr((wMsgFilterMin)),
		uintptr((wMsgFilterMax)),
		uintptr((wRemoveMsg)),
	)
	return bool(uintptrToBool(_ret))
}

var user32TranslateMessage = user32.MustFindProc("TranslateMessage")

func translateMessage(msg *_MSG) bool {
	_ret, _, _ := user32TranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)),
	)
	return bool(uintptrToBool(_ret))
}

var user32DispatchMessageW = user32.MustFindProc("DispatchMessageW")

func dispatchMessage(msg *_MSG) uintptr {
	_ret, _, _ := user32DispatchMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
	)
	return uintptr((_ret))
}

var user32SetTimer = user32.MustFindProc("SetTimer")

func setTimer(hwnd _HWND, nIDEvent uintptr, uElapse uint32, lpTimerFunc uintptr) uintptr {
	_ret, _, _ := user32SetTimer.Call(
		uintptr((hwnd)),
		uintptr((nIDEvent)),
		uintptr((uElapse)),
		uintptr((lpTimerFunc)),
	)
	return uintptr((_ret))
}

var user32KillTimer = user32.MustFindProc("KillTimer")

func killTimer(hwnd _HWND, uIDEvent uintptr) bool {
	_ret, _, _ := user32KillTimer.Call(
		uintptr((hwnd)),
		uintptr((uIDEvent)),
	)
	return bool(uintptrToBool(_ret))
}

var kernel32GetModuleHandleW = kernel32.MustFindProc("GetModuleHandleW")

func getModuleHandle(moduleName string) _HINSTANCE {
	_ret, _, _ := kernel32GetModuleHandleW.Call(
		uintptr(stringToUintptr(moduleName)),
	)
	return _HINSTANCE((_ret))
}

var kernel32GetLastError = kernel32.MustFindProc("GetLastError")

func getLastError() uint32 {
	_ret, _, _ := kernel32GetLastError.Call()
	return uint32((_ret))
}

var kernel32CreateTimerQueueTimer = kernel32.MustFindProc("CreateTimerQueueTimer")

func createTimerQueueTimer(phNewTimer *uintptr, TimerQueue, Callback, Parameter uintptr, DueTime, Period int32, Flags uint32) bool {
	_ret, _, _ := kernel32CreateTimerQueueTimer.Call(
		uintptr(unsafe.Pointer(phNewTimer)),
		uintptr((TimerQueue)),
		uintptr((Callback)),
		uintptr((Parameter)),
		uintptr((DueTime)),
		uintptr((Period)),
		uintptr((Flags)),
	)
	return bool(uintptrToBool(_ret))
}

var kernel32DeleteTimerQueueTimer = kernel32.MustFindProc("DeleteTimerQueueTimer")

func deleteTimerQueueTimer(TimerQueue, Timer, CompletionEvent uintptr) bool {
	_ret, _, _ := kernel32DeleteTimerQueueTimer.Call(
		uintptr((TimerQueue)),
		uintptr((Timer)),
		uintptr((CompletionEvent)),
	)
	return bool(uintptrToBool(_ret))
}

var gdi32SetPixelFormat = gdi32.MustFindProc("SetPixelFormat")

func setPixelFormat(hdc _HDC, iPixelFormat int32, ppfd *_PIXELFORMATDESCRIPTOR) bool {
	_ret, _, _ := gdi32SetPixelFormat.Call(
		uintptr((hdc)),
		uintptr((iPixelFormat)),
		uintptr(unsafe.Pointer(ppfd)),
	)
	return bool(uintptrToBool(_ret))
}

var gdi32ChoosePixelFormat = gdi32.MustFindProc("ChoosePixelFormat")

func choosePixelFormat(hdc _HDC, ppfd *_PIXELFORMATDESCRIPTOR) int32 {
	_ret, _, _ := gdi32ChoosePixelFormat.Call(
		uintptr((hdc)),
		uintptr(unsafe.Pointer(ppfd)),
	)
	return int32((_ret))
}

var gdi32SwapBuffers = gdi32.MustFindProc("SwapBuffers")

func swapBuffers(hdc _HDC) bool {
	_ret, _, _ := gdi32SwapBuffers.Call(
		uintptr((hdc)),
	)
	return bool(uintptrToBool(_ret))
}
