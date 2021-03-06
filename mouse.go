package main

/*
#cgo LDFLAGS: -lX11
#include <stdio.h>
#include <X11/Xlib.h>

extern void onMouseEvent(int, int, int);

Window root_window;
unsigned int mask;
Display *display;
Bool running;
XEvent evt;
static void init() {
	display = XOpenDisplay(NULL);
}

static void startMouseEventListener(){

    if(display == NULL){
		return;
	}

    XAllowEvents(display, AsyncBoth, CurrentTime);
	XGrabPointer(display, root_window, 0, PointerMotionMask | ButtonReleaseMask | ButtonPressMask, GrabModeAsync, GrabModeAsync, None, None, CurrentTime);
	running = True;
    while(running) {
		XNextEvent(display, &evt);
		if (evt.type == ButtonPress) {
			onMouseEvent(evt.xbutton.button, 1, 0);
			continue;
		}
		if (evt.type == ButtonRelease) {
			onMouseEvent(evt.xbutton.button, 0, 0);
			continue;
		}
		if (evt.type == MotionNotify) {
			onMouseEvent(-1, evt.xmotion.x_root, evt.xmotion.y_root);
			continue;
		}
	}
	XUngrabPointer(display, CurrentTime);
	XFlush(display);
}

static void releaseMouse() {
	running = False;
}

static int* getMouse () {
	if (display == NULL){
		printf("You need to run init first!");
		static int r[2];
		return r;
	}
	int root_x, root_y;
	XQueryPointer(display, DefaultRootWindow(display), &root_window, &root_window, &root_x, &root_y, &root_x, &root_y, &mask);
	static int  r[2];
	r[0] = root_x;
	r[1] = root_y;
	return r;
}

static void setMousePos(int x, int y) {
	if (display == NULL){
		printf("You need to run init first!");
		return;
	}
	root_window = XRootWindow(display, 0);
	XSelectInput(display, root_window, KeyReleaseMask);
	XWarpPointer(display, None, root_window, 0, 0, 0, 0, x, y);
	XFlush(display);
}
*/
import "C"
import (
	"unsafe"
)

var callback func(int, int, int)

//export onMouseEvent
func onMouseEvent(a, b, c C.int) {
	button := (int)(a)
	state := (int)(b)
	callback(button, state, (int)(c))
}

func mouseInit() {
	C.init()
}

func startMouseListener(call func(int, int, int)) {
	callback = call
	go C.startMouseEventListener()
}

func releaseMouse() {
	C.releaseMouse()
}

func setMousePos(x, y int) {
	C.setMousePos((C.int)(x), (C.int)(y))
}

func getMousePos() (int, int) {
	var mousePos *C.int = C.getMouse()
	length := 2
	slice := (*[1 << 28]C.int)(unsafe.Pointer(mousePos))[:length:length]
	return (int)(slice[0]), (int)(slice[1])
}
