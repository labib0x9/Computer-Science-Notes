#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<pthread.h>

void* Calc(void* args) {
   pthread_t calcId = pthread_self();
   printf("Calc() = %lu\n", (unsigned long)calcId);
   return NULL;
}

int main() {

   pthread_t mainT = pthread_self();
   printf("main() = %lu\n", (unsigned long)mainT);

   pthread_t t[10];
   for(int i = 0; i < 10; i++) {
      pthread_create(&t[i], NULL, Calc, NULL);
   }

   for(int i = 0; i < 10; i++) {
      pthread_join(t[i], NULL);
   }

   for(int i = 0; i < 10; i++) {
      printf("%d ThreadId = %lu\n", i, (unsigned long) t[i]);
   }

   return 0;
}