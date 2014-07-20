package gkdl

// #cgo LDFLAGS: -lX11 -lGL
// #include <X11/X.h>
// #include <X11/Xlib.h>
// #include <GL/gl.h>
// #include <GL/glx.h>
import "C"
import (
	"fmt"
	"runtime"
	//"unsafe"
)

type Context struct {
	dpy   *C.Display // Output Display
    edpy   *C.Display // Event Display
    glc   C.GLXContext
	win    C.Window
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

	c.glc = C.glXCreateContext(c.dpy, vi, nil, C.GL_TRUE)
	if c.glc == nil {
		return nil, fmt.Errorf("Could not create rendering context")
	}
	// x color map
	cmap := C.XCreateColormap(c.dpy, root, vi.visual, C.AllocNone)
	swa := C.XSetWindowAttributes{}
	swa.colormap = cmap
	swa.border_pixel = 0
	c.win = C.XCreateWindow(c.dpy, root, 0, 0, C.uint(width), C.uint(height), 0, vi.depth, C.InputOutput, vi.visual, C.CWBorderPixel|C.CWColormap, &swa)

	C.XSetStandardProperties(c.dpy, c.win, C.CString("xogl"), C.CString("xogl"), C.None, nil, 0, nil)

	C.glXMakeCurrent(c.dpy, C.GLXDrawable(c.win), c.glc)
	C.XMapWindow(c.dpy, c.win)

	events := make(chan Event, 128)
	c.Events = events
    //wmDelete := C.XInternAtom(c.dpy, C.CString("WM_DELETE_WINDOW"), 0)
    //C.XSetWMProtocols(c.dpy, c.win, &wmDelete, 1)

	go func() {
		runtime.LockOSThread()
        wmDelete2 := C.XInternAtom(c.edpy, C.CString("WM_DELETE_WINDOW"), 1)
        C.XSetWMProtocols(c.edpy, c.win, &wmDelete2, 1)
		//C.XSelectInput(evdisp, c.win, C.StructureNotifyMask|C.ExposureMask|C.KeyPressMask|C.KeyReleaseMask|C.ButtonPressMask|C.PointerMotionMask)
		C.XSelectInput(c.edpy, c.win,
            C.KeyPressMask|C.KeyReleaseMask| // Key Events
            C.ButtonPressMask|C.ButtonReleaseMask| // Mouse Button Events
            C.PointerMotionMask|
            C.StructureNotifyMask|C.ExposureMask, // Mouse Motion Events
        )
		handleEvents(c.edpy, c.win, events)
	}()

	return c, nil
}

func (c *Context) NextEvent() C.XEvent {
    var xev C.XEvent
    C.XNextEvent(c.dpy, &xev)
    return xev
}

func (c Context) SwapBuffers() {
	C.glXSwapBuffers(c.dpy, C.GLXDrawable(c.win))
}
