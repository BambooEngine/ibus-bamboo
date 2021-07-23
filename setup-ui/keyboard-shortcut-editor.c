#include <ctype.h>
#include <gtk/gtk.h>

#define TOTAL_ROWS 5
#define TOTAL_MASKS_PER_ROW 4
int row = 0;
int col = 0;
const int KEYVAL = 1;
const int MASK = 0;
int key_pairs[TOTAL_ROWS * 2];
char *labels[TOTAL_MASKS_PER_ROW] = {"Ctrl", "Alt", "Shift", "Super"};
int masks[TOTAL_MASKS_PER_ROW] = {GDK_CONTROL_MASK, GDK_MOD1_MASK, GDK_SHIFT_MASK,
                         GDK_SUPER_MASK};
int keyvals[TOTAL_MASKS_PER_ROW] = {GDK_KEY_Control_L, GDK_KEY_Alt_L, GDK_KEY_Shift_L,
                           GDK_KEY_Super_L};
char *text_arr[TOTAL_ROWS] = {"Chuyển chế độ gõ", "Khôi phục phím",
                                "Tạm tắt bộ gõ", "Emoji", "Hexadecimal"};
GtkWidget *maskWidgets[TOTAL_MASKS_PER_ROW * TOTAL_ROWS];
GtkWidget *keyWidgets[TOTAL_ROWS];

/*
 * Destroy
 *
 * Close down the application
 */
gint close_window_cb(GtkWidget *widget, gpointer *data) {
  gtk_main_quit();
  return FALSE;
}

gint btn_reset_cb(GtkWidget *widget, gpointer *data) {
  for (int i = 0 ; i < TOTAL_ROWS * 2; i++ ){
    key_pairs[i] = 0;
  }
  for (int i=0 ; i < TOTAL_ROWS * TOTAL_MASKS_PER_ROW ; i++) {
    gtk_toggle_button_set_active(GTK_TOGGLE_BUTTON(maskWidgets[i]), FALSE);
  }
  for (int i=0 ; i < TOTAL_ROWS ; i++) {
    gtk_entry_set_text(GTK_ENTRY(keyWidgets[i]), "");
  }
  return FALSE;
}

/*
 * btn_save_cb
 *
 * Some event happened and the name is passed in the
 * data field.
 */
void btn_save_cb(GtkWidget *widget, gpointer data) {
  for (int i = 0; i < TOTAL_ROWS * 2 - 1; i++) {
    printf("%d,", key_pairs[i]);
  }
  printf("%d", key_pairs[TOTAL_ROWS * 2 - 1]);
  fflush(stdout);
  close_window_cb(widget, data);
}

/*
 * check_event_cb
 *
 * Handle a checkbox signal
 */
void check_event_cb(GtkWidget *widget, gpointer data) {
  int pos = GPOINTER_TO_INT(data);
  int row = pos / TOTAL_MASKS_PER_ROW, mask_col = pos % TOTAL_MASKS_PER_ROW;
  if (gtk_toggle_button_get_active(GTK_TOGGLE_BUTTON(widget))) {
    key_pairs[row * 2] |= masks[mask_col];
  } else {
    key_pairs[row * 2] &= ~masks[mask_col];
  }
}

char * int_to_accel(int keyval) {
  gchar *accel = NULL;
  accel = gtk_accelerator_get_label(keyval, 0);

  // Convert to upper case
  char *s = accel;
  while (*s) {
    *s = toupper((unsigned char)*s);
    s++;
  }
  return accel;
}

static gboolean key_release_cb(GtkWidget *entry, GdkEventKey *event,
                           gpointer data) {
  int row = GPOINTER_TO_INT(data);
  int keyval = key_pairs[row * 2 + 1];

  /* --- Put text in the field. --- */
  gtk_entry_set_text(GTK_ENTRY(entry), int_to_accel(keyval));
  return TRUE;
}

static gboolean key_press_cb(GtkWidget *entry, GdkEventKey *event, gpointer data) {
  int row = GPOINTER_TO_INT(data);
  if (event->keyval == GDK_KEY_BackSpace || event->keyval == GDK_KEY_Delete) {
    key_pairs[row * 2 + 1] = 0;
    return FALSE;
  }
  key_pairs[row * 2 + 1] = gdk_keyval_to_lower(event->keyval);
  return TRUE;
}

void add_checkbox(GtkWidget *parent, char *text, int mask_pos) {
  // GtkWidget *check;
  int pad = 10;
  /*
   * --- Create a check button
   */
  maskWidgets[mask_pos] = gtk_check_button_new_with_label(text);
  /*
   * --- Active/Inactive check button
   */
  int row = mask_pos / TOTAL_MASKS_PER_ROW, mask_col = mask_pos % TOTAL_MASKS_PER_ROW;
  int mask = key_pairs[row * 2];
  gboolean active = FALSE;
  if (mask&masks[mask_col]) {
    active = TRUE;
  }
  gtk_toggle_button_set_active(GTK_TOGGLE_BUTTON(maskWidgets[mask_pos]), active);

  /* --- Pack the checkbox into the parent (expand? fill? padding?).  --- */
  gtk_box_pack_start(GTK_BOX(parent), maskWidgets[mask_pos], FALSE, FALSE, pad);

  g_signal_connect(maskWidgets[mask_pos], "toggled", G_CALLBACK(check_event_cb),
                   GINT_TO_POINTER(mask_pos));
}

