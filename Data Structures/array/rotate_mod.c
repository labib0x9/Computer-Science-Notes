#include<stdio.h>

void print_array(int arr[], int len) {
    for (int i = 0; i < len; i++) {
        printf("%d ", arr[i]);
    }
    printf("\n");
}

int main() {

    int n = 5;
    int arr[] = {1, 2, 3, 4, 5};

    print_array(arr, n);

    // Left - Roatation, k = 1
    int temp = arr[0];
    for (int i = 0; i + 1 < n; i++) {
        arr[i] = arr[i + 1];
    }
    arr[n - 1] = temp;
    print_array(arr, n);


    // Right - Rotation, k = 1
    temp = arr[n - 1];
    for (int i = n - 1; i - 1 >= 0; i--) {
        arr[i] = arr[i - 1];
    }
    arr[0] = temp;
    print_array(arr, n);


    // Left - Rotation
    int k = 1;
    int ans[n] = {};
    for (int i = 0; i < n; i++) {
        ans[i] = arr[(i + k + n) % n];
    }


    print_array(ans, n);

    for (int i = 0; i < n; i++) arr[i] = ans[i];
    ans[n] = {};

    // Right - Rotation
    for (int i = 0; i < n; i++) {
        ans[i] = arr[(i - k + n) % n];
    }
    print_array(ans, n);


    return 0;
}