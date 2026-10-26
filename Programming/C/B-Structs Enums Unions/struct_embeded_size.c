#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>
#include<stdint.h>
#include<stddef.h>

typedef struct Pair {
	int first, second;
} pair_t;

typedef struct {
	bool allocated;
	pair_t obj;
} pool_elem_t;

const int POOL_SIZE = 10;
pool_elem_t object_pool[POOL_SIZE] = {0};

int main() {

	// size of struct
	printf("pair_t = %zu\n", sizeof(pair_t));
	printf("pool_elem_t = %zu\n\n", sizeof(pool_elem_t));

	// address and size
	uintptr_t prev = (uintptr_t) &object_pool[0];
	for (int i = 0; i < POOL_SIZE; i++) {
		uintptr_t cur = (uintptr_t) &object_pool[i];
		uintptr_t obj_addr = (uintptr_t) &object_pool[i].obj;
		printf("%lu %lu   ::::    %lu %lu\n", cur, cur - prev, obj_addr, obj_addr - cur);
		prev = cur;
	}

	// find the starting offset of field obj
	printf("\n\n");
	printf("%zu\n", offsetof(pool_elem_t, obj));

	size_t offset = offsetof(pool_elem_t, obj);
	uintptr_t obj_addr = (uintptr_t) &object_pool[1].obj;
	uintptr_t field_addr = obj_addr - offset;
	printf("%zu\n", field_addr);

	return 0;
}