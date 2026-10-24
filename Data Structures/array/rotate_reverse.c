#include<stdio.h>

void print_array(int arr[], int len) {
    for (int i = 0; i < len; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");
}

void reverse_array(int arr[], int l, int r) {
    for (int i = 0; i < (r - l + 1) / 2; i++) {
        int temp = arr[l + i];
        arr[l + i] = arr[r - i];
        arr[r - i] = temp;
    }
}

int main() {

    int n = 5, k = 3;
    int arr[] = {1, 2, 3, 4, 5};

    print_array(arr, n);

    // Left - Rotation
    reverse_array(arr, 0, k - 1);
    reverse_array(arr, k, n - 1);
    reverse_array(arr, 0, n - 1);

    print_array(arr, n);

    // Right - Rotation
    reverse_array(arr, 0, n - k - 1);
    reverse_array(arr, n - k, n - 1);
    reverse_array(arr, 0, n - 1);

    print_array(arr, n);


    return 0;
}