const int FIXED_QUEUE_ELEM = 10;
int queue[FIXED_QUEUE_ELEM];
int insertAt = 0, readAt = 0;
int count = 0;

char isEmpty() {
   return count == 0;
}

char isFull() {
   return count == FIXED_QUEUE_ELEM;
}

char push(int val) {
  if(isFull()) {
   printf("QUEUE FULL\n");
   return 0;
  }
  queue[insertAt] = val;
  insertAt = (insertAt + 1) % FIXED_QUEUE_ELEM;
  count++;
  return 1;
}

int pop() {
   if (isEmpty()) {
      printf("EMPTY -> QUEUE\n");
      return -1;
   }
   int x = queue[readAt];
   readAt = (readAt + 1) % FIXED_QUEUE_ELEM;
   count--;
   return x;
}

void print() {
   for(int i = 0; i < count; i++) {
      printf("%d ", queue[(i + readAt) % FIXED_QUEUE_ELEM]);
   }
   printf("\n");
}