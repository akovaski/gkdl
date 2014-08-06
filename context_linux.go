package gkdl

// #cgo LDFLAGS: -lX11 -lGL
// #include <X11/X.h>
// #include <X11/Xlib.h>
// #include <GL/gl.h>
// #include <GL/glx.h>
// GLXContext CreateContextAttribsARB(
//              Display *dpy, GLXFBConfig config,
//              GLXContext share_context, Bool direct,
//              const int *attrib_list) {
//      PFNGLXCREATECONTEXTATTRIBSARBPROC glXCreateContextAttribsARB =
//      (PFNGLXCREATECONTEXTATTRIBSARBPROC)glXGetProcAddress((GLubyte*)"glXCreateContextAttribsARB");
//      return glXCreateContextAttribsARB(dpy, config, share_context, direct, attrib_list);
// }
import "C"
import (
	"encoding/binary"
	"fmt"
	"runtime"
	"unsafe"
)

type Context struct {
	dpy    *C.Display // Output Display
	edpy   *C.Display // Event Display
	glc    C.GLXContext
	win    C.Window
	events chan Event
	Events <-chan Event
}

func (c Context) KillWindow() error {
	C.glXMakeCurrent(c.dpy, C.None, nil)
	C.glXDestroyContext(c.dpy, c.glc)

	C.XDestroyWindow(c.dpy, c.win)

	fmt.Println("closing display")
	C.XCloseDisplay(c.dpy)
	return nil
}

func CreateContext(name string, width, height uint32, majorVersion, minorVersion int32) (*Context, error) {
	runtime.LockOSThread()
	C.XInitThreads()
	c := new(Context)
	c.dpy = C.XOpenDisplay(nil)

	if c.dpy == nil {
		return nil, fmt.Errorf("Could not open X11 display")
	}

	c.edpy = C.XOpenDisplay(nil)
	if c.edpy == nil {
		return nil, fmt.Errorf("Could not open event display")
	}

	var dummy C.int
	if 0 == C.glXQueryExtension(c.dpy, &dummy, &dummy) {
		return nil, fmt.Errorf("Could not open display")
	}

	root := C.XDefaultRootWindow(c.dpy)

	attributes := []C.int{C.GLX_RGBA, C.GLX_DEPTH_SIZE, 24, C.GLX_DOUBLEBUFFER, C.None}

	vi := C.glXChooseVisual(c.dpy, C.XDefaultScreen(c.dpy), &attributes[0])

	if vi == nil {
		return nil, fmt.Errorf("Could not create X11 visual")
	}

	// x color map
	cmap := C.XCreateColormap(c.dpy, root, vi.visual, C.AllocNone)
	swa := C.XSetWindowAttributes{}
	swa.colormap = cmap
	swa.border_pixel = 0
	c.win = C.XCreateWindow(c.dpy, root, 0, 0, C.uint(width), C.uint(height), 0, vi.depth, C.InputOutput, vi.visual, C.CWBorderPixel|C.CWColormap, &swa)

	C.XSetStandardProperties(c.dpy, c.win, C.CString(name), C.CString(name), C.None, nil, 0, nil)

	// Trying to get a version greater than 3.0 is untested, but might work
	if majorVersion > 3 || (majorVersion == 3 && minorVersion > 0) {
		if SupportExtension("GLX_ARB_create_context") == false {
			return nil, fmt.Errorf("Could not create gl context of specified version")
		}

		gl3attr := [...]C.int{C.GLX_CONTEXT_MAJOR_VERSION_ARB, C.int(majorVersion),
			C.GLX_CONTEXT_MINOR_VERSION_ARB, C.int(minorVersion),
			C.None}
		var elemc C.int
		fbcfg := C.glXChooseFBConfig(c.dpy, vi.screen, nil, &elemc)
		c.glc = C.CreateContextAttribsARB(c.dpy, *fbcfg, nil, 1, &gl3attr[0])
	} else {
		c.glc = C.glXCreateContext(c.dpy, vi, nil, C.GL_TRUE)
		if c.glc == nil {
			return nil, fmt.Errorf("Could not create rendering context")
		}
	}

	C.glXMakeCurrent(c.dpy, C.GLXDrawable(c.win), c.glc)

	major, minor := GetVersion()
	fmt.Println(major, minor)
	if major < int(majorVersion) || (major == int(majorVersion) && minor < int(minorVersion)) {
		return nil, fmt.Errorf("Could not create an appropriate OpenGL context")
	}

	C.XMapWindow(c.dpy, c.win)

	c.events = make(chan Event, 128)
	c.Events = c.events

	wmDelete := C.XInternAtom(c.dpy, C.CString("WM_DELETE_WINDOW"), 0)
	C.XSetWMProtocols(c.dpy, c.win, &wmDelete, 1)

	go func() {
		runtime.LockOSThread()
		//C.StructureNotifyMask|C.ExposureMask
		C.XSelectInput(c.edpy, c.win,
			C.KeyPressMask|C.KeyReleaseMask| // Key Events
				C.ButtonPressMask|C.ButtonReleaseMask| // Mouse Button Events
				C.PointerMotionMask, // Mouse Motion Events
		)
		handleEvents(c.edpy, c.win, c.events)
	}()

	return c, nil
}

func (c Context) SwapBuffers() {
	C.glXSwapBuffers(c.dpy, C.GLXDrawable(c.win))

	// retrieve any ClientMessage events (because they aren't caught by the event go-routine
	var xev C.XEvent
	for C.XEventsQueued(c.dpy, C.QueuedAlready) != 0 {
		C.XNextEvent(c.dpy, &xev)
		etype := binary.LittleEndian.Uint32(xev[0:4])
		switch etype {
		case C.ClientMessage:
			mess := *(*C.XClientMessageEvent)(unsafe.Pointer(&xev))
			switch *(*C.long)(unsafe.Pointer(&mess.data)) {
			case C.long(C.XInternAtom(c.dpy, C.CString("WM_DELETE_WINDOW"), 1)):
				fmt.Println("WM_DELETE_WINDOW message")
				go func() { c.events <- Quit{} }()
			default:
				fmt.Printf("Unhandled Client Message %+v\n", mess)
			}
		default:
			fmt.Printf("Non-Client Message event appeared where it shouldn't %+v\n", xev)
		}
	}
}
