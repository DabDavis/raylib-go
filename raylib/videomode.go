package rl

/*
#include "external/glfw/include/GLFW/glfw3.h"

typedef struct {
	int width;
	int height;
	int refreshRate;
} GoVideoMode;

static int goGetVideoModes(GoVideoMode* out, int maxModes) {
	GLFWmonitor* monitor = glfwGetPrimaryMonitor();
	if (!monitor) return 0;
	int count = 0;
	const GLFWvidmode* modes = glfwGetVideoModes(monitor, &count);
	if (!modes) return 0;
	if (count > maxModes) count = maxModes;
	for (int i = 0; i < count; i++) {
		out[i].width = modes[i].width;
		out[i].height = modes[i].height;
		out[i].refreshRate = modes[i].refreshRate;
	}
	return count;
}
*/
import "C"

// VideoMode represents a supported display resolution and refresh rate.
type VideoMode struct {
	Width, Height, RefreshRate int
}

// GetVideoModes returns all display modes supported by the primary monitor.
func GetVideoModes() []VideoMode {
	var buf [256]C.GoVideoMode
	n := C.goGetVideoModes(&buf[0], 256)
	modes := make([]VideoMode, int(n))
	for i := 0; i < int(n); i++ {
		modes[i] = VideoMode{
			Width:       int(buf[i].width),
			Height:      int(buf[i].height),
			RefreshRate: int(buf[i].refreshRate),
		}
	}
	return modes
}
