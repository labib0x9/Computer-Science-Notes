#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<string.h>
#include<stdbool.h>
#include<sys/socket.h>
#include<netinet/in.h>
#include<arpa/inet.h>

int main() {

	// create a socket endpoint
	int fd = socket(AF_INET, SOCK_STREAM, 0);
	if (fd < 0) {
		perror("socket create");
		return 0;
	}
	printf("%d\n", fd);

	// struct sockaddr_storage addr; // ipv4 + ipv6

	struct sockaddr_in addr;
	socklen_t addr_len = sizeof(addr);
	// bzero(&addr, sizeof(addr));	// fill with zero

	memset(&addr, 0, sizeof(addr));

	addr.sin_family = AF_INET;	// ipv4
	addr.sin_port = htons(8080);	// port -> endianess handles
	addr.sin_addr.s_addr = INADDR_ANY;	// accept any local interface

	// // bind to a specific ip address
	// if (inet_pton(AF_INET, "127.0.0.1", &addr.sin_addr) <= 0) {
	// 	perror("inet_pton");
	// 	return 0;
	// }

	if (bind(fd, (struct sockaddr*)&addr, addr_len) < 0 ) { close(fd); perror("bind"); return 0; };	// bind to port 127.0.0.1:8080
	if (listen(fd, 10) < 0) { close(fd); perror("listen"); return 0; }; // 10 backlog queue, why use backlog queue ? maximum number of pending connections waiting to be accepted

	int recieved;
	char request[1024];

	while (1) {
		int c_fd = accept(fd, NULL, NULL);	// blocks until a new client connects ..

		if (c_fd < 0) {
			continue;
		}

		printf("client = %d\n", c_fd);
		char x[] = "HELLO CLIENT";

		recieved = recv(c_fd, request, sizeof(request), 0);	// receieve message ..
		request[recieved] = '\0';
		
		send(c_fd, x, sizeof(x), 0);	// sends message ..

		close(c_fd);
		break;
	}

	close(fd);

	printf("%s\n", request);

	return 0;
}
