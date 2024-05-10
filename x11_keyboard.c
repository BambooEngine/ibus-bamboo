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

#include <sys/time.h>
#include <time.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/keysym.h> //xproto-devel
#include <X11/keysymdef.h>
#include <X11/extensions/XTest.h>

static void delay(int sec, long msec) {
    long pause;
    clock_t now,then;

    pause = msec*(CLOCKS_PER_SEC/1000);
    now = then = clock();
    while( (now-then) < pause )
        now = clock();
}

void x11SendShiftR() {
    Display *display = XOpenDisplay(NULL);
    if (display) {
        KeyCode xk_shift_r = XKeysymToKeycode(display, XK_Shift_R);
        XTestFakeKeyEvent(display, xk_shift_r, False, 0);
        XCloseDisplay(display);
    }
}

void x11SendShiftLeft(int n, int r, int timeout) {
    Display *display = XOpenDisplay(NULL);
    if (display) {
        XSynchronize(display, 1);
        KeyCode modcode;
        KeyCode xk_shift_l = XKeysymToKeycode(display, XK_Shift_L);
        KeyCode xk_shift_r = XKeysymToKeycode(display, XK_Shift_R);
        modcode = XKeysymToKeycode(display, XStringToKeysym("Left"));
        switch (r) {
        case 1:
            XTestFakeKeyEvent(display, xk_shift_l, True, 0);
            for (int i=0; i<n; i++) {
                XTestFakeKeyEvent(display, modcode, True, 0);
                XTestFakeKeyEvent(display, modcode, False, 0);
            }
            XTestFakeKeyEvent(display, xk_shift_l, False, 0);
            break;
        default:
            XTestFakeKeyEvent(display, xk_shift_r, True, 0);
            for (int i=0; i<n; i++) {
                XTestFakeKeyEvent(display, modcode, True, 0);
                XTestFakeKeyEvent(display, modcode, False, 0);
            }
            XTestFakeKeyEvent(display, xk_shift_r, False, 0);
            break;
        }
        XSynchronize(display, 0);
        XCloseDisplay(display);
    }
}

void x11SendBackspace(int n, int timeout) {
    Display *display = XOpenDisplay(NULL);
    if (display) {
        /* XSynchronize(display, 1); */
        KeyCode modcode;
        modcode = XKeysymToKeycode(display, XStringToKeysym("BackSpace"));
        for (int i=0; i<n; i++) {
            XTestFakeKeyEvent(display, modcode, True, 0);
            XTestFakeKeyEvent(display, modcode, False, 0);
            /* XSync(display, 0); */
            delay(0, timeout);
        }
        XFlush(display);
        /* XSynchronize(display, 0); */
        XCloseDisplay(display);
    }
}

void x11Paste(int n) {
    Display *display = XOpenDisplay(NULL);
    if (display) {
        KeyCode xk_shift_l = XKeysymToKeycode(display, XK_Shift_L);
        KeyCode xk_shift_r = XKeysymToKeycode(display, XK_Shift_R);
        KeyCode xk_control = XKeysymToKeycode(display, XK_Control_L);
        KeyCode xk_insert = XKeysymToKeycode(display, XK_Insert);
        KeyCode xk_v = XKeysymToKeycode(display, XK_V);

        switch (n) {
        case 0:
            XTestFakeKeyEvent(display, xk_shift_l, True, 0);
            XTestFakeKeyEvent(display, xk_insert, True, 0);
            XTestFakeKeyEvent(display, xk_shift_l, False, 0);
            XTestFakeKeyEvent(display, xk_insert, False, 0);
            break;
        case 1:
            XTestFakeKeyEvent(display, xk_shift_r, True, 0);
            XTestFakeKeyEvent(display, xk_insert, True, 0);
            XTestFakeKeyEvent(display, xk_shift_r, False, 0);
            XTestFakeKeyEvent(display, xk_insert, False, 0);
            break;
        case 2:
            XTestFakeKeyEvent(display, xk_control, True, 0);
            XTestFakeKeyEvent(display, xk_v, True, 0);
            XTestFakeKeyEvent(display, xk_control, False, 0);
            XTestFakeKeyEvent(display, xk_v, False, 0);
            break;
        }
        XSync(display, 0);
        XCloseDisplay(display);
    }
}

void x11SendString(char* str) {
    Display *display = XOpenDisplay(NULL);
    if (display) {
        for (int i=0; i<strlen(str); i++) {
            char chr[2] = {str[i], '\0'};
            int modcode = XKeysymToKeycode(display, XStringToKeysym(chr));
            XTestFakeKeyEvent(display, modcode, False, 0);
            XTestFakeKeyEvent(display, modcode, True, 0);
            XTestFakeKeyEvent(display, modcode, False, 0);
            XSync(display, 0);
        }
        XFlush(display);
        XCloseDisplay(display);
    }
}
