#include<iostream>
#include<stdlib.h>

using namespace std;

int32_t main() {

    // allocate 1-int sized memory
    int *p = new int;
    delete p;

    // allocate 3-int sized memory
    int *arr = new int[3];
    delete[] arr;

    // C-style allocation
    int *brr = (int*) malloc(3 * sizeof(int));  // in cpp, must typecast
    free(brr);

    // void* type
    double pi = 3.1416;
    void *ptr = &pi;    // typeless-pointer, can point anything

    int *pii = reinterpret_cast<int*>(ptr); // interpret the memory(double*) as *int

    cout << *pii << endl;

    int PII;
    memcpy(&PII, &pi, sizeof(int));

    cout << PII << endl;

    PII = static_cast<int>(pi);
    cout << PII << endl;

    return 0;
}