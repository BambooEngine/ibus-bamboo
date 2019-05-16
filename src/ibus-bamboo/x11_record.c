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
#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <X11/Xlibint.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/cursorfont.h>
#include <X11/keysymdef.h>
#include <X11/keysym.h>
#include <X11/extensions/record.h>
#include <X11/extensions/XTest.h>
#include "_cgo_export.h"
#define CAPTURE_MOUSE_MOVE_DELTA        100

/* for this struct, refer to libxnee */
typedef union {
  unsigned char    type ;
  xEvent           event ;
  xResourceReq     req   ;
  xGenericReply    reply ;
  xError           error ;
  xConnSetupPrefix setup;
} XRecordDatum;

/*
 * FIXME: We need define a private struct for callback function,
 * to store cur_x, cur_y, data_disp, ctrl_disp etc.
 */
Display *data_disp = NULL;
Display *ctrl_disp = NULL;

XRecordRange  *rr = NULL;
XRecordClientSpec  rcs;
XRecordContext rc;

/* recording flag */
static int mouse_recording = 0;

static int cur_x = 0;
static int cur_y = 0;

void event_callback (XPointer, XRecordInterceptData*);

static void* thread_mouse_recording (void* data)
{
  ctrl_disp = XOpenDisplay (NULL);
  data_disp = XOpenDisplay (NULL);

  if (!ctrl_disp || !data_disp) {
    fprintf (stderr, "Error to open local display!\n");
    return NULL;
  }

  /* 
   * we must set the ctrl_disp to sync mode, or, when we the enalbe 
   * context in data_disp, there will be a fatal X error !!!
   */
  XSynchronize(ctrl_disp,True);

  int major, minor;
  if (!XRecordQueryVersion (ctrl_disp, &major, &minor)) {
    fprintf (stderr, "RECORD extension is not supported on this X server!\n");
    mouse_recording = 0;
    return NULL;
  }
 
  printf ("RECORD extension for local server is version %d.%d\n", major, minor);

  rr = XRecordAllocRange ();
  if (!rr) {
    fprintf (stderr, "Could not alloc record range object!\n");
    return NULL;
  }

  rr->device_events.first = KeyPress;
  rr->device_events.last = MotionNotify;
  rcs = XRecordAllClients;

  rc = XRecordCreateContext (ctrl_disp, 0, &rcs, 1, &rr, 1);
  if (!rc) {
    fprintf (stderr, "Could not create a record context!\n");
    return NULL;
  }
 
  if (!XRecordEnableContext (data_disp, rc, event_callback, NULL)) {
    fprintf (stderr, "Cound not enable the record context!\n");
    return NULL;
  }

  Window w_root_return, w_child_return;
  int x_old, y_old, x_root_old, y_root_old, rt;
  unsigned int mask;
  /* Note: you should not use data_disp to do normal X operations !!!*/
  XQueryPointer(ctrl_disp, XDefaultRootWindow(ctrl_disp), &w_root_return, &w_child_return, &x_root_old, &y_root_old, &x_old, &y_old, &mask);
  cur_x = x_root_old;
  cur_y = y_root_old;
  while (mouse_recording) {
    XRecordProcessReplies (data_disp);
  }

  XRecordDisableContext (ctrl_disp, rc);
  XRecordFreeContext (ctrl_disp, rc);
  XFree (rr);
 
  XCloseDisplay (data_disp);
  XCloseDisplay (ctrl_disp);
  mouse_recording = 0;
  return NULL;
}

void event_callback(XPointer priv, XRecordInterceptData *hook)
{

  if (hook->category != XRecordFromServer) {
    XRecordFreeData (hook);
    return;
  }

  XRecordDatum *data = (XRecordDatum*) hook->data;

  int event_type = data->type;

  BYTE btncode, keycode;
  btncode = keycode = data->event.u.u.detail;

  int root_x = data->event.u.keyButtonPointer.rootX;
  int root_y = data->event.u.keyButtonPointer.rootY;
  int time = hook->server_time;

  switch (event_type) {
  case KeyPress:
    /** printf ("KeyPress: \t%s\n", XKeysymToString(XKeycodeToKeysym(ctrl_disp, keycode, 0))); */
    break;
  case KeyRelease:
    /** printf ("KeyRelease: \t%s\n", XKeysymToString(XKeycodeToKeysym(ctrl_disp, keycode, 0))); */
    break;
  case ButtonPress:
    /** printf ("ButtonPress: /t%d, rootX=%d, rootY=%d, recording=%d", btncode, cur_x, cur_y, mouse_recording); */
    mouse_click_handler();
    break;
  case ButtonRelease:
    /** printf ("ButtonRelease: /t%d, rootX=%d, rootY=%d", btncode, cur_x, cur_y); */
    break;
  case MotionNotify:
    /* printf ("MouseMove: /trootX=%d, rootY=%d",rootx, rooty); */
    if ((abs(root_x - cur_x) >= CAPTURE_MOUSE_MOVE_DELTA) ||
        (abs(root_y - cur_y) >= CAPTURE_MOUSE_MOVE_DELTA)) // mouse move at least CAPTURE_MOUSE_MOVE_DELTA
    {
        /** mouse_move_handler(); */
        cur_x = root_x;
        cur_y = root_y;
    }
    break;
  default:
    break;
  }
  XRecordFreeData (hook);
}

void mouse_recording_init()
{
    setbuf(stdout, NULL);
    setbuf(stderr, NULL);
    if (mouse_recording==1) {
        return;
    }
    XInitThreads();
    mouse_recording = 1;
    pthread_t th_mcap;
    pthread_create(&th_mcap, NULL, &thread_mouse_recording, NULL);
    pthread_detach(th_mcap);
}

void mouse_recording_exit()
{
    if (mouse_recording==0 || ctrl_disp == NULL) {
        return;
    }
    XRecordDisableContext (ctrl_disp, rc);
    mouse_recording = 0;
}

