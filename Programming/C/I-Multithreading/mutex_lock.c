#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<pthread.h>

pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;
int count;

void* Calc(void* args) {
   for(int i = 1; i <= 10000000; i++) {
      pthread_mutex_lock(&lock);
      count++;
      pthread_mutex_unlock(&lock);
   }
   return NULL;
}

int main() {

   count = 0;

   pthread_t t;
   pthread_create(&t, NULL, Calc, NULL);

   Calc(NULL);

   pthread_join(t, NULL);
   printf("%d\n", count);

   return 0;
}