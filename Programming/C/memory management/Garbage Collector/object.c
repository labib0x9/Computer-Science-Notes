#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

// object kind, what kind of data type
typedef enum {
	INT,
	FLOAT,
} object_kind_t;

// object value, original value of the object
typedef union {
	int v_int;
	float v_float;
} object_data_t;

// object, it stores the data and metadata of itself
typedef struct Object {
	object_kind_t Kind;	// what kind of object
	object_data_t Data;	// Actual value
} object_t;

object_t* NewObjectInt(int val) {
	object_t* obj = (object_t*) malloc(sizeof(object_t));
	if (obj == NULL) {
		return NULL;
	}
	obj->Kind = INT;
	obj->Data.v_int = val;
	return obj;
}

object_t* NewObjectFloat(float val) {
	object_t* obj = (object_t*) malloc(sizeof(object_t));
	if (obj == NULL) {
		return NULL;
	}
	obj->Kind = FLOAT;
	obj->Data.v_float = val;
	return obj;
}

void PrintObject(object_t* obj) {
	if (obj == NULL) {
		return;
	}
	switch (obj->Kind) {
		case INT:
			printf("%d\n", obj->Data.v_int);
			break;
		case FLOAT:
			printf("%0.3f\n", obj->Data.v_float);
			break;
		default:
			printf("unknown data types\n");
			break;
	}
}

int main() {

	object_t* obj1 = NewObjectInt(10);
	object_t* obj2 = NewObjectFloat(20.012);

	PrintObject(obj1);
	PrintObject(obj2);

    

	if (obj1) free(obj1);
	if (obj2) free(obj2);

	return 0;
}