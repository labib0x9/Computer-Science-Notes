#include<stdio.h>
#include<sys/stat.h>
#include<sys/mman.h>
#include<fcntl.h>
#include<unistd.h>
#include<inttypes.h>
#include<errno.h>
#include<string.h>

int main() {

    // File descriptor
    int fd = open("test.txt", O_RDONLY);
    if (fd == -1) {
        // error
        perror("open");
        return 0;
    }
    
    // File status
    struct stat status;
    if (fstat(fd, &status) == -1) {
        // error
        perror("read status");
        close(fd);
        return 0;
    }

    // empty files
    if (status.st_size == 0) {
        perror("mapped area cannot be empty");
        return 0;
    }

    // File size
    printf("%jd\n", (intmax_t)status.st_size);

    // mmap the area
    char *mmap_addr;
    mmap_addr = mmap(
        NULL,
        status.st_size,
        PROT_READ, 
        MAP_PRIVATE, 
        fd, 
        0
    );

    // check if mapped failed
    if (mmap_addr == MAP_FAILED) {
        perror("map failed");
        close(fd);
        return 0;
    }

    // // close the file descriptor..
    close(fd);

    // read from the mapped area
    // remember to add null terminator
    char threeByte[4];
    memcpy(threeByte, mmap_addr, sizeof(threeByte));
    threeByte[3] = '\0';
    printf("%s\n", threeByte);

    memcpy(threeByte, mmap_addr + 3, sizeof(threeByte));
    threeByte[3] = '\0';
    printf("%s\n", threeByte);

    // release mapped area
    if (munmap(mmap_addr, status.st_size) == -1) {
        perror("unmap failed");
        return 0;
    }

    return 0;
}