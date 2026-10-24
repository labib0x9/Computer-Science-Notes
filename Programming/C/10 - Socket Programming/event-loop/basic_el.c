#include<stdio.h>
#include<unistd.h>
#include<fcntl.h>

#include<sys/event.h>
#include<sys/types.h>

#include<pthread.h>

int pipe_fd[2];

// writes to a pipe
// thread worker
void* worker(void* arg) {
    WORK_TO_DO:
        sleep(10);
    char x = 'x';
    if (write(pipe_fd[1], &x, 1) != 1) goto WORK_TO_DO;

    printf("Worker= %zu, Wrote to pipe_fd[1]\n", (unsigned long) pthread_self());
    return NULL;
}

int main() {    

    // pipe creation
    int ok = pipe(pipe_fd);
    if (ok == -1) {
        perror("Pipe");
        return 0;
    }

    printf("===========================\n");

    // non-blocking
    int flags = fcntl(pipe_fd[0], F_GETFL, 0);
    printf("READER FLAG= %d\n", flags);
    fcntl(pipe_fd[0], F_SETFL, flags | O_NONBLOCK);
    printf("READER FLAG= %d\n", flags | O_NONBLOCK);

    flags = fcntl(pipe_fd[1], F_GETFL, 0);
    printf("WRITER FLAG= %d\n", flags);
    fcntl(pipe_fd[1], F_SETFL, flags | O_NONBLOCK);
    printf("WRITER FLAG= %d\n", flags | O_NONBLOCK);

    printf("===========================\n");

    // Data flows, fd[1] -> fd[0]
    printf("Read end, fd[0] = %d\n", pipe_fd[0]);
    printf("Write end, fd[1] = %d\n", pipe_fd[1]);

    printf("===========================\n");

    // thread
    pthread_t tid;
    pthread_create(&tid, NULL, worker, NULL);

    // kernel event queue is created 
    int e_fd = kqueue();
    if (e_fd == -1) {
        perror("Kqueue");
        return 0;
    }

    printf("KQUEUE_FD= %d\n", e_fd);

    // Register a event
    // will notify if pipe_fd[0] is ready to read.
    struct kevent ev;
    EV_SET(&ev, pipe_fd[0], EVFILT_READ, EV_ADD | EV_CLEAR, 0, 0, NULL);
    if (kevent(e_fd, &ev, 1, NULL, 0, NULL) == -1) {
        perror("Register");
        return 0;
    };

    printf("EL started\n");
    // event loop
    struct kevent events[16];
    int n = kevent(e_fd, NULL, 0, events, 16, NULL);    // hangs.... WHY ? BCZ, kernel blocks as no fd is ready to receieve...

    printf("event, n= %d\n", n);

    // iterate over all the event kernel returns...
    for (int i = 0; i < n; i++) {
        struct kevent* e = &events[i];
        if (e->filter == EVFILT_READ) {
            // ready to read
            // Read only one time, bcz Write is send to only once..
            char buf[64];
            int m = read(pipe_fd[0], &buf, sizeof(buf) - 1);
            buf[m] = '\0';
            printf("READ DATA[%d]= %s\n", m, buf);
        } else if (e->filter == EVFILT_WRITE) {
            // ready to write
        } else {
            // do nothing
        }
    }

    printf("===========================\n");

    // wait for thread worker to finish...
    pthread_join(tid, NULL);
 
    close(e_fd);
    close(pipe_fd[0]);
    close(pipe_fd[1]);
    
	return 0;
}