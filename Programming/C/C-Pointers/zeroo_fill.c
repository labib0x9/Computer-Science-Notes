#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

// For fun, i just changed the zero value, LOL :)
const int INT_ZERO_VAL = -1;
const float FLOAT_ZERO_VAL = -1.00;
const bool BOOL_ZERO_VAL = 1;

typedef enum {
	INT,
	FLOAT,
	BOOL,
} object_kind_t;

typedef struct INT_T {
	int val;
} obj_int_t;

typedef struct FLOAT_T {
	float val;
} obj_float_t;

typedef struct BOOL_T {
	bool val;
} obj_bool_t;

void zero_fill(void* ptr, object_kind_t kind) {
	switch (kind) {
	case INT: {
		obj_int_t* obj = (obj_int_t*) ptr;
		obj->val = INT_ZERO_VAL;
		break;
	}
	case FLOAT: {
		obj_float_t* obj = (obj_float_t*) ptr;
		obj->val = FLOAT_ZERO_VAL;
		break;
	}
	case BOOL: {
		obj_bool_t* obj = (obj_bool_t*) ptr;
		obj->val = BOOL_ZERO_VAL;
		break;
	}
	default:
		return;
	}
}

int main() {

	obj_int_t a;
	zero_fill(&a, INT);
	printf("%d\n", a.val);

	obj_float_t b;
	zero_fill(&b, FLOAT);
	printf("%f\n", b.val);

	obj_bool_t c;
	zero_fill(&c, BOOL);
	printf("%d\n", c.val);

	return 0;
}