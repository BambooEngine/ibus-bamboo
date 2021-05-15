#include <stdio.h>

FILE * open_for_read(char *fname) {
    FILE *fh;
    if ((fh = fopen(fname,"rb")) == NULL) {
        exit(1);
    }
    return fh;
}

FILE * open_for_write(char *fname) {
    FILE *fh;
    if ((fh = fopen(fname,"wb")) == NULL) {
        exit(1);
    }
    return fh;
}
