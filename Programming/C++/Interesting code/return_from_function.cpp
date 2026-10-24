#include<bits/stdc++.h>

using namespace std;

/*
/opt/homebrew/bin/g++-13 -w -std=c++20 test.cpp && (./a.out) -> same address, move is applied

/opt/homebrew/bin/g++-13 -w -std=c++20 -fno-elide-constructors test.cpp && (./a.out) -> different address, move is not applied by default
*/

vector<array<int, 2>> RLE(vector<int>& arr) {
    vector<array<int, 2>> rle;

    printf("%p\n", &rle);
    return rle;
}

int32_t main() {

    vector<int> arr{1, 2, 3, 3, 4};
    auto r = RLE(arr);

    printf("%p\n", &r);

    return 0;
}