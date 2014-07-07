package gkdl

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync/atomic"
	"syscall"
	"unsafe"
)

type Context struct {
	hwnd      _HWND
	hDC       _HDC
	hRC       _HGLRC
	hInstance _HINSTANCE
	events    chan Event
	Events    <-chan Event
}

func (c Context) KillWindow() error {
	//kill hRC rendering context
	if c.hRC != 0 {
		if wglMakeCurrent(0, 0) == false { // attempt to release DC and RC Context
			return errors.New("Release of DC and RC Failed")
		}

		if wglDeleteContext(c.hRC) == false {
			return errors.New("Release Rendering Contex Failed")
		}
		c.hRC = 0
	}

	//kill hDC Device Context
	if c.hDC != 0 && releaseDC(c.hwnd, c.hDC) == false {
		return errors.New("Release Device Context Failed")
		c.hDC = 0
	}

	//destroy hwnd window handler
	if c.hwnd != 0 && destroyWindow(c.hwnd) == false {
		return errors.New("Could not release hwnd")
		c.hwnd = 0
	}

	close(c.events)

	return nil
}

/* TODO:
- Customizable Pixel attributes and context attributes
*/
var classCounter uint64 = 0

func CreateContext(name string, width, height int32, majorVersion, minorVersion int32) (*Context, error) {
	runtime.LockOSThread()

	c := new(Context)
	errc := make(chan error)

	go func() {
		runtime.LockOSThread()
		c.events = make(chan Event, 128)
		c.Events = c.events
		c.hInstance = getModuleHandle("")
		log.Println("HINSTANCE: ", int64(c.hInstance))

		var className string = "OpenGL_" + strconv.FormatUint(atomic.AddUint64(&classCounter, 1), 36)
		wcname := syscall.StringToUTF16Ptr(className)
		wc := new(_WNDCLASSEX)
		wc.size = uint32(unsafe.Sizeof(*wc))
		wc.style = _CS_OWNDC
		wc.wndProc = syscall.NewCallback(CreateWndProc(c.events))
		wc.hInstance = c.hInstance
		wc.hIcon = loadIcon(0, (*uint16)(unsafe.Pointer(uintptr(_IDI_APPLICATION))))
		wc.hCursor = loadCursor(0, (*uint16)(unsafe.Pointer(uintptr(_IDC_ARROW))))
		wc.hbrBackground = _COLOR_BTNFACE + 1
		wc.menuName = nil
		wc.className = wcname
		wc.iconSm = loadIcon(0, (*uint16)(unsafe.Pointer(uintptr(_IDI_APPLICATION))))

		// Register Class
		if registerClassEx(wc) == 0 {
			errc <- errors.New("Unable to RegisterClassEx")
			return
		}
		println("registered")

		// Create Window Handle
		c.hwnd = createWindow(name, className, width, height, c.hInstance)
		if c.hwnd == 0 {
			errc <- errors.New("Unable to create window")
			return
		}
		println("window created", c.hwnd)

		// Create Device Context
		c.hDC = getDC(c.hwnd)
		log.Println("HDC:", c.hDC)
		if c.hDC == 0 {
			c.KillWindow()
			errc <- errors.New("Can't create a GL Device Context.")
			return
		}

		// Set Pixel Format
		if err := _setPixelFormat(c, 24); err != nil {
			errc <- err
			return
		}

		// Create OpenGL context
		c.hRC = wglCreateContext(c.hDC)
		if c.hRC == 0 {
			c.KillWindow()
			errc <- errors.New("Can't create a GL Rendering Context.")
			return
		}

		// Make OpenGL context current
		if wglMakeCurrent(c.hDC, c.hRC) == false {
			c.KillWindow()
			errc <- errors.New("Can't activate the GL Rendering Context.")
			return
		}

		// load needed wgl functions
		if err := wglInit(); err != nil {
			c.KillWindow()
			errc <- err
			return
		}

		// release the current context from this tread
		if wglMakeCurrent(0, 0) == false {
			c.KillWindow()
			errc <- errors.New("Can't release HDC and HRC.")
			return
		}

		// call wglCreateContextAttribsARB to re-create context if possible
		if supportWglExtension("WGL_ARB_pixel_format") {

			// kill window and contexts, then remake with extension
			c.KillWindow()

			c.hwnd = createWindow(name, className, width, height, c.hInstance)
			if c.hwnd == 0 {
				errc <- errors.New("Unable to create window")
				return
			}

			// Create Device Context
			c.hDC = getDC(c.hwnd)
			if c.hDC == 0 {
				c.KillWindow()
				errc <- errors.New("Can't create a GL Device Context.")
				return
			}

			// Set Pixel Format
			pixList := [...]int32{
				0x2010, 1, //WGL_SUPPORT_OPENGL_ARB, GL_TRUE,
				0x2001, 1, //WGL_DRAW_TO_WINDOW_ARB, GL_TRUE,
				0x2003, 0x2027, //WGL_ACCELERATION_ARB, WGL_FULL_ACCELERATION_ARB
				0x2014, 24, //WGL_COLOR_BITS_ARB, 32,
				0x2022, 24, //WGL_DEPTH_BITS_ARB, 24,
				0x2011, 1, //WGL_DOUBLE_BUFFER_ARB, GL_TRUE,
				0x2007, 0x2028, //WGL_SWAP_METHOD_ARB, WGL_SWAP_EXCHANGE_ARB
				0x2013, 0x202B, //WGL_PIXEL_TYPE_ARB, WGL_TYPE_RGBA_ARB,
				0x2023, 8, //WGL_STENCIL_BITS_ARB, 8,
				0, //End
			}

			var pixelFormat int32
			var numFormats uint32

			ok := wglChoosePixelFormatARB(c.hDC, &pixList[0], nil, 1, &pixelFormat, &numFormats)
			if ok == false {
				c.KillWindow()
				errc <- errors.New("Can't find a suitable ARB PixelFormat.")
				return
			}
			var pfd _PIXELFORMATDESCRIPTOR
			if setPixelFormat(c.hDC, pixelFormat, &pfd) == false {
				c.KillWindow()
				errc <- errors.New("Can't set the ARB PixelFormat.")
				return
			}

			attribList := [...]int32{
				0x2091, majorVersion, // WGL_CONTEXT_MAJOR_VERSION_ARB, majorVersion
				0x2092, minorVersion, // WGL_CONTEXT_MINOR_VERSION_ARB, minorVersion
				0x2094, 0, //WGL_CONTEXT_FLAGS_ARB, 0,
				0,
			}

			c.hRC = wglCreateContextAttribsARB(c.hDC, 0, &attribList[0])
			if c.hRC == 0 {
				c.KillWindow()
				errc <- errors.New("Can't create a GL ARB Rendering Context.")
				return
			}
		}

		// Some Unknown Error
		if e := getLastError(); e != 0 {
			c.KillWindow()
			errc <- fmt.Errorf("Unable to create valid OpenGL context, error: %x\n", e)
			return
		}

		showWindow(c.hwnd, _SW_SHOW)
		setFocus(c.hwnd)

		close(errc)

		//setTimer(0, 0, 1000, 0) // create a time so getMessage can't wait indefinitely
		var msg _MSG
		for {
			getMessage(&msg, 0, 0, 0)
			translateMessage(&msg)
			dispatchMessage(&msg)
		}
	}()

	if err := <-errc; err != nil {
		return nil, err
	}

	// Make OpenGL context current
	if wglMakeCurrent(c.hDC, c.hRC) == false {
		c.KillWindow()
		return nil, errors.New("Can't make current HDC and HRC")
	}

	return c, nil
}

