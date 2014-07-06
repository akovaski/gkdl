package gkdl

type _HWND uintptr
type _HDC uintptr
type _HGLRC uintptr
type _HINSTANCE uintptr
type _HICON uintptr
type _HCURSOR uintptr
type _HBRUSH uintptr
type _HMENU uintptr

type _POINT struct {
	x, y int32
}

type _MSG struct {
	hwnd _HWND
	message uint32
	WParam uintptr
	LParam uintptr
	time uint32
	pt _POINT
}

type rect struct {
	Left, Top, Right, Bottom int32
}

type _WNDCLASSEX struct {
	size uint32
	style uint32
	wndProc uintptr
	clsExtra int32
	wndExtra int32
	hInstance _HINSTANCE
	hIcon _HICON
	hCursor _HCURSOR
	hbrBackground _HBRUSH
	menuName *uint16
	className *uint16
	iconSm _HICON
}

type _PIXELFORMATDESCRIPTOR struct {
	size uint16
	version uint16
	dwFlags uint32
	iPixelType byte
	colorBits byte
	redBits byte
	redShift byte
	greenBits byte
	greenShift byte
	blueBits byte
	blueShift byte
	alphaBits byte
	alphaShift byte
	accumBits byte
	accumRedBits byte
	accumGreenBits byte
	accumBlueBits byte
	accumAlphaBits byte
	depthBits byte
	stencilBits byte
	auxBuffers byte
	iLayerType byte
	reserved byte
	dwLayerMask uint32
	dwVisibleMask uint32
	dwDamageMask uint32
}