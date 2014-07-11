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
	disp   *C.Display
	win    C.Window
	Events <-chan Event
}

func (c Context) KillWindow() error {
	return nil
}

func CreateContext(name string, width, height uint32, majorVersion, minorVersion int32) (*Context, error) {
	runtime.LockOSThread()
	C.XInitThreads()
	c := new(Context)
	c.disp = C.XOpenDisplay(nil)

	fmt.Println("disp:", *c.disp)

	if c.disp == nil {
		return nil, fmt.Errorf("Could not open X11 display")
	}

	var dummy C.int
	if 0 == C.glXQueryExtension(c.disp, &dummy, &dummy) {
		return nil, fmt.Errorf("Could not open display")
	}

	root := C.XDefaultRootWindow(c.disp)

	attributes := []C.int{C.GLX_RGBA, C.GLX_DEPTH_SIZE, 24, C.GLX_DOUBLEBUFFER, C.None}

	vi := C.glXChooseVisual(c.disp, C.XDefaultScreen(c.disp), &attributes[0])

	if vi == nil {
		return nil, fmt.Errorf("Could not create X11 visual")
	}

	cx := C.glXCreateContext(c.disp, vi, nil, C.GL_TRUE)
	if cx == nil {
		return nil, fmt.Errorf("Could not create rendering context")
	}
	// x color map
	cmap := C.XCreateColormap(c.disp, root, vi.visual, C.AllocNone)
	swa := C.XSetWindowAttributes{}
	swa.colormap = cmap
	swa.border_pixel = 0
	swa.event_mask = C.ExposureMask | C.KeyPressMask | C.StructureNotifyMask
	c.win = C.XCreateWindow(c.disp, root, 0, 0, C.uint(width), C.uint(height), 0, vi.depth, C.InputOutput, vi.visual, C.CWBorderPixel|C.CWColormap|C.CWEventMask, &swa)

	C.XSetStandardProperties(c.disp, c.win, C.CString("xogl"), C.CString("xogl"), C.None, nil, 0, nil)

	C.glXMakeCurrent(c.disp, C.GLXDrawable(c.win), cx)
	C.XMapWindow(c.disp, c.win)

	evdisp := C.XOpenDisplay(nil)
	if evdisp == nil {
		return nil, fmt.Errorf("Could not open event display")
	}

    events := make(chan Events, 128)
    c.Events = events

	go func() {
        runtime.LockOSThread()
		C.XSelectInput(evdisp, c.win, C.StructureNotifyMask|C.ExposureMask|C.KeyPressMask|C.ButtonPressMask|C.PointerMotionMask)
		xev := C.XEvent{}
		for {
			//C.XLockDisplay(evdisp)
			C.XNextEvent(evdisp, &xev)
			//C.XUnlockDisplay(evdisp)
			//fmt.Println(xev)
		}
	}()

	return c, nil
}

func (c Context) SwapBuffers() {
	C.glXSwapBuffers(c.disp, C.GLXDrawable(c.win))
}
