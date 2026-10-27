#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>
#include<stdint.h>
#include<stddef.h>

// Object Pool - Reuse allocated objects, instead of creating a new one..
// But the object should be same size.
// Free List - each pool element is node and connect these nodes using linked list..
// get the head of the linked list.
// back at the head of the linked list.

typedef struct Pair {
	int first, second;
} pair_t;

typedef struct PoolElem{
	pair_t obj;
	struct PoolElem* next;
} pool_elem_t;

const int POOL_SIZE = 10;
pool_elem_t object_pool[POOL_SIZE] = {{0}};
pool_elem_t* freeList = NULL;

void init_pool() {
	for (int i = 0; i + 1 < POOL_SIZE; i++) {
		object_pool[i].next = &object_pool[i + 1];
	}
	object_pool[POOL_SIZE - 1].next = NULL;
	freeList = &object_pool[0];
}

pair_t* Get() {
	if (freeList == NULL) {	// no free object
		return NULL;
	}
	pair_t* obj = &(freeList->obj);	// extract object
	freeList = freeList->next;	// move to the next
	return obj;
}

bool Put(pair_t* obj) {
	if (obj == NULL) return false;
	uintptr_t pool_addr = (uintptr_t) object_pool;	// address of the pool
	uintptr_t obj_addr = (uintptr_t) obj - (uintptr_t) offsetof(pool_elem_t, obj);	// object address inside pool, offset is used to find the initial memory address of where the obj is located in memory..
	size_t size = sizeof(pool_elem_t);	// size of the pool element struct with padding.
	int idx = (obj_addr - pool_addr) / size;	// index of the obj in pool
	if (idx >= POOL_SIZE || idx < 0) return false;
    if ((uintptr_t)(idx * size) + pool_addr != obj_addr) return false; // must be a multiplier..
	
	// make head
	pool_elem_t* which_to_put = &object_pool[idx];
	which_to_put->next = freeList;	// make chain
	freeList = which_to_put;	// make head
	return true;
}

int main() {

	init_pool();

	pair_t* obj = Get();
	obj = Get();
	Put(obj);
	obj = Get();
	Put(obj);

	return 0;
}