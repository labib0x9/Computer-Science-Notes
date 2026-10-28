#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<pthread.h>
#include<string.h>
#include<stdbool.h>

typedef struct {
   int* q;
   int cap;
   int len;
   int readAt, insertAt;

} queue;

bool init_queue(queue* q) {
   q->cap = 1;
   q->len = 0;
   q->readAt = q->insertAt = 0;
   int *qq = (int*) malloc(sizeof(int) * q->cap);
   if (qq == NULL) {
      return false;
   }
   q->q = qq;
   return true;
}

bool isFull(queue* q) {
   return q->len == q->cap;
}

bool isEmpty(queue* q) {
   return q->len == 0;
}

bool resize(queue* q) {
   if (!isFull(q)) {
      return false;
   }

   // resize by double of capacity
   int newCap = 2 * q->cap;
   int* newQ = (int*) malloc(sizeof(int) * newCap);   // how can i use realloc() ??
   if (newQ == NULL) {
      return false;
   }

   // move elements
   // well, here % is a costly operation, can i use memmove() ??
   for (int i = 0; i < q->len; i++) {
      newQ[i] = q->q[(q->readAt + i) % q->cap];
   }

   // fix the logic...
   // memmove(newQ, q->q + q->readAt, sizeof(int) * (q->len - q->readAt));
   // memmove(newQ + (q->len - q->readAt), q->q, sizeof(int) * (q->readAt + 1));

   // update
   free(q->q);
   q->q = newQ;
   q->cap = newCap;
   q->insertAt = q->len;
   q->readAt = 0;

   return true;
}

bool push(queue* q, int val) {
   if (isFull(q)) {
      bool ok = resize(q);
      if (!ok) return ok;
   }
   q->q[q->insertAt] = val;
   q->insertAt = (q->insertAt + 1) % q->cap; // how can i get rid of % operator ??
   // q->insertAt++;
   // if (q->insertAt == q->cap) q->insertAt = 0;
   q->len++;
   return true;
}

bool pop(queue* q, int *val) {
   if(isEmpty(q)) {
      return false;
   }
   *val = q->q[q->readAt];
   q->readAt = (q->readAt + 1) % q->cap; // how can i get rid of % operator ??
   // q->readAt++;
   // if (q->readAt == q->cap) q->readAt = 0;
   q->len--;
   return true;
}

void free_queue(queue* q) {
    free(q->q);
    free(q);
}

int main() {

   return 0;
}