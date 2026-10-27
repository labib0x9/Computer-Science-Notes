#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<pthread.h>

int count;

void* Calc(void* args) {
   for(int i = 1; i <= 10000000; i++) {
      count++;
   }
   return NULL;
}

int main() {

   count = 0;

   pthread_t alice;
   pthread_create(&alice, NULL, Calc, NULL);

   Calc(NULL);

   pthread_join(alice, NULL);
   printf("%d\n", count);

   return 0;
}