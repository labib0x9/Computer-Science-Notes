#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

typedef struct GStack {
	int* data;
	int len;
	int cap;
} stack;

stack* newStack(int cap) {
	if (cap == 0) cap = 1;
	stack* st = (stack*) malloc(sizeof(stack));
	if (st == NULL) {
		return NULL;
	}
	
	st->data = (int*) malloc(sizeof(int) * cap);
	if (st->data == NULL) {
		free(st);
		return NULL;
	}

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
	int* data = (int*) realloc(st->data, sizeof(int) * st->cap);
	if (data == NULL) {
		st->cap /= 2;
		return false;
	}
	st->data = data;
	return true;
}

bool push(stack* st, int x) {
	if (isFull(st)) { 
		if (!resize(st)) return false;
	}
	st->data[st->len] = x;
	st->len++;
	return true;
}

int pop(stack* st) {
	if (isEmpty(st)) {
		return -1;	// what if we store -1 in our stack ?
	}
	int x = st->data[st->len - 1];
	st->len--;
	return x;
}

void printStack(stack* st) {
	for (int i = 0; i < st->len; i++) {
		printf("%d ", st->data[i]);
	}
	printf("\n");
}

void freeStack(stack* st) {
	if (!st) {
		return;
	}
	if (st->data) {
		free(st->data);
	}
	free(st);
}

int main() {

	stack* st = newStack(0);
	if (st == NULL) return 0;

	push(st, 10);
	push(st, 20);
	push(st, 30);
	printf("%d\n", pop(st));
	printf("%d\n", pop(st));
	push(st, 40);
	push(st, 50);
	printf("%d\n", pop(st));
	push(st, 60);
	printf("%d\n", pop(st));

	printStack(st);

	pop(st);
	pop(st);

	printStack(st);

	printf("%d\n", pop(st));

	freeStack(st);

	return 0;
}