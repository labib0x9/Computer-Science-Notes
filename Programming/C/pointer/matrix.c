#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>

/*
    # # # # #
    # # # # #
    # # # # #
*/
typedef struct Matrix {
    int n, m;
    int** mat;
} matrix_t;

matrix_t NewMatrix(int n, int m) {
    matrix_t new = {.n = n, .m = m, .mat = NULL};
    new.mat = (int**) malloc(sizeof(int*) * n);
    if (new.mat == NULL) {
        return new;
    }
    for (int i = 0; i < n; i++) {
        int* elm = (int*) malloc(sizeof(int) * m);
        if (elm == NULL) {
            for (int j = 0; j < i; j++) {
                if (new.mat[j]) {
                    free(new.mat[j]);
                }
            }
            free(new.mat);
            new.mat = NULL;
            return new;
        }
        new.mat[i] = elm;
    }
    return new;
}

matrix_t MatrixAdd(matrix_t* a, matrix_t* b) {
    matrix_t tmp = NewMatrix(a->n, a->m);
    if (tmp.mat == NULL) {
        return;
    }
    // implement
}

void freeMatrix(matrix_t* mat) {
    if (mat == NULL) return;
    if (mat->mat == NULL) return;
    for (int i = 0; i < mat->n; i++) {
        if (mat->mat[i]) free(mat->mat[i]);
    }
}

int main() {

	return 0;
}