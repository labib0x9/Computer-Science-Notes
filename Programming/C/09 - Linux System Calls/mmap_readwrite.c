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
    int fd = open("test.txt", O_RDWR);
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

    // mmap the area
    char *mmap_addr;
    mmap_addr = mmap(
        NULL,
        status.st_size,
        PROT_READ | PROT_WRITE, 
        MAP_SHARED, 
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

    // Prev content
    write(STDOUT_FILENO, mmap_addr, status.st_size);
    printf("\n");

    // Modify the mapped area
    char *threeByte = "HI ";
    memcpy(mmap_addr, threeByte, strlen(threeByte));

    threeByte = "VAL";
    memcpy(mmap_addr + 3, threeByte, strlen(threeByte));

    // printf("%lu\n", sizeof(threeByte)); // 8-bit, because it is a pointer

    // Updated content
    write(STDOUT_FILENO, mmap_addr, status.st_size);
    printf("\n");

    // flash -> file
    if (msync(mmap_addr, status.st_size, MS_SYNC) == -1) {
        perror("sync");
        munmap(mmap_addr, status.st_size);
        return 0;
    }

    // release mapped area
    if (munmap(mmap_addr, status.st_size) == -1) {
        perror("unmap failed");
        return 0;
    }

    return 0;
}