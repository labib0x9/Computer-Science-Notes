#include<stdio.h>
#include<stdlib.h>
#include<string.h>

int32_t main() {

    // malloc

    // allocate memory for 1 int on heap..
    int *ptr = malloc(sizeof(int)); // garbage value
    if (ptr == NULL) {
        // malloc fails
    }
    // casting is required in c++
    // int *ptr1 = (int*) malloc(sizeof(int));
    // if (ptr1 == NULL) {
    //     // malloc fails
    // }

    memset(ptr, 0, sizeof(int)); // set to 0x00
    // memset(ptr, 1, sizeof(int)); // garbage value
    // sets bytes to 0x01 so ptr = 0x01010101 = some value (garbage value)

    printf("Address: %p, Value: %d\n", ptr, *ptr);

    if(ptr) {free(ptr);}
    // if(ptr1) {free(ptr1);}


    // calloc
    int *cptr = calloc(3, sizeof(int)); // calloc = malloc + memset
    if (cptr == NULL) {
        // calloc fails
    }

    printf("Address: %p, Value: %d\n", cptr, *cptr);

    if (cptr) {
        free(cptr);
    }

    // realloc
    int n = 5;
    int *arr = malloc(n * sizeof(int));
    if (arr == NULL) {
        // malloc fails
    }

    for(int i = 0; i < n; i++) {
        arr[i] = i * i;
    }

    printf("arr[3] = %d\n", arr[3]);
    printf("*(arr+3) = %d\n", *(arr + 3));

    printf("Address of arr (before realloc): %p\n", arr);
    printf("%p\n", arr + 3);

    printf("out of index: %d\n", arr[n]); // out of index
    
    int m = n + n;
    int *brr = realloc(arr, m * sizeof(int));
    if (brr == NULL) {
        // realloc fails, but doesn't free arr
        // free(arr);   // free memory 
    }

    if (brr) {
        // what happens to arr which was allocated before ??
        // well, if there is enough memory, realloc extens
        // else allocate new block, in this case prev allocation freed, so arr become dangling pointer
        // so no need to free arr, realloc does it for us...

        arr = brr;
        // free(brr); you cann't free brr as it is used in later
    }

    printf("arr[n]: %d\n", arr[n]);
    printf("arr[3]: %d\n", arr[3]);

    printf("Address of arr (after realloc): %p\n", arr);

    if(arr) {
        free(arr);
    }

    ///

    int A[10];

    // sets each byte
    memset(A, 0, sizeof(A));    // 0x00
    memset(A, -1, sizeof(A));   // 0xff
    memset(A, 1, sizeof(A));    // 0x01

    A[4] = 10;

    int B[5];
    // copy memory
    memcpy(B, A, sizeof(B));    // no overlap handling

    for(int i = 0; i < 5; i++) {
        printf("%d ", B[i]);
    }
    printf("\n");

    for(int i = 0; i < 10; i++) {
        A[i] = i + 1;
    }

    // copy a[0], a[1], a[2] to a[4], a[5], a[6] -> no overlap
    memcpy(A + 4, A, 3 * sizeof(int));
    for(int i = 0; i < 10; i++) {
        printf("%d ", A[i]);
    }
    printf("\n");

    // s[5] -> s[0]
    char s[] = "helloworld";
    memcpy(s + 5, s, 5 * sizeof(char)); // depends on optimization, else overwritten.

    printf("memcpy: %s\n", s);

    memmove(s, "helloworld", sizeof(s));
    memmove(s + 5, s, 5 * sizeof(char));    // safe overlap

    printf("memcpy: %s\n", s);

    // compare byte by byte
    int cmp = memcmp(A, B, 5);

    // memchr(s, 'd', sizeof(s));  // first occurence of a byte
    // memmem(str, sizeof(str), pattern, sizeof(pattern)); // finds subsequence
    // memccpy(dest, src, byte, sizeof(src)); // copy until byte is found or sizeof(src).

    // aligned_alloc();
    // restrict ;

    return 0;
}