void add_shortcut_box(GtkWidget *widget, char *text, int row) {
  GtkWidget *hbox, *label_hbox;
  GtkWidget *label;
  int pad = 10;
  /* Horizontal box to pack shortcut and label */
  hbox = gtk_box_new(GTK_ORIENTATION_HORIZONTAL, 0);

  /* Horizontal box to pack label */
  label_hbox = gtk_box_new(GTK_ORIENTATION_HORIZONTAL, 0);
  /* --- create a new label.  --- */
  label = gtk_label_new(text);
  gtk_label_set_xalign(GTK_LABEL(label), 0);
  /* --- Pack the label into the horizontal box (expand? fill? padding)  --- */
  gtk_box_pack_start(GTK_BOX(hbox), label, TRUE, TRUE, pad);

  for (int i = 0; i < TOTAL_MASKS_PER_ROW; i++) {
    add_checkbox(hbox, labels[i], row * TOTAL_MASKS_PER_ROW + i);
  }

  /* --- Create an entry field --- */
  keyWidgets[row] = gtk_entry_new();
  GtkWidget *entry = keyWidgets[row];

  /* --- Pack the entry into the vertical box (expand? fill?, padding?).  --- */
  gtk_box_pack_start(GTK_BOX(hbox), entry, FALSE, FALSE, 10);

  /* --- Put some text in the field. --- */
  int kvl = gdk_keyval_to_lower(key_pairs[row*2+1]);
  gtk_entry_set_text(GTK_ENTRY(entry), int_to_accel(kvl));
  gtk_entry_set_alignment(GTK_ENTRY(entry), 0.5);

  /* --- Pack it in. --- */
  gtk_box_pack_start(GTK_BOX(widget), hbox, FALSE, FALSE, 0);

  g_signal_connect(entry, "key_press_event", G_CALLBACK(key_press_cb),
                   GINT_TO_POINTER(row));
  g_signal_connect(entry, "key_release_event", G_CALLBACK(key_release_cb),
                   GINT_TO_POINTER(row));
}

void add_control_buttons(GtkWidget *widget) {
  GtkWidget *save_button;
  GtkWidget *cancel_button;
  GtkWidget *reset_button;
  GtkWidget *hbox;

  /* Horizontal box to pack OK and Cancel buttons */
  hbox = gtk_box_new(GTK_ORIENTATION_HORIZONTAL, 0);
  gtk_widget_set_halign(hbox, GTK_ALIGN_END);

  /* --- Create a Reset button. --- */
  reset_button = gtk_button_new_with_label("Reset");

  /* --- Pack the reset_button into the vertical box (vbox box1).  --- */
  gtk_box_pack_start(GTK_BOX(hbox), reset_button, FALSE, FALSE, 10);

  /* --- Create a Cancel button. --- */
  cancel_button = gtk_button_new_with_label("Cancel");

  /* --- Pack the cancel_button into the vertical box (vbox box1).  --- */
  gtk_box_pack_start(GTK_BOX(hbox), cancel_button, FALSE, FALSE, 10);

  /* --- Create a Save button. --- */
  save_button = gtk_button_new_with_label("Save");

  /* --- Pack the button into the vertical box (vbox box1).  --- */
  gtk_box_pack_start(GTK_BOX(hbox), save_button, FALSE, FALSE, 10);

  gtk_container_add(GTK_CONTAINER(widget), hbox);

  g_signal_connect(reset_button, "clicked", G_CALLBACK(btn_reset_cb), "clicked");
  g_signal_connect(save_button, "clicked", G_CALLBACK(btn_save_cb), "clicked");
  g_signal_connect(cancel_button, "clicked", G_CALLBACK(close_window_cb),
                   "clicked");
}

void shortcut_init(char *s) {
  int count = 0, n = 0;
  while (*s) {
    if (*s == ',') {
      key_pairs[count] = n;
      count++;
      n = 0;
    } else if (*s >= '0' || *s <= '9') {
      int t = *s - '0';
      n = n * 10 + t;
    }
    s++;
  }
  key_pairs[count] = n;
}

/*
 * Main - program begins here
 */
int main(int argc, char *argv[]) {
  GtkWidget *window;
  GtkWidget *dialog, *content_area;
  GtkWidget *vbox, *vcbox;
  int which;
  int pad = 15;

  /* --- GTK initialization --- */
  gtk_init(&argc, &argv);
  if (argc > 1) {
    shortcut_init(argv[1]);
  }

  /* --- Create the top level window --- */
  window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
  dialog = gtk_dialog_new();
  content_area = gtk_dialog_get_content_area(GTK_DIALOG(dialog));
  gtk_window_set_transient_for(GTK_WINDOW(dialog), GTK_WINDOW(window));
  gtk_widget_set_size_request(dialog, 600, 250);

  /* --- You should always remember to connect the delete_event
   *     to the main window.
   */
  g_signal_connect(window, "delete_event", G_CALLBACK(close_window_cb), NULL);

  /* --- Give the window a border --- */
  gtk_container_set_border_width(GTK_CONTAINER(content_area), pad);

  /* --- We create a vertical box (vbox) to pack
   *     the horizontal boxes into.
   */
  vbox = gtk_box_new(GTK_ORIENTATION_VERTICAL, pad);

  for (int i = 0; i < TOTAL_ROWS; i++) {
    add_shortcut_box(vbox, text_arr[i], i);
  }

  vcbox = gtk_box_new(GTK_ORIENTATION_VERTICAL, pad);
  add_control_buttons(vcbox);

  /* --- Align the controls box to the bottom.   --- */
  gtk_widget_set_valign(vcbox, GTK_ALIGN_END);
  gtk_widget_set_vexpand(vcbox, TRUE);
  gtk_box_pack_start(GTK_BOX(vbox), vcbox, TRUE, TRUE, 0);

  /*
   * --- Make the main window visible
   */
  gtk_container_add(GTK_CONTAINER(content_area), vbox);
  gtk_window_set_title(GTK_WINDOW(dialog), "IBus Bamboo Shortcuts");
  gtk_widget_show_all(dialog);

  gtk_main();
  exit(0);
}
