// gcc -o xtest xtest.c -lX11 -lXtst
#define _GNU_SOURCE
#include <sys/time.h>
#include <time.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/keysym.h> //xproto-devel
#include <X11/extensions/XTest.h>


/**
 * milliseconds over 1000 will be ignored
 */
static void delay(time_t sec, long msec) {
    struct timespec sleep;

    sleep.tv_sec  = sec;
    sleep.tv_nsec = (msec % 1000) * 1000 * 1000;

    if (nanosleep(&sleep, NULL) == -1) {
    }
}

void SendKeys(long keys[], int len) {
	Display *display = XOpenDisplay(NULL);
	if (display) {
		//find out window with current focus:
        Window winfocus;
        int    revert, modcode;
        XGetInputFocus(display, &winfocus, &revert);

		for (;;) {
            for (int i=0;i<len;i++) {
                modcode = XKeysymToKeycode(display, keys[i]);;
                // key press
                XTestFakeKeyEvent(display, modcode, True, 0);
                XSync(display, 0);
                // key release
                XTestFakeKeyEvent(display, modcode, False, 0);
                XSync(display, 0);
                delay(0, 50);
            }
		}
		XSync(display, 1);
		XCloseDisplay(display);
	}
}

void str_copy(char* dest, char* src, int start, int end) {
    for (int i=start;i<end;i++) {
        dest[i] = src[i];
    }
}

int main() {
    FILE* sample = NULL;
    sample = fopen("xtest.data", "r+");
    int chr = 0;
    int len = 0;
    char command[1000] = "";
    long xks[100000];
    if (sample != NULL) {
        do {
            chr = fgetc(sample);
            if (chr == ' ' || chr == '\n') {
                xks[len] = XStringToKeysym("space");
            } else if (chr == '\\') {
                fgets(command, 1000, sample);
                char keyname[1000] = "";
                str_copy(keyname, command, 1, strlen(command));
                xks[len] = XStringToKeysym(keyname);
            } else {
                char str[2] = {chr, '\0'};
                xks[len] = XStringToKeysym(str);
            }
            len += 1;
        } while (chr != EOF);
        fclose(sample);
    }
    SendKeys(xks, len);
    return 0;
}
