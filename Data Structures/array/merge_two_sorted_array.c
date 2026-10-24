// Merge Two Sorted Array.
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
    int arr[] = {1, 4, 6};
    int brr[] = {2, 3, 4, 5};

    print_array(arr, lenA);
    print_array(brr, lenB);

    // merge using two pointer.
    // one pointer will point at arr.
    // another will point at brr.

    int ptr_a = 0, ptr_b = 0, ptr_idx = 0;
    int merged_arr[lenA + lenB];

    // copy one array entirely. why not the smallest array?
    // A = {100}
    // B = {1, 2, 3}
    while (ptr_a < lenA && ptr_b < lenB) {
        if (arr[ptr_a] < brr[ptr_b]) merged_arr[ptr_idx++] = arr[ptr_a++];
        else merged_arr[ptr_idx++] = brr[ptr_b++];
    }

    // if arr didn't finish
    while (ptr_a < lenA) merged_arr[ptr_idx++] = arr[ptr_a++];

    // if brr didn't finish
    while (ptr_b < lenB) merged_arr[ptr_idx++] = brr[ptr_b]++;

    print_array(merged_arr, lenA + lenB);

    return 0;
}