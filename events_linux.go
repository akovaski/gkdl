package gkdl

// #include <X11/X.h>
// #include <X11/Xlib.h>
import "C"
import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

func handleEvents(dpy *C.Display, win C.Window, events chan<- Event) {
	var xev C.XEvent
	var nev C.XEvent // this is used to store peeked events

	for {
		C.XNextEvent(dpy, &xev)
		//fmt.Println(xev)
		etype := binary.LittleEndian.Uint32(xev[0:4])
		switch etype {
		case C.KeyPress:
			var ev KeyDown

			keyv := *(*C.XKeyEvent)(unsafe.Pointer(&xev))
			ev.Scancode = byte(keyv.keycode)
			ev.Keycode = scanCodeToKeyCode(ev.Scancode)

			events <- ev
		case C.KeyRelease:
			var ev KeyUp

			keyv := *(*C.XKeyEvent)(unsafe.Pointer(&xev))

			// Check if this is an autorepeat, ignore it if it is
			if C.XEventsQueued(dpy, C.QueuedAlready) != 0 {
				C.XPeekEvent(dpy, &nev)
				ntype := binary.LittleEndian.Uint32(nev[0:4])
				if ntype == C.KeyPress {
					keyn := *(*C.XKeyEvent)(unsafe.Pointer(&nev))
					if keyn.time == keyv.time && keyn.keycode == keyv.keycode {
						// verified autorepeat key
						C.XNextEvent(dpy, &nev)
						continue
					}
				}
			}

			ev.Scancode = byte(keyv.keycode)
			ev.Keycode = scanCodeToKeyCode(ev.Scancode)

			events <- ev
		case C.ButtonPress:
			mv := *(*C.XButtonEvent)(unsafe.Pointer(&xev))

            switch mv.button {
            case 4,5:
                // this is a scroll event
                var ev MouseWheel
                if mv.button == 4 {
                    ev.WheelDelta = 1
                } else {
                    ev.WheelDelta = -1
                }
                ev.getMouseState(mv.x, mv.y, mv.state)

                events <- ev
            case 1,2,3:
                // regulare button press
                var ev MouseButtonDown
                ev.Button = buttonToButton(mv.button)

                ev.getMouseState(mv.x, mv.y, mv.state)

                events <- ev
            }

        case C.ButtonRelease:
			mv := *(*C.XButtonEvent)(unsafe.Pointer(&xev))

            switch mv.button {
            case 1,2,3:
                var ev MouseButtonUp
                ev.Button = buttonToButton(mv.button)

                ev.getMouseState(mv.x, mv.y, mv.state)

                events <- ev
            }

        case C.MotionNotify:
            mv := *(*C.XMotionEvent)(unsafe.Pointer(&xev))

            var ev MouseMotion
            ev.getMouseState(mv.x, mv.y, mv.state)

            events <- ev

        case C.ClientMessage:
            fmt.Println("Destroy Notify Event")
            events <- Quit{}

        default:
            fmt.Println("Etype", etype)
		}
	}
}

// TODO: properly convert scancode, find how to function the same as windows
func scanCodeToKeyCode(scancode byte) byte {
	return scancode
}

func buttonToButton(button C.uint) MouseButtons {
	var b MouseButtons
	switch button {
	case 1:
		b.L = true
	case 2:
		b.M = true
	case 3:
		b.R = true
	case 8:
		b.X1 = true
    case 9:
		b.X2 = true
	}
	return b
}

func (s *mouseState) getMouseState(x, y C.int, state C.uint) {
	s.X = int(x)
	s.Y = int(y)
	s.State.L = state&(1<<8) != 0
	s.State.M = state&(1<<9) != 0
	s.State.R = state&(1<<10) != 0
	s.State.X1 = state&(1<<15) != 0
	s.State.X2 = state&(1<<16) != 0
}
