#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

void print(int x) {
	printf("%p = %d\n", __builtin_frame_address(0), x);
	print(x + 1);
}

int main() {

	// overflow - 01
	const int size = 1024 * 1024 * 10;	// 10mb ?
	char buf[size];
	buf[0] = 'A';
	buf[1] = 'C';
	buf[size - 1] = '\0';

	// Now we get segmentation fault ? why ?
	// stack overflow ....

	// overflow - 02
	print(0);

	return 0;
}