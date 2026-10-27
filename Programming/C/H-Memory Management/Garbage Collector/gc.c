// I don't have any idea how to implement, need to study more and more....
/// Incomplete Code


#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>
#include<stdint.h>
#include<stddef.h>

typedef struct GarbageObject {
	void* ptr;
	int ref_count;
	size_t size;
	struct GarbageObject* next;
} garbaje_object_t;

typedef struct GarbageCollector {
	garbaje_object_t* root;
} gc_t;

void* gc_alloc(gc_t* gc, size_t size) {
	void* ptr = malloc(size);
	return ptr;
}

int main() {

	int* arr = (int*) gc_alloc(NULL, sizeof(int) * 10);
	memset(arr, 0, sizeof(int) * 10);

	arr[0] = 10;
	arr[1] = 30;

	for (int i = 0; i < 10; i++) {
		printf("%d ", arr[i]);
	}
	printf("\n");

	free(arr);
	return 0;
}