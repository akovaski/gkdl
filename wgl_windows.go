package gkdl

import (
	"syscall"
	"unsafe"
	"strings"
	"fmt"
)

var (
	opengl32 = syscall.MustLoadDLL("opengl32.dll")

	wglChoosePixelFormatARB func(hDC _HDC, piAttribIList *int32, pfAttribFList *float32, nMaxFormats uint32, piFormats *int32, nNumFormats *uint32) bool
	wglCreateContextAttribsARB func(hDC _HDC, hsharecontext _HGLRC, attribList *int32) _HGLRC
)

func goString(cp *uint8) (s string) {
	cs := uintptr(unsafe.Pointer(cp))
	if cs != 0 {
		us := make([]byte, 0, 256)
		for {
			u := *(*byte)(unsafe.Pointer(cs))
			if u == 0 {
				return string(us)
			}
			us = append(us, u)
			cs += 1
		}
	}
	return ""
}

var addrWglCreateContext  = opengl32.MustFindProc("wglCreateContext")
func wglCreateContext(hdc _HDC) _HGLRC {
	ret, _, _ := addrWglCreateContext.Call(uintptr(hdc))

	return _HGLRC(ret)
}

var addrWglDeleteContext  = opengl32.MustFindProc("wglDeleteContext")
func wglDeleteContext(hglrc _HGLRC) bool {
	ret, _, _ := addrWglDeleteContext.Call(uintptr(hglrc))

	return ret == 1
}

var addrWglGetCurrentDC   = opengl32.MustFindProc("wglGetCurrentDC")
func wglGetCurrentDC() _HDC {
	ret, _, _ := addrWglGetCurrentDC.Call()

	return _HDC(ret)
}

var wglGetProcAddress = opengl32.MustFindProc("wglGetProcAddress")
func getProcAddress(name string) uintptr {
	ret, _, _ := wglGetProcAddress.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr(name))))

	return ret
}

var addrWglGetExtensionsStringARB uintptr
func wglGetExtensionsStringARB (hDC _HDC) string {
	if addrWglGetExtensionsStringARB == 0 {
		addrWglGetExtensionsStringARB = getProcAddress("wglGetExtensionsStringARB")
	}
	ret, _, _ := syscall.Syscall(addrWglGetExtensionsStringARB, 1, uintptr(hDC), 0, 0)

	return goString((*uint8)(unsafe.Pointer(ret)))
}

var addrWglMakeCurrent = opengl32.MustFindProc("wglMakeCurrent")
func wglMakeCurrent(hDC _HDC, hRC _HGLRC) bool {
	ret, _, _ := addrWglMakeCurrent.Call(uintptr(hDC), uintptr(hRC))
	return ret != 0
}

var wgl_extensions map[string]bool
func supportWglExtension(ext string) bool {
	if wgl_extensions == nil {
		wgl_extensions = make(map[string]bool)
		exts := strings.Split(wglGetExtensionsStringARB(wglGetCurrentDC()), " ")
		for _, e := range(exts) {
			wgl_extensions[e] = true
		}
	}
	return wgl_extensions[ext]
}

//var wglChoosePixelFormatARB func(hDC HDC, piAttribIList *int32, pfAttribFList *float32, nMaxFormats uint32, piFormats *int32, nNumFormats *uint32) bool

func wglInit() error {
	exts := []string{"WGL_ARB_pixel_format", "WGL_ARB_create_context"}
	for _, e := range exts {
		if supportWglExtension(e) == false {
			return fmt.Errorf("Could not load needed Wgl extension: %s", e)
		}
	}

	{// WGL_ARB_pixel_format
	addr := getProcAddress("wglChoosePixelFormatARB")
	wglChoosePixelFormatARB = func(hDC _HDC, piAttribIList *int32, pfAttribFList *float32, nMaxFormats uint32, piFormats *int32, nNumFormats *uint32) bool {
		ret, _, _ := syscall.Syscall6(addr, 6, uintptr(hDC), uintptr(unsafe.Pointer(piAttribIList)), uintptr(unsafe.Pointer(pfAttribFList)), uintptr(nMaxFormats), uintptr(unsafe.Pointer(piFormats)), uintptr(unsafe.Pointer(nNumFormats)))

		return ret != 0
	}
	}

	{// WGL_ARB_create_context
	addr := getProcAddress("wglCreateContextAttribsARB")
	wglCreateContextAttribsARB = func(hDC _HDC, hsharecontext _HGLRC, attribList *int32) _HGLRC {
		ret, _, _ := syscall.Syscall(addr, 3, uintptr(hDC), uintptr(hsharecontext), uintptr(unsafe.Pointer(attribList)))

		return _HGLRC(ret)
	}
	}

	return nil
}