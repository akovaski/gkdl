package gkdl

import (
	"log"
)

const (
	lButtonMask  = 0x1
	rButtonMask  = 0x2
	mButtonMask  = 0x10
	x1ButtonMask = 0x20
	x2ButtonMask = 0x40
)

/* TODO:
Deal with alt, ctrl, ... keys ghosting
*/
func CreateWndProc(events chan<- Event) func(hwnd _HWND, msg uint32, wparam, lparam uintptr) uintptr {
	return func(hwnd _HWND, msg uint32, wparam, lparam uintptr) uintptr {
		switch msg {
		case _WM_SYSKEYDOWN, _WM_KEYDOWN:
			if lparam&(1<<30) != 0 {
				// key was held down, not pressed
				return 0
			}

			var ev KeyDown
			ev.Scancode = byte((lparam >> 16) & 0xff)
			ev.Keycode = byte(wparam)
			events <- ev
			return 0

		case _WM_KEYUP, _WM_SYSKEYUP:
			var ev KeyUp
			ev.Scancode = byte((lparam >> 16) & 0xff)
			ev.Keycode = byte(wparam)
			events <- ev
			return 0
		case _WM_MOUSEMOVE:
			var ev MouseMotion
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 0

		case _WM_LBUTTONUP:
			var ev MouseButtonUp
			ev.Button.L = true
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 0

		case _WM_RBUTTONUP:
			var ev MouseButtonUp
			ev.Button.R = true
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 0

		case _WM_MBUTTONUP:
			var ev MouseButtonUp
			ev.Button.L = true
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 0

		case _WM_XBUTTONUP:
			var ev MouseButtonUp

			if wparam&(1<<16) != 0 {
				ev.Button.X1 = true
			} else {
				ev.Button.X2 = true
			}
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 1

		case _WM_LBUTTONDOWN:
			var ev MouseButtonDown
			ev.Button.L = true
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 0

		case _WM_RBUTTONDOWN:
			var ev MouseButtonDown
			ev.Button.R = true
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 0

		case _WM_MBUTTONDOWN:
			var ev MouseButtonDown
			ev.Button.M = true
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 0

		case _WM_XBUTTONDOWN:
			var ev MouseButtonDown

			if wparam&(1<<16) != 0 {
				ev.Button.X1 = true
			} else {
				ev.Button.X2 = true
			}
			ev.getMouseState(wparam, lparam)
			events <- ev
			return 1

		case _WM_MOUSEWHEEL:
			var ev MouseWheel
			ev.getMouseState(wparam, lparam)
            if (wparam >> 16) > 0 {
                ev.WheelDelta = 1
            } else {
                ev.WheelDelta = -1
            }
			events <- ev
			return 0

		case _WM_CLOSE:
			var ev Quit
			events <- ev
			return 0

		case _WM_SYSCOMMAND:
			switch wparam {
			case 0xf140, 0xf170: // w32.SC_SCREENSAVE, w32.SC_MONITORPOWER:
				return 0
			}
			/*
				case WM_SHOWWINDOW:
				case WM_ACTIVATE:

				case WM_MOUSELEAVE:
				case WM_UNICHAR:
				case WM_CHAR:
				case WM_INPUTLANGCHANGE:
				case WM_GETMINMAXINFO:
				case WM_WINDOWPOSCHANGED:
				case WM_SETCURSOR:
				case WM_PAINT:
				case WM_ERASEBKGND:
				case w32.WM_INPUT:
				case WM_TOUCH:
				case WM_DROPFILES:
				case WM_INPUT: // Mouse relative, ...
			*/
		}

		return defWindowProc(uintptr(hwnd), msg, wparam, lparam)
	}
}

const (
	lButtonMask  = 0x1
	rButtonMask  = 0x2
	mButtonMask  = 0x10
	x1ButtonMask = 0x20
	x2ButtonMask = 0x40
)

func (s *mouseState) getMouseState(wparam, lparam uintptr) {
	s.X = int(int16(lparam & 0x0000ffff))
	s.Y = int(int16(lparam >> 16))
	s.State.L = wparam&lButtonMask != 0
	s.State.R = wparam&rButtonMask != 0
	s.State.M = wparam&mButtonMask != 0
	s.State.X1 = wparam&x1ButtonMask != 0
	s.State.X2 = wparam&x2ButtonMask != 0
}
