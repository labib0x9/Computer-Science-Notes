#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

void test() {
	static int a = 0;	// stored at data segment, once created it persists on function calls...
	a += 10;
	printf("a = %d\n", a);
}

int main() {

	test();	// 10
	test();	// 20
	test();	// 30

	return 0;
}