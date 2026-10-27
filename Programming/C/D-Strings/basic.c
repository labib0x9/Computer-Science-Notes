#include<stdio.h>
#include<string.h>
#include<stdlib.h>

int main() {

    // string array -> stores at stack -> modifiable
    char s[] = "Hello";
    s[0] = 'X';
    printf("%s\n", s);
    printf("%lu\n", sizeof(s)); // 6 = strlen + '\0'
    printf("%lu\n", strlen(s)); // 5

    // string literal -> stores at read-only section -> non modifiable
    char *ss = "vauvau";
    // ss[0] = 'V'; // bus error -> SIGBUS/SIGSEGV
    printf("%s\n", ss); // no change
    printf("%lu\n", sizeof(ss));  // 8 = sizeof(pointer)
    printf("%lu\n", sizeof(*ss)); // 1 = s[0] -> one char
    printf("%lu\n", strlen(ss));  // 5

    ss = "but i can !"; // can assign another string..
    printf("%s\n", ss);

    ss = s;
    printf("%s\n", ss);


    /* */
    // c string always ends with '\0' -> null terminator
    char sss[4] = "abc";
    sss[3] = '\0';
    printf("%s\n", sss);

    // strncpy() doesn't copy null-terminator. you must do it your own.
    strncpy(sss, "xyz", sizeof(sss));
    sss[3] = '\0';
    printf("%s\n", sss);

    /* find index of delimeter char */
    char s8[8] = "abcd\0efg";
    printf("%s\n", s8);

    // IMPORTANT TO NOTEEE......., give the index from where the iteration begins.
    size_t idx = strcspn(s8, "\0"); // 4
    printf("%zu\n", idx);

    size_t idx = strcspn(s8 + 5, "f"); // 2
    printf("%zu\n", idx);

    /* */
    char s3[3], s6[6];

    // strcpy(s6, s3);

    // strncpy(s3, s6, 3);


    // strcat(s3, s6);
    // strncat(s3, s6, 3);

    // strcmp(s3, s6);
    // strncmp(s3, s6, 3);

    // compare case insensitive
    if(strcasecmp("Tcp", "tcp") == 0) printf("equal\n");    // 0 = equal

    // find first char in sss, that matches pattern. 
    // strpbrk(sss, pattern)
    char* sss = "aggg:ggg";
    char* found = strpbrk(sss, ":");
    printf("%s\n", found + 1);

    // strspn()
    char* sss = "aggg:ggg";
    size_t len = strspn(sss, "ag");
    printf("%zu\n", len);

    // strcspn()
    char* sss2 = "aggg:ggg";
    size_t idx = strcspn(sss2, ":");
    printf("%zu\n", idx);

    // strchr(ss, c)
    // Find first occurances of c in ss
    char* ss = "aggg:ggg";
    char* sep = strchr(ss, ':');
    printf("%d\n", sep - ss);
    printf("%s\n", sep + 1);
    
    // strstr()
    char* sss = "aggg:ggg";
    char* found = strstr(sss, ":");
    printf("%s\n", found + 1);

    // strdup()
    char* dup = strdup("hello");
    printf("%s\n", dup);
    free(dup);

    // // In memory section..
    // memset();
    // memmem();
    // memcpy();
    // memcmp();
    // memchr();


    /// strcpy vs memcpy ///
    // strcpy -> copies until found NULL terminator('\0)
    // memcpy -> you tell how many should be copied ..

    /* input from stdin */
    char buffer[1024];
    fgets(buffer, sizeof(buffer) - 1, stdin);
    buffer[strcspn(buffer, "\n")] = '\0';

    char *newStr;
    newStr = malloc(strlen(buffer) + 1);    // reserve one room for '\0', must must free() later
    strcpy(newStr, buffer);

    printf("%s\n", newStr);

    /* seperate through delimeter */
    // strtok() modifies the original string ...
    char msg[] = "hello world";
    char *token = strtok(msg, " "); // token, msg is same
    printf("%s %s\n", token, msg);

    token = strtok(NULL, " ");
    printf("%s\n", token);

    token = strtok(NULL, " ");
    printf("%s\n", token);

    if (token == NULL) {
        printf("No more token\n");
    }

    // toknize from newStr
    token = strtok(newStr, " ");    // newStr from 'input from stdin' section
    printf("%s\n", token);

    free(newStr);

    /* */
    // sprintf();
    // snprintf();
    char buf[64];
    char *path = "/index.html";
    int n = snprintf(buf, sizeof(buf), "GET %s HTTP/1.1\r\n", path);

    // Conversion
    int n = atoi("10"); // atoi("1232bb");
    float f = atof("10.5");
    long long nn = atoll("111000000000");
    long n2 = strtol("123", NULL, 10);
    char* end;
    long nn2 = strtol("12434n", &end, 10);
    if (end == 'n') {
        printf("No more characters after number\n");
    }

    return 0;
}