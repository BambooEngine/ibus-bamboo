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

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <pthread.h>
#include <X11/Xlib.h>
#include <string.h> // strlen

#define MAX_TEXT_LEN 100
static pthread_t th_input_watch;
#define MaxPropertyLen 128
#define MaxWmClassesLen 5
static char * WM_CLASS = "WM_CLASS";
static char * WM_NAME = "WM_NAME";
static char * text = NULL;

static int ignore_x_error(Display *display, XErrorEvent *error) {
    return 0;
}

void setXIgnoreErrorHandler() {
    XSetErrorHandler(ignore_x_error);
}

char* uchar2char(unsigned char* uc, unsigned long len) {
    for (int i=0; i<len; i++) {
        if (uc[i] == 0 && i+1 < len) {
            uc[i] = ':';
        }
    }
    return (char*)uc;
 }

char * x11GetStringProperty(Display *display, Window window, char * propName) {
    Atom actualType, filterAtom;
    int status, actualFormat = 0;
    unsigned long len, bytesAfter;
    unsigned char * uc = NULL;

    if (display == NULL) {
        return NULL;
    }
    filterAtom = XInternAtom(display, propName, True);
    status = XGetWindowProperty(display, window, filterAtom, 0, MaxPropertyLen, False, AnyPropertyType,
        &actualType, &actualFormat, &len, &bytesAfter, &uc);
    if (status == Success && uc != NULL) {
        return uchar2char(uc, len);
    }
    if (uc != NULL) {
        XFree(uc);
    }
    return NULL;
}

char * x11GetFocusWindowClassByProp(Display *display, char * propName) {
    Window w;
    int revertTo;
    if (display == NULL) {
        return NULL;
    }
    XGetInputFocus(display, &w, &revertTo);
    for (int i=0; i<MaxWmClassesLen; i++) {
        char * strClass = x11GetStringProperty(display, w, propName);
        if (strClass != NULL && strstr(strClass, "FocusProxy") == NULL) {
            return strClass;
        }
        if (strClass != NULL) {
            XFree(strClass);
        }
        Window * childrenWindows = NULL;
        Window parentWindow, rootWindow;
        unsigned int nChild = 0;
        XQueryTree(display, w, &rootWindow, &parentWindow, &childrenWindows, &nChild);
        if (childrenWindows != NULL) {
            XFree(childrenWindows);
        }
        if (rootWindow == parentWindow) {
            break;
        }
        w = parentWindow;
    }
    return NULL;
}

char * x11GetFocusWindowClassByDpy(Display *display) {
    char * strClass = x11GetFocusWindowClassByProp(display, WM_CLASS);
    if (strClass == NULL) {
        strClass = x11GetFocusWindowClassByProp(display, "_GTK_APPLICATION_ID");
    }
    return strClass;
}

char * x11GetFocusWindowClass() {
    Display * dpy;
    dpy = XOpenDisplay(NULL);
    if (dpy == NULL) {
        // No X display (e.g. pure Wayland session): no class to report.
        return NULL;
    }
    char * wm = x11GetFocusWindowClassByDpy(dpy);
    XCloseDisplay(dpy);
    return wm;
}

static int input_watching = 0;
static int th_count = 0;
static void* thread_input_watching(void* data)
{
    XEvent event;
    int x_old, y_old, x_root_old, y_root_old, rt;
    unsigned int mask;
    Window w, w_root_return, w_child_return;
    Display * dpy;

    dpy = XOpenDisplay(NULL);
    setXIgnoreErrorHandler();
    if (!dpy) {
        return NULL;
    }
    int revertTo;
    XGetInputFocus(dpy, &w, &revertTo);
    XSelectInput(dpy, w, FocusChangeMask);
    char * name = NULL;
    text = (char*)calloc(MAX_TEXT_LEN, sizeof(char));
    char * cl = x11GetFocusWindowClassByDpy(dpy);
    if (cl != NULL) {
      strncpy(text, cl, MAX_TEXT_LEN - 1);
      XFree(cl);
    }
    while (input_watching == 1) {
        XNextEvent(dpy, &event);
        /* text = (char*)calloc(MAX_TEXT_LEN, sizeof(char)); */
        memset(text, 0, MAX_TEXT_LEN * sizeof(char));
        if (event.type == FocusIn) {
            cl = x11GetFocusWindowClassByDpy(dpy);
            if (cl != NULL) {
                strncpy(text, cl, MAX_TEXT_LEN - 1);
                XFree(cl);
            }
        }
        XSync(dpy, 0);
        XGetInputFocus(dpy, &w, &revertTo);
        if (XFetchName(dpy, w, &name) != 0 && name != NULL) {
            XFree(name);
            name = NULL;
        }
        XSelectInput(dpy, w, FocusChangeMask);
    }
    input_watching = 0;
    th_count--;
    free(text);
    XCloseDisplay(dpy);
    return NULL;
}

void x11StartWindowInspector()
{
    setbuf(stdout, NULL);
    setbuf(stderr, NULL);
    if (input_watching || th_count) {
        return;
    }
    XInitThreads();
    input_watching = 1;
    th_count++;
    pthread_create(&th_input_watch, NULL, &thread_input_watching, NULL);
    pthread_detach(th_input_watch);
}

void x11StopWindowInspector() {
    input_watching = 0;
}
