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

#include <string.h> // strlen
#include <X11/Xlib.h>
#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#define MAX_TEXT_LEN 100

static pthread_t th_clipboard;
static int clipboard_running;
static char * text = NULL;
static char * old_text = NULL;
static int done = 0;

Atom targets_atom, text_atom, UTF8, XA_ATOM = 4, XA_STRING = 31;

void clipboard_exit()
{
    clipboard_running = 0;
    //free(text);
}

static void* thread_clipboard_copy(void* data) {
    Display* display = XOpenDisplay(0);
    if (!display) {
        return NULL;
    }
    int N = DefaultScreen(display);
    Window window = XCreateSimpleWindow(display, RootWindow(display, N), 0, 0, 1, 1, 0,
        BlackPixel(display, N), WhitePixel(display, N));
    targets_atom = XInternAtom(display, "TARGETS", 0);
    text_atom = XInternAtom(display, "TEXT", 0);
    UTF8 = XInternAtom(display, "UTF8_STRING", 1);
    if (UTF8 == None) UTF8 = XA_STRING;
    Atom selection = XInternAtom(display, (char*)data, 0);

    XEvent event;
    Window owner;
    XSetSelectionOwner (display, selection, window, 0);
    if (XGetSelectionOwner (display, selection) != window) return NULL;
    while (clipboard_running==1) {
        XNextEvent (display, &event);
        switch (event.type) {
            case SelectionRequest:
                if (event.xselectionrequest.selection != selection) break;
                XSelectionRequestEvent * xsr = &event.xselectionrequest;
                XSelectionEvent ev = {0};
                int R = 0;
                ev.type = SelectionNotify, ev.display = xsr->display, ev.requestor = xsr->requestor,
                ev.selection = xsr->selection, ev.time = xsr->time, ev.target = xsr->target, ev.property = xsr->property;

                if (done == -1) {
                    ev.property = None;
                    XSendEvent (display, ev.requestor, 0, 0, (XEvent *)&ev);
                    break;
                }
                if (text == NULL) break;
                int size = strlen(text);
                if (ev.target == targets_atom) {
                    R = XChangeProperty (ev.display, ev.requestor, ev.property, XA_ATOM, 32, PropModeReplace, (unsigned char*)&UTF8, 1);
                }
                else if (ev.target == XA_STRING || ev.target == text_atom) {
                    R = XChangeProperty(ev.display, ev.requestor, ev.property, XA_STRING, 8, PropModeReplace, (unsigned char*)text, size);
                    done = 1;
                }
                else if (ev.target == UTF8) {
                    R = XChangeProperty(ev.display, ev.requestor, ev.property, UTF8, 8, PropModeReplace, (unsigned char*)text, size);
                    done = 1;
                }
                else ev.property = None;
                if ((R & 2) == 0) XSendEvent (display, ev.requestor, 0, 0, (XEvent *)&ev);
                break;
            case SelectionClear:
                clipboard_exit();
                break;
        }
    }
    XCloseDisplay(display);
    return NULL;
}

void clipboard_init()
{
    if (clipboard_running==1) {
        return;
    }
    XInitThreads();
    clipboard_running = 1;
    char *selections[] = {"PRIMARY", "CLIPBOARD"};
    for (int i=0; i<sizeof selections/sizeof selections[0]; i++) {
        pthread_create(&th_clipboard, NULL, &thread_clipboard_copy, selections[i]);
        pthread_detach(th_clipboard);
    }
}

void x11Copy(char *str) {
    if (text == NULL) {
        text = (char*)calloc(MAX_TEXT_LEN, sizeof(char));
    }
    strcpy(text, str);
    done = 0;
    fprintf(stderr, "...x11Clipboard text=%s\n", text);
}
