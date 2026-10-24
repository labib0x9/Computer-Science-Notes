// Array Insertion.
/*
    Given an array and an index i, insert at index i.
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

void insert(int arr[], int idx, int x, int* len, int capacity) {
    if (*len >= capacity) { return; }
    memmove(&arr[idx + 1], &arr[idx], sizeof(int) * (*len - idx));
    (*len)++;
    arr[idx] = x;
}

int main() {

    int cap = 10, len = 5;
    int arr[cap];
    for (int i = 0; i < len; i++) arr[i] = i + 1;

    int idx = 0, x = 100;
    insert(arr, idx, x, &len, cap);
    print_array(arr, len);

    idx = 3, x = 200;
    insert(arr, idx, x, &len, cap);
    print_array(arr, len);

    idx = len - 1, x = 300;
    insert(arr, idx, x, &len, cap);
    print_array(arr, len);

    // using loop
    idx = 3; x = 400;
    for (int i = len - 1; i >= idx; i--) {
        arr[i + 1] = arr[i];
    }
    arr[idx] = x;
    len++;
    print_array(arr, len);

    return 0;
}