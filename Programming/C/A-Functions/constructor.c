#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>
#include<stdint.h>
#include<stddef.h>

// runs before main().
__attribute__((constructor)) void init() {
	printf("Hello From init()\n");
}

// alias
typedef uint8_t int8;
typedef uint8_t byte;
typedef uint16_t int16;
typedef uint32_t int32;
typedef uint64_t int64;

// public and private, linux based
#define public __attribute__((visibility("default")))
#define private __attribute__((visibility("hidden")))

private void hello_world() {
	printf("Hello World\n");
}

public void hello_world_public() {
	printf("Hello world public\n");
}

int main() {

	printf("Hello From main()\n");

	return 0;
}