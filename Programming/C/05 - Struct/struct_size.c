#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdint.h>
#include<stdbool.h>

typedef struct {
    bool alloc;
    int a;
    int b;
    long long c;
} a_t;

typedef struct {
    long long c;
    int a;
    int b;
    bool alloc;
} b_t;

int32_t main() {

    printf("A = %zu\n", sizeof(a_t));
    printf("B = %zu\n", sizeof(b_t));

    return 0;
}