func createWindow(name, className string, width, height int32, hInstance _HINSTANCE) _HWND {
	winRect := rect{Left: 0, Top: 0, Right: 0 + width, Bottom: 0 + height}

	var dwExStyle uint32 = _WS_EX_APPWINDOW | _WS_EX_WINDOWEDGE
	var dwStyle uint32 = _WS_OVERLAPPEDWINDOW

	// Fix window dimensions for borders
	adjustWindowRectEx(&winRect, dwStyle, false, dwExStyle)

	hwnd := createWindowEx(dwExStyle,
		syscall.StringToUTF16Ptr(className),
		syscall.StringToUTF16Ptr(name),
		_WS_CLIPSIBLINGS|_WS_CLIPCHILDREN|dwStyle,
		0, 0,
		winRect.Right-winRect.Left,
		winRect.Bottom-winRect.Top,
		0, 0,
		hInstance,
		unsafe.Pointer(nil))

	return hwnd
}

func _setPixelFormat(c *Context, bits byte) error {
	var pfd _PIXELFORMATDESCRIPTOR
	pfd = _PIXELFORMATDESCRIPTOR{
		uint16(unsafe.Sizeof(pfd)),
		1,
		_PFD_DRAW_TO_WINDOW | _PFD_SUPPORT_OPENGL | _PFD_DOUBLEBUFFER,
		_PFD_TYPE_RGBA,
		bits,
		0, 0, 0, 0, 0, 0,
		0,
		0,
		0,
		0, 0, 0, 0,
		16,
		0,
		0,
		_PFD_MAIN_PLANE,
		0,
		0, 0, 0}

	// Set Pixel Format
	PixelFormat := choosePixelFormat(c.hDC, &pfd)
	if PixelFormat == 0 {
		c.KillWindow()
		return errors.New("Can't find a suitable PixelFormat.")
	}
	if setPixelFormat(c.hDC, PixelFormat, &pfd) == false {
		c.KillWindow()
		return errors.New("Can't set the PixelFormat.")
	}
	return nil
}

func (c Context) SwapBuffers() {
	swapBuffers(c.hDC)
}
