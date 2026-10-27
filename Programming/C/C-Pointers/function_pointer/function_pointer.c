#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<string.h>
#include<stdbool.h>

void* add(void* args) {
	int* A = (int*) args;
	int* res = (int*) malloc(sizeof(int));
	*res = A[0] + A[1];
	return res;
}

int minus(int a, int b) {
	return a - b;
}

int mult(int a, int b) {
	return a * b;
}

typedef void* (*func)(void* args);

int main() {

	int args[2] = {10, 20};
	func f = add;
	int* sum = f(args);

	printf("%d\n", *sum);
	free(sum);

	/**** ****/
	int (*exec)(int, int);
	exec = minus;
	printf("%d\n", exec(20, 10));	// minus

	exec = mult;
	printf("%d\n", exec(10, 5));	// mult

	return 0;
}
