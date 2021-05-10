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

#include <glib-object.h>
#include <glib/gi18n.h>

#include "keyboard-shortcut-editor.h"

#include "keyboard.c"
struct _CcKeyboardShortcutEditor
{
  GtkDialog           parent;

  GtkButton          *cancel_button;
  GtkHeaderBar       *headerbar;
  GtkButton          *revert_button;
  GtkButton          *clear_button;
  GtkButton          *set_button;
  GtkLabel           *top_info_label;

  /* Custom shortcuts */
  GtkBox             *standard_box;
  GtkBox             *edit_box;
  GtkStack           *stack;
  GtkLabel   *shortcut_accel_label;
  guint               grab_idle_id;

  CcKeyCombo         *custom_combo;
};

static void          set_button_clicked_cb                       (CcKeyboardShortcutEditor *self);

G_DEFINE_TYPE (CcKeyboardShortcutEditor, cc_keyboard_shortcut_editor, GTK_TYPE_DIALOG)


static void
setup_keyboard_item(CcKeyboardShortcutEditor *self)
{
  gtk_label_set_text (self->top_info_label, "Enter the new shortcut");

  /* Headerbar */
  gtk_header_bar_set_title (self->headerbar,
                            ("Set Shortcut"));

  gtk_widget_hide (GTK_WIDGET (self->set_button));
  gtk_widget_hide (GTK_WIDGET (self->revert_button));

  gtk_label_set_markup (self->top_info_label, "Enter the new shortcut");

  /* Show the appropriate view */
  gtk_stack_set_visible_child (self->stack, GTK_WIDGET (self->edit_box));
}
static void
clear_custom_entries (CcKeyboardShortcutEditor *self)
{
  memset (self->custom_combo, 0, sizeof (CcKeyCombo));
}

static void
clear_button_clicked_cb (CcKeyboardShortcutEditor *self)
{
  gtk_stack_set_visible_child (self->stack, GTK_WIDGET (self->edit_box));
  gtk_widget_hide (GTK_WIDGET (self->set_button));
  gtk_widget_hide (GTK_WIDGET (self->revert_button));
}

static void
cancel_editing (CcKeyboardShortcutEditor *self)
{
  clear_custom_entries (self);

  gtk_widget_destroy (GTK_WIDGET (self));
}

static void
cancel_button_clicked_cb (GtkWidget                *button,
                          CcKeyboardShortcutEditor *self)
{
  cancel_editing (self);
}

static void
revert_button_clicked_cb (CcKeyboardShortcutEditor *self)
{
  gtk_widget_destroy (GTK_WIDGET (self));
}


static void
set_button_clicked_cb (CcKeyboardShortcutEditor *self)
{
  CcKeyCombo *combo;
  combo = self->custom_combo;
  printf("%d,%d", combo->keyval, combo->mask);
  fflush(stdout);
  gtk_widget_destroy (GTK_WIDGET (self));
}

static void
cc_keyboard_shortcut_editor_finalize (GObject *object)
{
  CcKeyboardShortcutEditor *self = (CcKeyboardShortcutEditor *)object;

  g_clear_pointer (&self->custom_combo, g_free);

  G_OBJECT_CLASS (cc_keyboard_shortcut_editor_parent_class)->finalize (object);
}

static void
cc_keyboard_shortcut_editor_get_property (GObject    *object,
                                          guint       prop_id,
                                          GValue     *value,
                                          GParamSpec *pspec)
{
  CcKeyboardShortcutEditor *self = CC_KEYBOARD_SHORTCUT_EDITOR (object);

  switch (prop_id)
    {
    default:
      G_OBJECT_WARN_INVALID_PROPERTY_ID (object, prop_id, pspec);
    }
}

static void
cc_keyboard_shortcut_editor_set_property (GObject      *object,
                                          guint         prop_id,
                                          const GValue *value,
                                          GParamSpec   *pspec)
{
  CcKeyboardShortcutEditor *self = CC_KEYBOARD_SHORTCUT_EDITOR (object);

  switch (prop_id)
    {
    default:
      G_OBJECT_WARN_INVALID_PROPERTY_ID (object, prop_id, pspec);
    }
}

static gboolean
cc_keyboard_shortcut_editor_key_press_event (GtkWidget   *widget,
                                             GdkEventKey *event)
{
  CcKeyboardShortcutEditor *self;
  GdkModifierType real_mask;
  gboolean editing;
  guint keyval_lower;

  self = CC_KEYBOARD_SHORTCUT_EDITOR (widget);

  real_mask = event->state;
  /* real_mask = event->state & gtk_accelerator_get_default_mod_mask (); */

  /* keyval_lower = gdk_keyval_to_lower (event->keyval); */
  keyval_lower = event->keyval;

  /* Normalise <Tab> */
  if (keyval_lower == GDK_KEY_ISO_Left_Tab)
    keyval_lower = GDK_KEY_Tab;

  /* Put shift back if it changed the case of the key, not otherwise. */
  if (keyval_lower != event->keyval)
    real_mask |= GDK_SHIFT_MASK;

  if (keyval_lower == GDK_KEY_Sys_Req &&
      (real_mask & GDK_MOD1_MASK) != 0)
    {
      /* HACK: we don't want to use SysRq as a keybinding (but we do
       * want Alt+Print), so we avoid translation from Alt+Print to SysRq */
      keyval_lower = GDK_KEY_Print;
    }

  /* A single Escape press cancels the editing */
  if (!event->is_modifier && real_mask == 0 && keyval_lower == GDK_KEY_Escape)
    {

      cancel_editing (self);

      return GDK_EVENT_STOP;
    }

  self->custom_combo->keycode = event->hardware_keycode;
  self->custom_combo->keyval = keyval_lower;
  self->custom_combo->mask = real_mask;

  /* CapsLock isn't supported as a keybinding modifier, so keep it from confusing us */
  self->custom_combo->mask &= ~GDK_LOCK_MASK;

  GtkWidget *label;

  g_autofree gchar *accel = NULL;

  #if GTK_MAJOR_VERSION <= 3 && GTK_MINOR_VERSION < 24
    accel = gtk_accelerator_get_label_with_keycode (NULL, self->custom_combo->keyval, self->custom_combo->keycode, self->custom_combo->mask);
    gtk_label_set_text (self->shortcut_accel_label, accel);
  #else
    accel = gtk_accelerator_name (keyval_lower, real_mask);
    gtk_shortcut_label_set_accelerator (GTK_SHORTCUT_LABEL(self->shortcut_accel_label), accel);
  #endif
  gtk_stack_set_visible_child (self->stack, GTK_WIDGET (self->standard_box));
  gtk_widget_show (GTK_WIDGET (self->set_button));
  /* gtk_widget_show (GTK_WIDGET (self->revert_button)); */
  return GDK_EVENT_STOP;
}

