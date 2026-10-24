#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"khash.h"

// initailizes a hash table named map, which stores as string-int
// also map is just a name, not a variable
// uses Open addressing + linear probing
KHASH_MAP_INIT_STR(map, int);

// why free key with const char ? 
// if you allocate the key, khash will still expose it as const char *.
// myself own's the memory, but khash forbids mutation.
static inline void free_khash_key(const char *k) {
    free((char *)k);
}

int main() {

    // allocates the hash table named mp; returns a pointer to that table.
    // khash never allocates keys or values
    khash_t(map) *mp = kh_init(map);
    if (mp == NULL) { /* failed */ }

    // Insert a value to hash table
    // Find the key, then store the value

    // ret = 0 -> key exist
    // ret = 1 -> inserted
    // ret = -1 -> allocation failed
    int ret;
    char *key = strdup("labib");    // on heap, must be freed.
    int value = 10;
    khiter_t it = kh_put(map, mp, key, &ret);   // it = the slot index where the key is stored or already existed

    if (ret < 0) { /* allocation failed */ }
    else if (ret == 0) { 
        // free(kh_value(mp, it));  // free the existing value, if dynamically allocated
    }

    kh_value(mp, it) = value;   // put value to key


    // Lookup a key
    khiter_t found = kh_get(map, mp, key);
    if (found == kh_end(mp)) {
        // not found
    } else {
        printf("Value found = %d\n", kh_value(mp, found));
    }

    // Remove a key
    it = kh_get(map, mp, key);
    if (it != kh_end(mp)) {
        // free the memory, if dynamically allocated.
        free((char *)kh_key(mp, it));   // note to cast to (char*)
        // free(kh_value(mp, it));

        // remove the key, but it only releases the slot
        kh_del(map, mp, it);
    }

    // iteration. for debug, delete. free key, value
    for (khiter_t k = kh_begin(mp); k != kh_end(mp); k++) {
        if (kh_exist(mp, k)) {
            printf("%s => %d\n", kh_key(mp, k), kh_value(mp, k));
        }
    }
    
    // destroy / free the hash table
    kh_destroy(map, mp);

    /***************/
    // OTHER
    // char *val = strdup("hello");
    // (void) val;

    /*. NOTE */
    // 1. key -> string literal (char *p = "FD"), do not free
    // 2. key -> heap || strdup, free with cast. free((char*) kh_key(mp, it))

    return 0;
}