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
#include <X11/Xlib.h>

#define MaxPropertyLen 128
#define MaxWmClassesLen 10
static char * WM_CLASS = "WM_CLASS";
static char * WM_NAME = "WM_NAME";

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
    Atom actualType, filterAtom, XA_STRING = 31, XA_ATOM = 4;
    int status, actualFormat = 0;
    unsigned long len, bytesAfter;
    unsigned char * uc = NULL;

    filterAtom = XInternAtom(display, propName, True);
    status = XGetWindowProperty(display, window, filterAtom, 0, MaxPropertyLen, False, AnyPropertyType,
        &actualType, &actualFormat, &len, &bytesAfter, &uc);
    if (status == Success) {
        return uchar2char(uc, len);
    }
    return NULL;
}

char * x11GetFocusWindowClasses(Display *display, char * propName) {
    Window w;
    int revertTo;
    XGetInputFocus(display, &w, &revertTo);
    for (int i=0; i<MaxWmClassesLen; i++) {
        char * strClass = x11GetStringProperty(display, w, propName);
        if (strClass != NULL && strstr(strClass, "FocusProxy") == NULL) {
            return strClass;
        }
        Window * childrenWindows;
        Window parentWindow, rootWindow;
        unsigned int nChild = 0;
        XQueryTree(display, w, &rootWindow, &parentWindow, &childrenWindows, &nChild);
        if (childrenWindows != NULL) {
            //XFree(childrenWindows);
        }
        if (rootWindow == parentWindow) {
            break;
        }
        w = parentWindow;
    }
    return NULL;
}

char * x11GetFocusWindowClass() {
    Display *display = XOpenDisplay(NULL);
    if (!display) {
        return NULL;
    }
    char * strClass = x11GetFocusWindowClasses(display, WM_CLASS);
    if (strClass == NULL) {
        strClass = x11GetFocusWindowClasses(display, WM_NAME);
    }
    XCloseDisplay(display);
    return strClass;
}
