#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<string.h>
#include<stdbool.h>
#include<sys/socket.h>
#include<netinet/in.h>
#include<arpa/inet.h>
#include<errno.h>

/*
    -> Client address
    -> Bind to specific address
    -> Socket behave setting
    -> timeout on client socket fd
*/

typedef struct Client {
	int fd;
	struct sockaddr_storage addr;
} client_t;

void handle_conn(client_t client) {

    // set recv() timeout on client socket fd
    struct timeval tv = {.tv_sec = 5, .tv_usec = 0};    // 5seconds
    if( setsockopt(client.fd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv)) < 0) {
        perror("setsockopt SO_RCVTIMEO");
        close(client.fd);
        return;
    }

    // ipv4, client address print
    if (client.addr.ss_family == AF_INET) {
        // // Method - 01
        // struct sockaddr_in temp_client;
        // memcpy(&temp_client, &client.addr, sizeof(temp_client));
        // char client_ip[INET_ADDRSTRLEN];
        // inet_ntop(AF_INET, &temp_client->sin_addr, client_ip, INET_ADDRSTRLEN);
        // printf("[CLIENT] %s:%d\n", client_ip, htons(temp_client->sin_port));

        // // Method - 02
        struct sockaddr_in *temp_client = (struct sockaddr_in*) &client.addr;
        char client_ip[INET_ADDRSTRLEN];
        inet_ntop(AF_INET, &temp_client->sin_addr, client_ip, INET_ADDRSTRLEN);
        printf("[CLIENT] %s:%d\n", client_ip, htons(temp_client->sin_port));
    }

    int recieved;
	char request[1024];

    printf("client = %d\n", client.fd);
	char x[] = "HELLO CLIENT";

	while(1) {
        recieved = recv(client.fd, request, sizeof(request) - 1, 0);	// receieve message ..
        if (recieved <= 0) {
            printf("%d bytes received\n", recieved);

            // client disconnect
            if (errno == EAGAIN || errno == EWOULDBLOCK) {
                printf("EAGAIN\n");
            }

            break;
        }
        request[recieved] = '\0';
        
        // sends message ..
        recieved = send(client.fd, x, strlen(x), MSG_NOSIGNAL);
        if (recieved <= 0) {
            printf("%d bytes send\n", recieved);
            break;
        }
    }

    printf("client disconnected\n");
	close(client.fd);
}

int main() {

    // ignore sigpipe
	signal(SIGPIPE, SIG_IGN);

	int fd = socket(AF_INET, SOCK_STREAM, 0);
	if (fd < 0) {
		perror("socket create");
		return 0;
	}
	printf("%d\n", fd);

	struct sockaddr_in s_addr;
	socklen_t addr_len = sizeof(s_addr);

	memset(&s_addr, 0, sizeof(s_addr));

	s_addr.sin_family = AF_INET;
	s_addr.sin_port = htons(8080);

	// bind to a specific ip address
	if (inet_pton(AF_INET, "127.0.0.1", &s_addr.sin_addr) <= 0) {
		perror("inet_pton");
		return 0;
	}

    // Control socket behaviour
    int opt = 1;
    // reuse address
    if (setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) < 0) {
        perror("setsockopt SO_REUSEADDR");
        close(fd);
        return 0;
    }

	if (bind(fd, (struct sockaddr*)&s_addr, addr_len) < 0 ) { close(fd); perror("bind"); return 0; };	// bind to port 127.0.0.1:8080
	if (listen(fd, 10) < 0) { close(fd); perror("listen"); return 0; }; // 10 backlog queue, why use backlog queue ? maximum number of pending connections waiting to be accepted

	while(1) {
        // store client info
        client_t client;
        socklen_t len = sizeof(client.addr);
		int c_fd = accept(fd, (struct sockaddr*) &client.addr, &len);	// blocks until a new client connects ..

		if (c_fd < 0) {
			continue;
		}

        client.fd = c_fd;
        handle_conn(client);
	}

	close(fd);

	return 0;
}