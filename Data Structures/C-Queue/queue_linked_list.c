#include<stdio.h>
#include<stdlib.h>
#include<stdbool.h>

typedef struct node {
   int val;
   struct node* next;
} node;

typedef struct {
   node* tail;
   node* head;
} queue;

queue* init_queue() {
   queue* q = (queue*) malloc(sizeof(queue));   
   q->head = q->tail = NULL;
   return q;
}

bool push(queue* q, int val) {
   // new node
   node* newNode = (node*) malloc(sizeof(node));
   if (newNode == NULL) {  // malloc fails
      return false;
   }
   newNode->val = val;
   newNode->next = NULL;

   // first node..
   if(q->head == NULL) {
      q->head = newNode;
      q->tail = newNode;
      return true;
   }

   // chain with prev node ..
   q->tail->next = newNode;
   q->tail = newNode;
   return true;
}

int pop(queue* q) {
   if (q->head == NULL) {  // empty queue..
      return -1;           // -1 ? might be a problem :)
   }

   // get result & free the head and pass to next node
   node* next = q->head->next;
   int val = q->head->val;
   free(q->head);
   q->head = next;
   
   // if head is nil, then tail must be nil, no node 
   if(q->head == NULL) {
      q->tail = NULL;
   }
   return val;
}

int peek(queue* q) {
   if (q->head == NULL) {
      return -1;
   }
   return q->head->val;
}

void free_queue(queue* q) {
  node* cur = q->head;
  while(cur != NULL) {
   node* next = cur->next;
   free(cur);
   cur = next;
  }
  // for reusing, don't free
  // q->head = q->tail = NULL;
  free(q);
}

int main() {

   queue* q = init_queue();

   push(q, 10);
   push(q, 20);
   push(q, 111);

   printf("%d\n", pop(q));
   printf("%d\n", pop(q));
   printf("%d\n", pop(q));
   printf("%d\n", pop(q));

   push(q, 100);
   push(q, 210);

   printf("%d\n", pop(q));
   printf("%d\n", pop(q));

   free_queue(q);

   return 0;
}