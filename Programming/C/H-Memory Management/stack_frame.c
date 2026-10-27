#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

// stack frames
// level = 0 -> current function
// level = 1 -> caller function
// level = 2 -> caller of caller

void call_me() {
	int x = 10;
	printf("Call Me: (top) = %p\n", __builtin_frame_address(0));	// call_me() stack frame address -> lower than the main()
	printf("Call Me: (top) = %lu\n", (uintptr_t) __builtin_frame_address(0));
	printf("Call Me: (ret) = %p\n", __builtin_return_address(0));	// .text section
	printf("Call Me: (ret) = %lu\n", (uintptr_t) __builtin_return_address(0));
}

int main() {

	printf("Main: (top) = %p\n", __builtin_frame_address(0));	// main() stack frame address
	printf("Main: (top) = %lu\n", (uintptr_t) __builtin_frame_address(0));
	call_me();

	return 0;
}