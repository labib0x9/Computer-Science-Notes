// Array Reverse.
/*
    Given an array, reverse the array.
    So array has a length and capacity & index is zero-based.
*/

#include<stdio.h>
#include<string.h>

void print_array(int arr[], int len) {
    for (int i = 0; i < len; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");
}

int main() {

    int cap = 10, len = 10;
    int arr[cap];
    for (int i = 0; i < len; i++) arr[i] = i + 1;

    print_array(arr, len);
    
    // Reverse
    // swap is using bitmask xor.
    for (int i = 0; i < len / 2; i++) {
        arr[i] ^= arr[len - i - 1];
        arr[len - i - 1] ^= arr[i];
        arr[i] ^= arr[len - i - 1];
    }

    print_array(arr, len);


    // using two pointer.
    int lo = 0, hi = len - 1;
    while(lo <= hi) {
       int temp = arr[lo];
       arr[lo] = arr[hi];
       arr[hi] = temp;
       lo++; hi--;
    }

    print_array(arr, len);

    return 0;
}