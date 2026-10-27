#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>

typedef struct Object object_t; // alias for ordering

// length and data
typedef struct Array {
	int len;
	object_t** data;	// (object*) * 
} array_t;

// combination of 3 objects.
typedef struct Vector3 {
	object_t* x;
	object_t* y;
	object_t* z;
} vector3_t;

// object kind, what kind of data type
typedef enum {
	INT,
	FLOAT,
	STRING,
	ARRAY,
	VECTOR3,
} object_kind_t;

// object value, original value of the object
typedef union {
	int v_int;
	float v_float;
	char* v_string;
	array_t v_array;
	vector3_t v_vector3;
} object_data_t;

// object, it stores the data and metadata of itself
struct Object {
	object_kind_t Kind;	// what kind of object
	object_data_t Data;	// Actual value
};

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

object_t* NewObjectString(char* val) {
	object_t* obj = (object_t*) malloc(sizeof(object_t));
	if (obj == NULL) {
		return NULL;
	}
	obj->Kind = STRING;
	// obj->Data.v_string = val;	// is there any problem ?? we just copied the pointer to obj.
	// 							// or just strcpy() string...
	int n = strlen(val);
	obj->Data.v_string = (char*) malloc(sizeof(char) * (n + 1));	// don't forget for '\0'
	if (obj->Data.v_string == NULL) {
		free(obj);
		return NULL;
	}
	strncpy(obj->Data.v_string, val, sizeof(char) * n);
	obj->Data.v_string[n] = '\0';	// null terminator..
	return obj;
}

object_t* NewObjectArray(int size) {
	object_t* obj = (object_t*) malloc(sizeof(object_t));
	if (obj == NULL) {
		return NULL;
	}
	obj->Kind = ARRAY;
	obj->Data.v_array.len = size;
	obj->Data.v_array.data = (object_t**) malloc(sizeof(object_t*) * size);
	memset(obj->Data.v_array.data, 0, sizeof(object_t*) * size);
	if (obj->Data.v_array.data == NULL) {
		free(obj);
		return NULL;
	}
	return obj;
}

object_t* NewObjectVector3(object_t* x, object_t* y, object_t* z) {
	// if (x == NULL || y == NULL || z == NULL) {
	// 	return NULL;
	// }
	object_t* obj = (object_t*) malloc(sizeof(object_t));
	if (obj == NULL) {
		return NULL;
	}
	obj->Kind = VECTOR3;
	obj->Data.v_vector3.x = x;
	obj->Data.v_vector3.y = y;
	obj->Data.v_vector3.z = z;
	return obj;
}

// obj[index] = val
bool set_array(object_t* obj, int index, object_t* val) {
	if (obj == NULL || val == NULL) return false;
	if (obj->Kind != ARRAY) return false;
	if (index >= obj->Data.v_array.len) return false;
	obj->Data.v_array.data[index] = val;
	return true;
}

// retrive obj[index]
object_t* get_array(object_t* obj, int index) {
	if (obj == NULL) return NULL;
	if (obj->Kind != ARRAY) return NULL;
	if (index >= obj->Data.v_array.len) return NULL;
	return obj->Data.v_array.data[index];
}

// length of object
int len(object_t* obj) {
	if (obj == NULL) return -1;
	switch (obj->Kind) {
		case INT:
			return 1;
		case FLOAT:
			return 1;
		case STRING:
			return strlen(obj->Data.v_string);
		case ARRAY:
			return obj->Data.v_array.len;
		case VECTOR3:
			return 3;
		default:
			return -1;
	}
}