static void
cc_keyboard_shortcut_editor_close (GtkDialog *dialog)
{
  CcKeyboardShortcutEditor *self = CC_KEYBOARD_SHORTCUT_EDITOR (dialog);

  GTK_DIALOG_CLASS (cc_keyboard_shortcut_editor_parent_class)->close (dialog);
}

static void
cc_keyboard_shortcut_editor_response (GtkDialog *dialog,
                                      gint       response_id)
{
  CcKeyboardShortcutEditor *self = CC_KEYBOARD_SHORTCUT_EDITOR (dialog);
}

static gboolean
grab_idle (gpointer data)
{
  CcKeyboardShortcutEditor *self = data;

  self->grab_idle_id = 0;

  return G_SOURCE_REMOVE;
}

static void
cc_keyboard_shortcut_editor_show (GtkWidget *widget)
{
  CcKeyboardShortcutEditor *self = CC_KEYBOARD_SHORTCUT_EDITOR (widget);

  /* Map before grabbing, so that the window is visible */
  GTK_WIDGET_CLASS (cc_keyboard_shortcut_editor_parent_class)->show (widget);
  setup_keyboard_item(self);

  self->grab_idle_id = g_timeout_add (100, grab_idle, self);
}

static void
cc_keyboard_shortcut_editor_unrealize (GtkWidget *widget)
{
  CcKeyboardShortcutEditor *self = CC_KEYBOARD_SHORTCUT_EDITOR (widget);

  if (self->grab_idle_id) {
    g_source_remove (self->grab_idle_id);
    self->grab_idle_id = 0;
  }

  GTK_WIDGET_CLASS (cc_keyboard_shortcut_editor_parent_class)->unrealize (widget);
}

static void
cc_keyboard_shortcut_editor_class_init (CcKeyboardShortcutEditorClass *klass)
{
  GtkWidgetClass *widget_class = GTK_WIDGET_CLASS (klass);
  GtkDialogClass *dialog_class = GTK_DIALOG_CLASS (klass);
  GObjectClass *object_class = G_OBJECT_CLASS (klass);

  object_class->finalize = cc_keyboard_shortcut_editor_finalize;
  object_class->get_property = cc_keyboard_shortcut_editor_get_property;
  object_class->set_property = cc_keyboard_shortcut_editor_set_property;

  widget_class->show = cc_keyboard_shortcut_editor_show;
  widget_class->unrealize = cc_keyboard_shortcut_editor_unrealize;
  widget_class->key_press_event = cc_keyboard_shortcut_editor_key_press_event;

  dialog_class->close = cc_keyboard_shortcut_editor_close;
  dialog_class->response = cc_keyboard_shortcut_editor_response;
#if GTK_MAJOR_VERSION <= 3 && GTK_MINOR_VERSION < 24
  gtk_widget_class_set_template_from_resource (widget_class, "/org/input/bamboo/setup-ui/keyboard-shortcut-editor-v3.ui");
#else
  gtk_widget_class_set_template_from_resource (widget_class, "/org/input/bamboo/setup-ui/keyboard-shortcut-editor.ui");
#endif

  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, cancel_button);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, headerbar);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, revert_button);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, clear_button);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, set_button);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, top_info_label);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, shortcut_accel_label);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, standard_box);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, stack);
  gtk_widget_class_bind_template_child (widget_class, CcKeyboardShortcutEditor, edit_box);

  gtk_widget_class_bind_template_callback (widget_class, clear_button_clicked_cb);
  gtk_widget_class_bind_template_callback (widget_class, cancel_button_clicked_cb);
  gtk_widget_class_bind_template_callback (widget_class, revert_button_clicked_cb);
  gtk_widget_class_bind_template_callback (widget_class, set_button_clicked_cb);
}

static void
cc_keyboard_shortcut_editor_init (CcKeyboardShortcutEditor *self)
{
  gtk_widget_init_template (GTK_WIDGET (self));

  self->custom_combo = g_new0 (CcKeyCombo, 1);
}

GtkWidget*
cc_keyboard_shortcut_editor_new ()
{
  return g_object_new (CC_TYPE_KEYBOARD_SHORTCUT_EDITOR,
                       "use-header-bar", 1,
                       NULL);
}

int main(int argc, char *argv[]) {
  GtkWidget *window;
  gtk_init(&argc, &argv);
  window = cc_keyboard_shortcut_editor_new();

  gtk_widget_show (window);

  g_signal_connect(window, "destroy",
      G_CALLBACK(gtk_main_quit), NULL);

  gtk_main();
  return 0;
}
