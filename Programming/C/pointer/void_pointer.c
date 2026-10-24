#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

// just 10 references..
const int size = 10;
int current_at = 0;
int* HEAP_REF[size];

void test(void* a) {
	int* x = (int*) a;	// must type-cast
	printf("%d\n", *x);
}

void* return_from() {
	int x = 40;
	return (void*) (&x);	// hehehehehe, if you know, you know :)
}

void* return_from_escape_to_heap() {
	int *x = (int*) malloc(sizeof(int));
	*x = 60;
	HEAP_REF[current_at++] = x;
	return (void*) x;
}

int main() {

	int x = 10;
	// test(&x);

	void* tmp = return_from();
	test(&x);
	// printf("%p\n", *tmp);	// not works

	// not the value we stored. why ?
	// because stack frame is destroyed, and we lost whatever data we stored in that region...
	// so what's the solution ? escape to heap, like golang, right ?
	printf("Lost (stored in stack frame) = %d\n", *(int*) tmp);

	// escape to heap
	tmp = return_from_escape_to_heap();
	printf("Got (escape to heap) = %d\n", *(int*) tmp);

	// free the heap memory
	for (int i = 0; i < current_at; i++) {
		if (HEAP_REF[i]) {
			free(HEAP_REF[i]);
		}
	}

	return 0;
}