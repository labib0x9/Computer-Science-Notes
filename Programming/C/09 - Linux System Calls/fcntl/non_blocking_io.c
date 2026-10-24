#include<stdio.h>
#include<unistd.h>
#include<fcntl.h>
#include<sys/socket.h>

int set_fd_nonblock(int fd) {
    int flags = fcntl(fd, F_GETFL, 0);
    if (flags == -1) {
        return -1;
    }
    return fcntl(fd, F_SETFL, flags | O_NONBLOCK);
}

int main() {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd == -1) {
        perror("socket");
        return 1;
    }
    if (set_fd_nonblock(fd) == -1) {
        perror("set_fd_nonblock");
        close(fd);
        return 1;
    }
    printf("Socket set to non-blocking mode\n");
    close(fd);
    return 0;
}