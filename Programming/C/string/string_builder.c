#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>

typedef struct String {
    int len;
    char* data;
} string_t;

string_t NewString(char* str) {
    string_t tmp = {.len = 0, .data = NULL};
    if (str == NULL) {
        return tmp;
    }
    int len = strlen(str);
    char* data = (char*) malloc(sizeof(char) * (len + 1));
    if (data == NULL) {
        return tmp;
    }
    tmp.len = len;
    tmp.data = data;
    strncpy(tmp.data, str, sizeof(char) * len);
    tmp.data[len] = '\0';
    return tmp;
}

// a = a + b
bool AppendString(string_t *a, string_t* b) {
    if (a == NULL || b == NULL) return false;
    if (a->data == NULL || b->data == NULL) return false;
    // what happens here ? we store the pointer in another variable, waht if realloc fails ?? ptr will be null, but a->data ?
    char *ptr = (char*) realloc(a->data, sizeof(char) * (a->len + b->len + 1));
    if (ptr == NULL) {
        return false;
    }
    strncpy(ptr + a->len, b->data, sizeof(char) * b->len);
    ptr[a->len + b->len] = '\0';
    a->data = ptr;
    a->len += b->len;
    // free(ptr);
    return true;
}

string_t SubString(string_t* a, int len) {
    string_t tmp = {.data = NULL, .len = 0};
    if (a == NULL) return tmp;
    if (a->data == NULL) return tmp;
    char *ptr = (char*) malloc(sizeof(char) * (len + 1));
    if (ptr == NULL) {
        return tmp;
    }
    tmp = NewString(ptr);
    // free(ptr);
    return tmp;
}

string_t* SplitString(string_t* a, char delim) {
    if (a == NULL) return NULL;
    if (a->data == NULL) return NULL;
    // implement
}

void freeString(string_t* str) {
    if (str->data) free(str->data);
    str->data = NULL;
    str->len = 0;
}

int main() {

	return 0;
}