#include <X11/Xlib.h>
#include <stdlib.h>
#include <pthread.h>
#include "_cgo_export.h"
#define CAPTURE_MOUSE_MOVE_DELTA        50

static pthread_t th_mcap;
static pthread_mutex_t mutex_mcap;
static Display* dpy;
static int mcap_running;

void* thread_mouse_capture(void* data)
{
    XEvent event;
    int x_old, y_old, x_root_old, y_root_old, rt;
    int mask;
    Window w, w_root_return, w_child_return;

    dpy = XOpenDisplay(NULL);
    w = XDefaultRootWindow(dpy);

    XQueryPointer(dpy, w, &w_root_return, &w_child_return, &x_root_old, &y_root_old, &x_old, &y_old, &mask);
    while (mcap_running == 1)
    {
        pthread_mutex_lock(&mutex_mcap);
        if (mcap_running == 0)
            return NULL;
        rt = XGrabPointer(dpy, w, 0, ButtonPressMask | PointerMotionMask, GrabModeAsync, GrabModeAsync, None, None, CurrentTime);
        pthread_mutex_trylock(&mutex_mcap); // set mutex to lock status, so this thread will wait until next unlock (by update preedit string)
        if (rt != 0)
            continue;
        XPeekEvent(dpy, &event);
        XUngrabPointer(dpy, CurrentTime);
        XSync(dpy, 1);

        if (event.type == MotionNotify) // mouse move
        {
            if ((abs(event.xmotion.x_root - x_root_old) >= CAPTURE_MOUSE_MOVE_DELTA) ||
                (abs(event.xmotion.y_root - y_root_old) >= CAPTURE_MOUSE_MOVE_DELTA)) // mouse move at least CAPTURE_MOUSE_MOVE_DELTA
            {
        		mouse_move_handler();

                x_root_old = event.xmotion.x_root;
                y_root_old = event.xmotion.y_root;
            }
            else // if don't reset -> unlock mutex so mouse continue to be grab
                pthread_mutex_unlock(&mutex_mcap);
        }
        else
    		  mouse_move_handler();
    }

    XCloseDisplay(dpy);

    return NULL;
}

void mouse_capture_init()
{
    mcap_running = 1;
    pthread_mutex_init(&mutex_mcap, NULL);
    pthread_mutex_trylock(&mutex_mcap); // lock mutex after init so mouse capture not start
    pthread_create(&th_mcap, NULL, &thread_mouse_capture, NULL);
    pthread_detach(th_mcap);
}

void mouse_capture_exit()
{
    mcap_running = 0;
    pthread_mutex_unlock(&mutex_mcap); // unlock mutex, so thread can exit
}

// every time have preedit text -> unlock mutex -> start capture mouse
void mouse_capture_unlock()
{
    // unlock capture thread (start capture)
    pthread_mutex_unlock(&mutex_mcap);
}