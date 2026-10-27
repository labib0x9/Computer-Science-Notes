#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

typedef struct {
	char name[10];
	int age;
} test_t;

void swap(void* a, void* b, int size) {
	if (a == NULL || b == NULL || size == 0) return;
	void* tmp = (void*) malloc(size);
	if (tmp == NULL) return;

	// tmp = a;
	memmove(tmp, a, size);
	// a = b;
	memmove(a, b, size);
	// b = tmp;
	memmove(b, tmp, size);

	free(tmp);
}

int main() {

	/* */
	int x = 10, y = 20;
	swap(&x, &y, sizeof(int));

	printf("x = %d, y = %d\n", x, y);

	/* */
	char* a = "afafadf";
	char* b = "xyz";

	swap(&a, &b, sizeof(char*));
	printf("a = %s, b = %s\n", a, b);

	/* */
	test_t u1 = {.name = "labib", .age = 24};
	test_t u2 = {.name = "usha", .age = 21};

	swap(&u1, &u2, sizeof(test_t));
	printf("u1 = %s %d\n", u1.name, u1.age);
	printf("u2 = %s %d\n", u2.name, u2.age);

	/* buggy */
	printf("int = %zu\n", sizeof(int));
	printf("float = %zu\n", sizeof(float));

	int int_num = 5;
	float float_num = 10.45;
	swap(&int_num, &float_num, sizeof(int));
	printf("int_num = %d, float_num = %f\n", int_num, float_num);	// not what i want, memmove just copies the raw bytes..

	return 0;
}