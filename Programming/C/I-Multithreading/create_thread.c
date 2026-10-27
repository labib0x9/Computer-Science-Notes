#include<stdio.h>
#include<unistd.h>
#include<pthread.h>

// func must be void*(void*)

void* AliceTurn(void* args) {
   for (int i = 0; i <= 10; i++) {
      sleep(1);
      printf("Alice = %d\n", i);
   }
   return NULL;
}

void* BobTurn(void* args) {
   for (int i = 0; i <= 4; i++) {
      sleep(1);
      printf("Bob = %d\n", i);
   }
   return NULL;
}

int main() {

   pthread_t alice;
   pthread_create(&alice, NULL, AliceTurn, NULL);

   pthread_t bob;
   pthread_create(&bob, NULL, BobTurn, NULL);

   // We have to wait, if the main() function executes -> threads are also executed .....
   // Just like go's concurrency :)))))))....
   sleep(10);

   return 0;
}