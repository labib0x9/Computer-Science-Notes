
#include<stdio.h>
#include<string.h>
#include<llhttp.h>

const char CONN[] =
    "GET /text HTTP/1.1\r\n"
    "Host: 127.0.0.1:8080\r\n"
    "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:146.0) Gecko/20100101 Firefox/146.0\r\n"
    "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n"
    "Accept-Language: en-US,en;q=0.5\r\n"
    "Accept-Encoding: gzip, deflate, br, zstd\r\n"
    "Connection: keep-alive\r\n"
    "Upgrade-Insecure-Requests: 1\r\n"
    "Sec-Fetch-Dest: document\r\n"
    "Sec-Fetch-Mode: navigate\r\n"
    "Sec-Fetch-Site: none\r\n"
    "Sec-Fetch-User: ?1\r\n"
    "Priority: u=0, i\r\n"
    "\r\n";

const char CONN_P[] =
    "POST /submit HTTP/1.1\r\n"
    "Host: 127.0.0.1:8080\r\n"
    "Content-Length: 11\r\n"
    "\r\n"
    "hello world";


// stores headers as key-value
typedef struct {
    char key[128];
    char value[128];
} header_t;

// stores parsed request
// multiple value in single key, stores as one entity. 
typedef struct Request {
    char method[8];
    char path[256];
    int path_len;
    header_t headers[128];
    int header_count;
    header_t cur_header;
    char* body; // heap allocated
} http_request_t;

// llhttp = A byte-stream HTTP/1.x state machine
// handles partial inputs
// recv() → bytes → llhttp_execute() → callbacks → your state

// set all memory to 0.
void http_request_t_init(http_request_t* req) {
    memset(req, 0, sizeof(*req));
}

// define callbacks, Callbacks are the only output mechanism.

// Get the method
static int on_method(llhttp_t* parser, const char* at, size_t len) {
    http_request_t *r = parser->data;

    // switch(parser->method) {
    //     case HTTP_GET: {
    //         strcpy(r->method, "GET");
    //         break;
    //     }
    //     case HTTP_POST: {
    //         strcpy(r->method, "POST");
    //         break;
    //     }
    //     default: {
    //         strcpy(r->method, "OTHER");
    //     }
    // }
    memcpy(r->method, at, len);
    r->method[len] = '\0';
    return 0;
}

// at = pointer from where path begins.
// length = path length
static int on_url(llhttp_t* parser, const char* at, size_t length) {
    // fwrite(at, 1, length, stdout);
    // (void) parser;
    http_request_t *r = parser->data;
    // (void) r;
    // (void) at;
    // (void) length;

    // char path[length + 1];
    // memcpy(path, at, length);
    // path[length] = '\0';

    // printf("PATH = %s\n", path);
    // printf("URL(on_url) = %s\n\nLength = %zu\n", at, length);

    memcpy(r->path, at, length);
    r->path[length] = '\0';

    return 0;
}

// Gives the header fields only, store it in cur_header variable to access in on_header_value()
static int on_header_field(llhttp_t* parser, const char* at, size_t length) {
    // printf("Header field: %.*s\n", (int)length, at);
    // (void) parser;
    http_request_t *r = parser->data;
    // (void) r;
    // (void) at;
    // (void) length;

    memcpy(r->cur_header.key, at, length);
    r->cur_header.key[length] = '\0';

    return 0;
}

// Gives the header values, parse it for multiple types seperated by ;
static int on_header_value(llhttp_t* parser, const char* at, size_t length) {
    // printf("Header value: %.*s\n", (int)length, at);
    // (void) parser;
    http_request_t *r = parser->data;
    // (void) r;
    // (void) at;
    // (void) length;

    memcpy(r->cur_header.value, at, length);
    r->cur_header.value[length] = '\0';

    memcpy(&r->headers[r->header_count], &r->cur_header, sizeof(r->cur_header));
    r->header_count++;

    return 0;
}

static int on_header_complete(llhttp_t* parser) {
    // printf("Request complete\n");
    // (void) parser;
    http_request_t *r = parser->data;
    (void) r;
    return 0;
}

static int on_body(llhttp_t* parser, const char* at, size_t length) {
    http_request_t *r = parser->data;
    (void) r;
    (void) at;
    (void) length;
    return 0;
}

static int on_message_complete(llhttp_t* parser) {
    // printf("Request complete\n");
    // (void) parser;
    http_request_t *r = parser->data;
    (void) r;
    return 0;
}

int main() {

    // request struct
    http_request_t req;
    http_request_t_init(&req);

    // http parser and settings
    llhttp_t parser;
    llhttp_settings_t settings;

    // initialize settings
    llhttp_settings_init(&settings);

    // defines callbacks on settings
    settings.on_url = on_url;
    settings.on_header_field = on_header_field;
    settings.on_header_value = on_header_value;
    settings.on_message_complete = on_message_complete;
    settings.on_body = on_body;
    settings.on_headers_complete = on_header_complete;
    settings.on_method = on_method;

    // initialize parser with settings
    llhttp_init(&parser, HTTP_REQUEST, &settings);
    parser.data = &req;

    // But, how can i store parsed req into http_request_t struct ???
    llhttp_errno_t err = llhttp_execute(&parser, CONN, strlen(CONN));

    // invalid input, you should close the connection.
    // 400 Bad Request
    if (err != HPE_OK) {
        fprintf(stderr, "Parse error: %s\n", llhttp_errno_name(err));
    }

    if (parser.finish) {
        printf("finished\n");
    }

    printf("PATH = %s\n", req.path);
    printf("HEADER-LEN = %d\n", req.header_count);
    printf("2nd header KEY = %s\n", req.headers[1].key);
    printf("2nd header VALUE = %s\n", req.headers[1].value);
    printf("MATHOD = %s\n", req.method);


    // // Handle Connection: Keep-Alive 
    // int keep_alive = llhttp_should_keep_alive(&parser);
    // if (keep_alive) {
    //     // reset the parser, but how ?
    //     llhttp_init(&parser, HTTP_REQUEST, &settings);
    //     // parser.data = ?
    // } else {
    //     // close the connection.
    // }


    return 0;
}