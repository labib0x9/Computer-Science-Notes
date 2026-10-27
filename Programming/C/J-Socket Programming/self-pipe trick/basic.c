#include<stdio.h>
#include<unistd.h>

/*
    uni directional
    parent-child process
    same process threading

    fork()
    thread wroker()
*/

int main() {    

    // pipe creation
    int pipe_fd[2];
    int ok = pipe(pipe_fd);
    if (ok == -1) {
        perror("Pipe");
        return 0;
    }

    // Data flows, fd[1] -> fd[0]
    printf("Read end, fd[0] = %d\n", pipe_fd[0]);
    printf("Write end, fd[1] = %d\n", pipe_fd[1]);

    // write to pipe
    int n = write(pipe_fd[1], "Hello Labib\n", 12);
    printf("%d bytes are written to fd[1]\n", n);

    // read from pipe
    char buf[16];
    n = read(pipe_fd[0], &buf, sizeof(buf));
    printf("%d bytes are read from fd[0]\n", n);
    printf("Data= %s\n", buf);

    // You must close the fd...
    close(pipe_fd[0]);
    close(pipe_fd[1]);

	return 0;
}