#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<pthread.h>
#include<string.h>

/*
   we have a queue of n fixed element, if the queue is full then the task will wait
   for the slot to free. 
*/

typedef struct Task {
   int a, b;
} Task;

const int QueueCount = 10;
Task queue[QueueCount];
int insertAt = 0, readAt = 0;
int size = 0, finish = 0;

const int MAX_THREAD_NUM = 4;
pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;
pthread_cond_t empty = PTHREAD_COND_INITIALIZER;
pthread_cond_t full = PTHREAD_COND_INITIALIZER;

// Mutex locking to avoid race condition...
void push(Task task) {
   pthread_mutex_lock(&lock); // lock mutex
   while (size == QueueCount) {
      // We will wait here until a task is done, blocking the push ooperation...
      pthread_cond_wait(&full, &lock);
   }
   queue[insertAt] = task;
   insertAt = (insertAt + 1) % QueueCount;
   size++;
   pthread_cond_signal(&empty); // wakeup empty
   pthread_mutex_unlock(&lock);  // unlock mutex
}

// Mutex locking to avoid race condition...
int pop(Task* t) {
   pthread_mutex_lock(&lock); // lock mutex
   while (size == 0 && finish == 0) {
      // We will wait here until a task is inserted, blocking the pop operation...
     pthread_cond_wait(&empty, &lock);
   }
   if(size == 0 && finish) {
      pthread_mutex_unlock(&lock);  // unlock mutex
      return 0;
   }
   *t = queue[readAt];
   readAt = (readAt + 1) % QueueCount;
   size--;
   pthread_cond_signal(&full);   // wakeup full
   pthread_mutex_unlock(&lock);  // unlock mutex
   return 1;
}

void executeTask(Task task) {
   printf("ThreadId = %lu  ->  %d + %d = %d\n", (unsigned long) pthread_self(), task.a, task.b, task.a + task.b);
}

void* startThread(void* args) {
   Task t;
   while(1) {
      int ok = pop(&t);
      if(!ok) break;
      executeTask(t);
   }
   return NULL;
}

int main() {

   pthread_t thread[MAX_THREAD_NUM];

   for(int i = 0; i < MAX_THREAD_NUM; i++) {
      if (pthread_create(&thread[i], NULL, startThread, NULL) != 0) {
         perror("Creating thread");
      }
   }

   for(int i = 1; i <= 20; i++) {
      Task t = {.a = i, .b = i * i};
      push(t);
   }

   pthread_mutex_lock(&lock);
   finish = 1;
   pthread_cond_broadcast(&empty); // wakeup all sleeping pop threads...
   pthread_mutex_unlock(&lock);

   for(int i = 0; i < MAX_THREAD_NUM; i++) {
      if(pthread_join(thread[i], NULL) != 0) {
         perror("Synchorizing thread");
      }
   }

   return 0;
}