#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

// const int pop_sized = 256;
// void* POP_REF[pop_sized];
// int track_ptr = 0;

typedef struct GStack {
	void** data;	// (void*)*
	// bool* heapTrack;
	int len;
	int cap;
} stack;

stack* newStack(int cap) {
	if (cap == 0) cap = 1;
	stack* st = (stack*) malloc(sizeof(stack));
	if (st == NULL) {
		return NULL;
	}
	
	st->data = (void*) malloc(sizeof(void*) * cap);
	if (st->data == NULL) {
		free(st);
		return NULL;
	}

	// st->heapTrack = (bool*) malloc(sizeof(bool) * cap);
	// if (st->heapTrack == NULL) {
	// 	free(st->data);
	// 	free(st);
	// 	return NULL;
	// }

	st->cap = cap;
	st->len = 0;
	return st;
}

bool isEmpty(stack* st) {
	return st->len == 0;
}

bool isFull(stack* st) {
	return st->len == st->cap;
}

bool resize(stack* st) {
	st->cap += st->cap;	// resize by double
	void** data = (void*) realloc(st->data, sizeof(void*) * st->cap);
	if (data == NULL) {
		st->cap /= 2;
		return false;
	}
	// bool* heapTrack = (bool*) realloc(st->heapTrack, sizeof(bool) * st->cap);
	// if (heapTrack == NULL) {
	// 	st->cap /= 2;
	// 	free(data);	// can be dangerous :)
	// 	return false;
	// }
	st->data = data;
	// st->heapTrack = heapTrack;
	return true;
}

// bool push(stack* st, void* x, bool onHeap) {
bool push(stack* st, void* x) {
	if (isFull(st)) { 
		if (!resize(st)) return false;
	}
	st->data[st->len] = x;
	// st->heapTrack[st->len] = onHeap;
	st->len++;
	return true;
}

// what if we allocate elements in heap ?
// poping would give the address, but how do we know when to free them ?
// We used another two array :|||
void* pop(stack* st) {
	if (isEmpty(st)) {
		return NULL;	// what if we store -1 in our stack ?
	}
	void* x = st->data[st->len - 1];
	// if (st->heapTrack[st->len - 1]) POP_REF[track_ptr++] = x;
	st->len--;
	return x;
}

void printStack(stack* st) {
	for (int i = 0; i < st->len; i++) {
		printf("%d ", *(int*) st->data[i]);
	}
	printf("\n");
}

void freeStack(stack* st) {
	if (!st) {
		return;
	}

	if (st->data) {
		// for (int i = 0; i < st->len; i++) {
		// 	if (st->heapTrack[i]) free(st->data[i]);
		// }
		free(st->data);
	}
	// if (st->heapTrack) free(st->heapTrack);
	free(st);

	// for (int i = 0; i < track_ptr; i++) {
	// 	free(POP_REF[i]);
	// }
}

int main() {

	stack* st = newStack(0);
	if (st == NULL) return 0;

	int arr[10] = {10, 20, 30, 40, 50, 60, 70, 80, 90, 100};

	push(st, (void*) &arr[0]);

	int* x = (int*) malloc(sizeof(int));
	*x = 30;

	push(st, x);
	// push(st, 30);
	// printf("%d\n", pop(st));
	// printf("%d\n", pop(st));
	// push(st, 40);
	// push(st, 50);
	printf("%d\n", *(int*) pop(st));
	// push(st, 60);
	// printf("%d\n", pop(st));

	printStack(st);

	// pop(st);
	// pop(st);

	// printStack(st);

	// printf("%d\n", pop(st));

	freeStack(st);
	free(x);

	// push diff types
	// *int
	// *float
	// int
	// *char

	return 0;
}