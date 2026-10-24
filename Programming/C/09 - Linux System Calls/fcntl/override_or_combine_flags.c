#include<stdio.h>
#include<unistd.h>
#include<fcntl.h>

#include<sys/event.h>
#include<sys/types.h>

#include<pthread.h>

int pipe_fd[2];

int main() {    

    // pipe creation
    int ok = pipe(pipe_fd);
    if (ok == -1) {
        perror("Pipe");
        return 0;
    }

    printf("=========================\n");

    // Combined flag
    int flags = fcntl(pipe_fd[0], F_GETFL, 0);
    printf("READER FLAG= %d\n", flags);
    fcntl(pipe_fd[0], F_SETFL, flags | O_NONBLOCK);
    printf("READER FLAG= %d\n", fcntl(pipe_fd[0], F_GETFL, 0));


    // Flag override, can be dangerous
    flags = fcntl(pipe_fd[1], F_GETFL, 0);
    printf("WRITER FLAG= %d\n", flags);
    fcntl(pipe_fd[1], F_SETFL, O_NONBLOCK);
    printf("WRITER FLAG= %d\n", fcntl(pipe_fd[1], F_GETFL, 0));

    printf("===========================\n");

    // Data flows, fd[1] -> fd[0]
    printf("Read end, fd[0] = %d\n", pipe_fd[0]);
    printf("Write end, fd[1] = %d\n", pipe_fd[1]);

    close(pipe_fd[0]);
    close(pipe_fd[1]);
    
	return 0;
}