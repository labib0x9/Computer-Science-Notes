// Merge Array.
/*
    Given two array, merge the two array.
    So array has a length & index is zero-based.
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

    int lenA = 3, lenB = 4;
    int arr[lenA], brr[lenB];

    for (int i = 0; i < lenA; i++) arr[i] = i + 1;
    for (int i = 0; i < lenB; i++)  brr[i] = i + lenA + 1;

    print_array(arr, lenA);
    print_array(brr, lenB);

    // merge array
    // using loop
    int lenC = lenA + lenB;
    int merged_arr[lenC];
    for (int i = 0; i < lenA; i++) merged_arr[i] = arr[i];
    for (int i = 0; i < lenB; i++) merged_arr[i + lenA] = brr[i];

    print_array(merged_arr, lenC);

    // using memmove
    int drr[lenC];
    memmove(&drr, &arr, sizeof(int) * lenA);
    memmove(&drr[lenA], &brr, sizeof(int) * lenB);

    print_array(drr, lenC);

    return 0;
}