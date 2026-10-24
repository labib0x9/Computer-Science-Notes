#include<stdio.h>
#include<stdbool.h>
#include<stddef.h>
#include<string.h>

struct A {
    bool a;
    int b;
    char c;
    long long d;
    char e;
};


struct B {
    bool a;
    int b;
    char c;
    long long d;
    char e;
}__attribute__((packed));


int main() {

    printf("%zu\n", sizeof(struct A));
    printf("%zu\n", sizeof(struct B));

    printf("%zu\n", offsetof(struct A, b));
    printf("%zu\n", offsetof(struct B, b));

    struct A first_st = {
        .a = true,
        .b = 1023,
        .c = 'c',
        .d = 1232ll,
        .e = 'e'
    };

    int b;
    memmove(&b, (char *)&first_st + offsetof(struct A, b), sizeof(b));

    printf("%d\n", b);

    return 0;
}