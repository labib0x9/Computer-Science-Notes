#include<stdio.h>
#include<sys/stat.h>
#include<fcntl.h>
#include<unistd.h>
#include<inttypes.h>
#include<errno.h>

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

    // File size
    printf("%jd\n", (intmax_t)status.st_size);

    mode_t mode = status.st_mode; // File type + permissions
    gid_t group = status.st_gid;  // group
    uid_t owner = status.st_uid;  // owner

    S_ISDIR(mode);
    S_ISREG(mode);

    close(fd);
    return 0;
}