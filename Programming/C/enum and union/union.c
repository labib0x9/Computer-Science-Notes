#include<stdio.h>
#include<stdlib.h>
#include<stdint.h>

// Use the same memory. allocates the largeset fileds size. 
typedef union SAME_MEMORY {
    int val;        // 4byte
    long long err;  // 8byte
} union_type_t;     // max(8, 4) = 8byte

// helper-filed, access different parts of the memory.
typedef union Color {
    struct {
        uint8_t r;
        uint8_t g;
        uint8_t b;
        uint8_t a;
    } components;
    uint32_t rgba;
} color_t;

int32_t main() {

    union_type_t a = {.val = 10};

    printf("%d\n", a.val);
    printf("%lld\n", a.err);    // same memory interpret as long long

    color_t x = {.rgba = 2555555};  // hex = 0x00270F13, depending on endianess value will be interpreted.. r = 0x13 or r = 0x00
    printf("%d\n", x.components.r);
    printf("%d\n", x.components.g);
    printf("%d\n", x.components.b);
    printf("%d\n", x.components.a);

    return 0;
}