// merge two objects of same kind except INT and FLOAT..
object_t* add(object_t* a, object_t* b) {
	if (a == NULL && b == NULL) return NULL;
	// if (a->Kind != b->Kind) return NULL;

	switch (a->Kind) {
		case INT: {
			switch (b->Kind) {
				case INT:
					return NewObjectInt(a->Data.v_int + b->Data.v_int);
				case FLOAT:
					return NewObjectFloat((float) a->Data.v_int + b->Data.v_float);
				default:
					return NULL;
			}
		}
		case FLOAT: {
			switch (b->Kind) {
				case INT:
					return NewObjectFloat((float) b->Data.v_int + a->Data.v_float);
				case FLOAT:
					return NewObjectFloat(a->Data.v_float + b->Data.v_float);
				default:
					return NULL;
			}
		}
		case STRING: {
			if (b->Kind != STRING) return NULL;
			int n = strlen(a->Data.v_string), m = strlen(b->Data.v_string);
			char* tmpStr = (char*) malloc(sizeof(char) * (n + m + 1));
			strncpy(tmpStr, a->Data.v_string, sizeof(char) * n);
			strncpy(tmpStr + n, b->Data.v_string, sizeof(char) * m);
			tmpStr[n + m] = '\0';
			object_t* tmp = NewObjectString(tmpStr);
			free(tmpStr);
			return tmp;
		}
		case ARRAY: {
			if (b->Kind != ARRAY) return NULL;
			int n = a->Data.v_array.len, m = b->Data.v_array.len;
			object_t* tmp = NewObjectArray(n + m);
			for (int i = 0; i < n; i++) {
				set_array(tmp, i, get_array(a, i));
			}
			for (int i = 0; i < m; i++) {
				set_array(tmp, i + n, get_array(b, i));
			}
			return tmp;
		}
		case VECTOR3: {
			if (b->Kind != VECTOR3) return NULL;
			return NewObjectVector3(
				add(a->Data.v_vector3.x, b->Data.v_vector3.x),
				add(a->Data.v_vector3.y, b->Data.v_vector3.y),
				add(a->Data.v_vector3.z, b->Data.v_vector3.z)
			);
		}
		default:
			return NULL;
	}
}

// Print
void PrintObject(object_t* obj, char delim) {
	if (obj == NULL) {
		return;
	}
	switch (obj->Kind) {
		case INT:
			printf("%d%c", obj->Data.v_int, delim);
			break;
		case FLOAT:
			printf("%0.3f%c", obj->Data.v_float, delim);
			break;
		case ARRAY:
			for (int i = 0; i < obj->Data.v_array.len; i++) {
				PrintObject(get_array(obj, i), ' ');
			}
			printf("%c", delim);
			break;
		case STRING:
			printf("%s%c", obj->Data.v_string, delim);
			break;
		case VECTOR3:
			printf("{\n");
			PrintObject(obj->Data.v_vector3.x, ' ');
			printf("\n");
			PrintObject(obj->Data.v_vector3.y, ' ');
			printf("\n");
			PrintObject(obj->Data.v_vector3.z, ' ');
			printf("\n}");
			printf("%c", delim);
			break;
		default:
			printf("unknown data types\n");
			break;
	}
}

// do we need to assign NULL to freed pointer ????
void freeObject(object_t* obj) {
	if (obj == NULL) return;
	switch (obj->Kind) {
		case INT: {
			// printf("INT - FREE\n");
			free(obj);
			return;
		}
		case FLOAT: {
			// printf("FLOAT - FREE\n");
			free(obj);
			return;
		}
		case STRING: {
			// printf("STRING - FREE\n");
			if (obj->Data.v_string) {
				free(obj->Data.v_string);
			}
			free(obj);
			return;
		}
		case ARRAY: {
			// printf("ARRAY - FREE\n");
			for (int i = 0; i < obj->Data.v_array.len; i++) {
				freeObject(get_array(obj, i));
			}
			if (obj->Data.v_array.data) {
				free (obj->Data.v_array.data);
			}
			free(obj);
			return;
		}
		case VECTOR3: {
			// printf("VECTOR3 - FREE\n");
			freeObject(obj->Data.v_vector3.x);
			freeObject(obj->Data.v_vector3.y);
			freeObject(obj->Data.v_vector3.z);
			free(obj);
			return;
		}
		default:
			free(obj);
			return;
	}
}

int main() {

	/* int, float */

	object_t* obj1 = NewObjectInt(10);
	object_t* obj2 = NewObjectFloat(20.012);

	// PrintObject(obj1, '\n');
	// PrintObject(obj2, '\n');

	/* array */

	int arr_size = 10;
	object_t* arr = NewObjectArray(arr_size);
	for (int i = 0; i < arr_size; i++) {
		set_array(arr, i, NewObjectInt(i));
	}
	// PrintObject(arr, '\n');

	/* string */
	object_t* str = NewObjectString("HELLO WORLD");
	if (str == NULL) {
		printf("STR NULL");
	}
	// PrintObject(str, '\n');

	/* vector3 */
	object_t* vec = NewObjectVector3(arr, str, obj1);
	// PrintObject(vec, '\n');

	object_t* vec1 = NewObjectVector3(vec, obj2, NULL);
	if (vec1 == NULL) { printf("vec1 is NULL\n"); }
	PrintObject(vec1, '\n');

	printf("SUCCESFULLY PRINTTED ALL\n");

	// // vec will free them
	// freeObject(obj1);
	// freeObject(str);
	// freeObject(arr);

	// // vec1 will free them now
	// freeObject(obj2);
	// freeObject(vec);

	freeObject(vec1);

	return 0;
}