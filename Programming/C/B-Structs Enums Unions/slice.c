#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>
#include<stdint.h>
#include<stddef.h>

// Slice -> Vector -> Dynamic Array
typedef struct Slice {
	int* arr;
	int len, cap;
} slice_t;

slice_t NewSlice(int len, int cap) {
	slice_t tmp = {.arr = NULL, .len = 0, .cap = 0};
	if (cap == 0) {
		cap = 1;
	}
	// assert(len <= cap);
	tmp.arr = (int*) malloc(sizeof(int) * cap);
	if (tmp.arr == NULL) {
		return tmp;
	}
	tmp.len = len;
	tmp.cap = cap;
	return tmp;
}
// slice grows by 2 * cap
bool resize(slice_t* s) {
	s->cap += s->cap;
	// i have some confusion here ....
	int* tmp = (int*) realloc(s->arr, sizeof(int) * s->cap);
	if (tmp == NULL) {
		s->cap /= 2;
		return false;
	}

	printf("Resize\n");

	s->arr = tmp;
	return true;
}

slice_t Append(slice_t s, int x) {
	if (s.len == s.cap) {
		if (!resize(&s)) {
			return s;
		}
	}
	s.arr[s.len] = x;
	s.len++;
	return s;
}

// just copy the elements...
slice_t SubSlice(slice_t s, int at, int len) {
	slice_t tmp = NewSlice(len, len);
	memcpy(tmp.arr, s.arr + at, sizeof(int) * len);
	return tmp;
}

int At(slice_t s, int idx) {
	if (idx >= s.len) {
		return -1;	// segmentation fault
	}
	return s.arr[idx];
}

int len(slice_t arr) {
	return arr.len;
}

void FreeSlice(slice_t s) {
	if (s.arr) free(s.arr);
}

void PrintSlice(slice_t s) {
	for (int i = 0; i < len(s); i++) {
		printf("%d ", At(s, i));
	}
	printf("\n");
}

int main() {

	printf("slice = %zu\n", sizeof(slice_t));
	slice_t arr = NewSlice(0, 0);
	if (arr.arr == NULL) {
		printf("Malloc fails\n");
		return 0;
	}
	arr = Append(arr, 10);
	arr = Append(arr, 20);
	arr = Append(arr, 30);
	arr = Append(arr, 40);
	arr = Append(arr, 50);
	arr = Append(arr, 60);
	arr = Append(arr, 70);
	arr = Append(arr, 80);
	arr = Append(arr, 90);
	arr = Append(arr, 100);

	PrintSlice(arr);

	slice_t sArr = SubSlice(arr, 4, 3);	// O(n)
	PrintSlice(sArr);

	printf("arr = %p\n", arr.arr);
	printf("sArr = %p\n", sArr.arr);

	FreeSlice(arr);
	FreeSlice(sArr);

	return 0;
}