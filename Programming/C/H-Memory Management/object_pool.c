#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>
#include<stdint.h>
#include<stddef.h>

// Object Pool - Reuse allocated objects, instead of creating a new one..

typedef struct Pair {
	int first, second;
} pair_t;

typedef struct {
	bool allocated;
	pair_t obj;
} pool_elem_t;

const int POOL_SIZE = 10;
pool_elem_t object_pool[POOL_SIZE] = {{0}};

pair_t* Get() {
	for (int i = 0; i < POOL_SIZE; i++) {
		if (object_pool[i].allocated) continue;
		object_pool[i].allocated = true;
		return &object_pool[i].obj;	// return the address
	}
	return NULL; // no free object
}

bool Put(pair_t* obj) {
	if (obj == NULL) return false;
	uintptr_t pool_addr = (uintptr_t) object_pool;	// address of the pool
	uintptr_t obj_addr = (uintptr_t) obj - (uintptr_t) offsetof(pool_elem_t, obj);	// object address inside pool, offset is used to find the initial memory address of where the obj is located in memory..
	size_t size = sizeof(pool_elem_t);	// size of the pool element struct with padding.
	int idx = (obj_addr - pool_addr) / size;	// index of the obj in pool
	if (idx >= POOL_SIZE || idx < 0) return false;
    if ((uintptr_t)(idx * size) + pool_addr != obj_addr) return false; // must be a multiplier..
	object_pool[idx].allocated = false;	// pool claims that memory
	return true;
}

int main() {

	pair_t* obj = Get();
	obj = Get();
	Put(obj);
	obj = Get();
	Put(obj);

	return 0;
}