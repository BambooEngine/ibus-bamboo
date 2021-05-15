#include <gtk/gtk.h>
#include "keyboard.c"
#include "utils.h"

GtkTextBuffer *textBuff;
GtkWidget *textView;

FILE *fh;
char fileBuf[1284];
char *ui_path = "/org/input/bamboo/setup-ui/macro-editor.ui";
int modified = 0;

int main(int argc, char *argv[]) {
    GtkBuilder      *builder;
    GtkWidget       *window;

    gtk_init(&argc, &argv);

    builder = gtk_builder_new();
    gtk_builder_add_from_resource (builder, ui_path, NULL);

    window = GTK_WIDGET(gtk_builder_get_object(builder, "macro_window"));
    gtk_builder_connect_signals(builder, NULL);

    textBuff = GTK_TEXT_BUFFER(gtk_builder_get_object(builder, "textbuff"));
    textView = GTK_WIDGET(gtk_builder_get_object(builder, "textview"));

    g_object_unref(builder);
    gtk_widget_show(window);

    if (argc > 1) {  // load file from command line
        strcpy(fileBuf, argv[1]);
        // read the file and insert it into the textview
        fh = open_for_read(fileBuf);
        fseek(fh, 0, SEEK_END);
        long fsize = ftell(fh);
        fseek(fh, 0, SEEK_SET);
        char *string = malloc(fsize + 1);
        fread(string, 1, fsize, fh);
        fclose(fh);
        string[fsize] = 0;
        gtk_text_buffer_set_text(GTK_TEXT_BUFFER(textBuff), string, -1);
        gtk_window_set_title(GTK_WINDOW(window), fileBuf);
        free(string);
        gtk_text_buffer_set_modified (GTK_TEXT_BUFFER(textBuff) , FALSE);
        modified = 0;
    }

    gtk_main();
    return 0;
}


void macro_editor_close() {
    gtk_main_quit();
}

void macro_editor_save() {
    if (!modified) {
        gtk_main_quit();
    }
    GtkTextIter start, end;
    GtkTextBuffer *buffer = gtk_text_view_get_buffer (GTK_TEXT_VIEW(textView));
    gchar *text;

    gtk_text_buffer_get_bounds (buffer, &start, &end);
    text = gtk_text_buffer_get_text (buffer, &start, &end, FALSE);
    fh = open_for_write(fileBuf);
    fprintf(fh, "%s", text);
    fclose(fh);
    g_free(text);
    gtk_text_buffer_set_modified (buffer , FALSE);
    modified = 0;
    gtk_main_quit();
}

void on_textbuff_modified_changed() {
    modified = 1;
}

