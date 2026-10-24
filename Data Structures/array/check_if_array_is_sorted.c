// Merge Two Sorted Array.
/*
    Given two array, merge the two array.
    So array has a length & index is zero-based.
*/

#include<stdio.h>
#include<string.h>
#include<stdbool.h>

void print_array(int arr[], int len) {
    for (int i = 0; i < len; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");
}

// check if the array is sorted (ascending order)
bool is_sorted(int arr[], int len) {
    for (int i = 0; i + 1 < len; i++) {
        if (arr[i] > arr[i + 1]) { return false; }
    }
    return true;
}

// descending order
bool is_dsorted(int arr[], int len) {
    for (int i = 0; i + 1 < len; i++) {
        if (arr[i] < arr[i + 1]) { return false; }
    }
    return true;
}

int main() {

    int lenA = 3, lenB = 4, lenC = 4;
    int arr[] = {1, 9, 6};
    int brr[] = {2, 3, 4, 5};
    int crr[] = {9, 5, 5, 2};

    print_array(arr, lenA);
    print_array(brr, lenB);

    // ascending
    bool a = is_sorted(arr, lenA);
    bool b = is_sorted(brr, lenB);

    printf("A = %d, B = %d\n", a, b);

    // descending
    bool c = is_dsorted(brr, lenB);
    bool d = is_dsorted(crr, lenC);

    printf("B = %d, C = %d\n", c, d);

    return 0;
}