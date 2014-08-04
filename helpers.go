package gkdl

// #include <GL/gl.h>
import "C"
import (
    "strings"
    "unsafe"
)

var gl_extensions map[string]bool

func SupportExtension(ext string) bool {
	/*if strings.HasPrefix(ext, "GL_VERSION_") {
		var major, minor int
		fmt.Sscanf(ext, "GL_VERSION_%d_%d", &major, &minor)
		return SupportVersion(major, minor)
	}*/
	if gl_extensions == nil { // load up a map of available extensions
		gl_extensions = make(map[string]bool)

        // pre 3.0 way of checking extensions
        exts := C.GoString((*C.char)(unsafe.Pointer(C.glGetString(C.GL_EXTENSIONS))))
        split := strings.Split(exts, " ")
        for _, s := range split {
            gl_extensions[s] = true
        }
	}
	return gl_extensions[ext]
}
