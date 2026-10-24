#include<stdio.h>
#include<stdlib.h>
#include<string.h>

typedef struct {
    void* ptr;
    int len, cap;
} slice;

// shares memory, only one active at a time
union Data{
    int len;
    char *buffer;
};


int32_t main() {

    // printf("%lu\n", sizeof root);
    // printf("%lu\n", sizeof *root);

    slice x;
    x.ptr = NULL;
    x.len = 1;
    x.cap = 2;

    slice *y = &x;

    if (y->len == y->cap) {
        y->cap *= 2;
    }

    printf("%d : %d\n", y->len, x.len);

    printf("%lu\n", sizeof(slice));
    printf("%lu\n", sizeof(x));
    printf("%lu\n", sizeof(x.cap));
    printf("%lu\n", sizeof(x.len));
    printf("%lu\n", sizeof(x.ptr));

    return 0;
}