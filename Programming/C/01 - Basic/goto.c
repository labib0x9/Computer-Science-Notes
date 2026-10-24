#include<stdio.h>

int main() {

    int a = 20;
    
    if (a > 10) {
        goto RESET;
    }

    RESET:
        a -= 10;

    printf("a=%d\n", a);

    return 0;
}