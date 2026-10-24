- gcc -save-temps test.c

---

1. source code : test.c
2. prepocessing : gcc -E test.c -o test.i
3. compilation: gcc -S test.c -o test.s
4. assembly: gcc -c test.s -o test.o
5. linking: (h.o)