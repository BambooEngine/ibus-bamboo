/* cc-keyboard-shortcut-editor.h
 *
 * Copyright (C) 2016 Endless, Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
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
 * Authors: Georges Basile Stavracas Neto <georges.stavracas@gmail.com>
 */

#pragma once

#include <gtk/gtk.h>

G_BEGIN_DECLS

#define CC_TYPE_KEYBOARD_SHORTCUT_EDITOR (cc_keyboard_shortcut_editor_get_type ())
G_DECLARE_FINAL_TYPE (CcKeyboardShortcutEditor, cc_keyboard_shortcut_editor, CC, KEYBOARD_SHORTCUT_EDITOR, GtkDialog)

typedef struct
{
  guint           keyval;
  guint           keycode;
  GdkModifierType mask;
} CcKeyCombo;
typedef enum
{
  CC_SHORTCUT_EDITOR_CREATE,
  CC_SHORTCUT_EDITOR_EDIT
} CcShortcutEditorMode;

GtkWidget*           cc_keyboard_shortcut_editor_new             ();

// CcKeyboardItem*      cc_keyboard_shortcut_editor_get_item        (CcKeyboardShortcutEditor *self);

void                 cc_keyboard_shortcut_editor_set_item        (CcKeyboardShortcutEditor *self
                                                                  );

CcShortcutEditorMode cc_keyboard_shortcut_editor_get_mode        (CcKeyboardShortcutEditor *self);

void                 cc_keyboard_shortcut_editor_set_mode        (CcKeyboardShortcutEditor *self,
                                                                  CcShortcutEditorMode      mode);

G_END_DECLS

