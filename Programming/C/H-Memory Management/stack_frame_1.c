#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

void stack_frame() {
	static void* prev_sp = NULL;
	void* current_sp = __builtin_frame_address(0);
	long diff = (char*) prev_sp - (char*) current_sp;

	if (prev_sp) {
		printf("Diff = %ld\n", diff);
	}

	if (current_sp) {
		prev_sp = current_sp;
	}
}

void A() {
	int x = 10;
	printf("A() = %p\n", __builtin_frame_address(0));
	stack_frame();
}

void B() {
	int y = 20;
	printf("B() = %p\n", __builtin_frame_address(0));
	stack_frame();
}

void C() {
	int z = 30;
	printf("C() = %p\n", __builtin_frame_address(0));
	stack_frame();
}

void D() {
	A();
	printf("D() = %p\n", __builtin_frame_address(0));
	stack_frame();
}

int main() {

	// A();
	// B();
	// C();

	// stack grows downward, so diff value is negative.
    // D() -> A()
    // higher addr -> lower addr
	D();

	return 0;
}