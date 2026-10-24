#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<string.h>
#include<stdbool.h>

typedef struct task_node {
	void* (*func)(void* args);
	void* args;
} task_node;

void* print(void* args) {
	int x = *(int*) args;
	printf("%d\n", x);
	return NULL;
}

void* mult(void* args) {
	int x = *(int*) args;
	int *s = (int*) malloc(sizeof(int));
	*s = x * x;
	return (void*) s;
}

int main() {

	int* x = (int*) malloc(sizeof(int));
	*x = 10;

	task_node t = {print, x};
	t.func(t.args);

	t = (task_node){mult, x};
	int* s = t.func(t.args);
	printf("%d\n", *s);

	free(x);
	free(s);

	return 0;
}
