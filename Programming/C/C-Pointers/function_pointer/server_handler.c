#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<string.h>
#include<stdbool.h>

typedef int (*handler_func)(int client_id, char* body);
// typedef void (*route_handler_fn)(client_t*, request_ctx_t*);

typedef struct {
	char* method;
	char* route;
	handler_func handler;
} Route;

void register_route(char* method, char* route, handler_func func) {
	Route x = {method, route, func};
}

int handle_home(int client_id, char* body) {
	printf("Hello\n");
	return 0;
}

int handle_login(int client_id, char* body) {
	printf("LOGIN\n");
	return 0;
}

int main() {

	register_route("GET", "/", handle_home);
	register_route("GET", "/login", handle_login);

	return 0;
}
