// Array Deletion.
/*
    Given an array and an index i, delete index i.
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

void delete(int arr[], int idx, int* len) {
    // if (*len >= capacity) { return; }
    memmove(&arr[idx], &arr[idx + 1], sizeof(int) * (*len - idx - 1));
    (*len)--;
}

int main() {

    int cap = 10, len = 10;
    int arr[cap];
    for (int i = 0; i < len; i++) arr[i] = i + 1;

    int idx = 0;
    delete(arr, idx, &len);
    print_array(arr, len);

    idx = 3;
    delete(arr, idx, &len);
    print_array(arr, len);

    idx = len - 1;
    delete(arr, idx, &len);
    print_array(arr, len);

    // using loop
    idx = 3;
    for (int i = idx; i < len - 1; i++) {
        arr[i] = arr[i + 1];
    }
    len--;
    print_array(arr, len);

    return 0;
}