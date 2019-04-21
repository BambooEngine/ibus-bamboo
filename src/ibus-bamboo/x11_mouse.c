/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#define _GNU_SOURCE
#include <X11/Xlib.h>
#include <stdlib.h>
#include <stdio.h>
#include <pthread.h>
#include <sys/time.h>
#include <unistd.h>
#include <time.h>
#include "_cgo_export.h"
#define CAPTURE_MOUSE_MOVE_DELTA        50

static pthread_t th_mcap;
static pthread_mutex_t mutex_mcap;
static Display* dpy;
static int mcap_running;

static void signalHandler(int signo) {
    mcap_running = signo;
}
/**
 * milliseconds over 1000 will be ignored
 */
static void delay(time_t sec, long msec) {
    struct timespec sleep;

    sleep.tv_sec  = sec;
    sleep.tv_nsec = (msec % 1000) * 1000 * 1000;

    if (nanosleep(&sleep, NULL) == -1) {
        signalHandler(1);
    }
}
/**
 * returns 0 for failure, 1 for success
 */
static int grabPointer(Display *dpy, Window w, unsigned int mask) {
    int rc;

    /* retry until we actually get the pointer (with a suitable delay)
     * or we get an error we can't recover from. */
    while (mcap_running == 1) {
        rc = XGrabPointer(dpy, w, 0, ButtonPressMask | PointerMotionMask, GrabModeAsync, GrabModeAsync, None, None, CurrentTime);

        switch (rc) {
            case GrabSuccess:
                fprintf(stderr, "XGrabPointer: successfully grabbed mouse pointer\n");
                return 1;

            case AlreadyGrabbed:
                fprintf(stderr, "XGrabPointer: already grabbed mouse pointer, retrying with delay\n");
                delay(1, 500);
                break;

            case GrabFrozen:
                fprintf(stderr, "XGrabPointer: grab was frozen, retrying after delay\n");
                delay(1, 500);
                break;

            case GrabNotViewable:
                fprintf(stderr, "XGrabPointer: grab was not viewable, exiting\n");
                return 0;

            case GrabInvalidTime:
                fprintf(stderr, "XGrabPointer: invalid time, exiting\n");
                return 0;

            default:
                fprintf(stderr, "XGrabPointer: could not grab mouse pointer (%d), exiting\n", rc);
                return 0;
        }
    }

    return 0;
}

static void* thread_mouse_capture(void* data)
{
    XEvent event;
    int x_old, y_old, x_root_old, y_root_old, rt;
    unsigned int mask;
    Window w, w_root_return, w_child_return;

    dpy = XOpenDisplay(NULL);
    if (!dpy) {
        return NULL;
    }
    w = XDefaultRootWindow(dpy);

    XQueryPointer(dpy, w, &w_root_return, &w_child_return, &x_root_old, &y_root_old, &x_old, &y_old, &mask);
    while (mcap_running == 1 && grabPointer(dpy, w, mask)) {
        while (mcap_running == 1) {
            if (XPending(dpy) > 0) {
                XPeekEvent(dpy, &event);
                break;
            }
            delay(0, 50);
        }
        XUngrabPointer(dpy, CurrentTime);
        XSync(dpy, 1);
        pthread_mutex_lock(&mutex_mcap); // set mutex to lock status, so this thread will wait until next unlock (by update preedit string)
        if (mcap_running == 0)
            break;

        if (event.type == MotionNotify) // mouse move
        {
            if ((abs(event.xmotion.x_root - x_root_old) >= CAPTURE_MOUSE_MOVE_DELTA) ||
                (abs(event.xmotion.y_root - y_root_old) >= CAPTURE_MOUSE_MOVE_DELTA)) // mouse move at least CAPTURE_MOUSE_MOVE_DELTA
            {
                fprintf(stderr, "MotionNotify: delta_x=%d delta_y=%d\n", abs(event.xmotion.x_root - x_root_old), abs(event.xmotion.y_root - y_root_old));
                mouse_move_handler();

                x_root_old = event.xmotion.x_root;
                y_root_old = event.xmotion.y_root;
            }
            else { // if don't reset -> unlock mutex so mouse continue to be grab
                pthread_mutex_unlock(&mutex_mcap);
            }
        }
        else {
              mouse_move_handler();
        }
    }
    mcap_running = 0;
    XCloseDisplay(dpy);
    return NULL;
}

void mouse_capture_init()
{
    setbuf(stdout, NULL);
    setbuf(stderr, NULL);
    if (mcap_running==1) {
        return;
    }
    XInitThreads();
    mcap_running = 1;
    pthread_mutex_init(&mutex_mcap, NULL);
    pthread_mutex_trylock(&mutex_mcap); // lock mutex after init so mouse capture not start
    pthread_create(&th_mcap, NULL, &thread_mouse_capture, NULL);
    pthread_detach(th_mcap);
}

void mouse_capture_exit()
{
    if (mcap_running==0) {
        return;
    }
    mcap_running = 0;
    pthread_mutex_unlock(&mutex_mcap); // unlock mutex, so thread can exit
}

// every time have preedit text -> unlock mutex -> start capture mouse
void mouse_capture_unlock()
{
    if (mcap_running==0) {
        return;
    }
    // unlock capture thread (start capture)
    pthread_mutex_unlock(&mutex_mcap);
}
