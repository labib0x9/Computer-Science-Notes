#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<pthread.h>

// Passing argument....

void* AliceTurn(void* args) {
   int* ptr = (int*)malloc(sizeof(int)); // alloc 4bytes
   for (int i = 0; i <= 10; i++) {
      sleep(1);
      printf("Alice = %d\n", i);
      (*ptr) += 2;
   }
   return ptr; // return result pointer
}

int main() {

   pthread_t alice;
   pthread_create(&alice, NULL, AliceTurn, NULL);

   int* result;
   pthread_join(alice, (void*) &result); // wait alice to finish and we should pass the pointer of resulting variable..

   printf("%d\n", *result);

   return 0;
}