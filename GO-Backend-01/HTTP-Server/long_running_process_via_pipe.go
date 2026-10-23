package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

/*

labib0x9@TestGoCode $  curl -i http://localhost:8080/job
HTTP/1.1 200 OK
Content-Type: text/plain
X-Content-Type-Options: nosniff
Date: Fri, 23 Oct 2026 16:18:09 GMT
Transfer-Encoding: chunked			// added by flush()

{"id": 0, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 1, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 2, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 3, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 4, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 5, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 6, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 7, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 8, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 9, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 10, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 11, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 12, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 13, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 14, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 15, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 16, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 17, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 18, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 19, "user_ip": "127.0.0.1", "event": "testing_on_background"}


labib0x9@TestGoCode $  curl -i http://localhost:8080/job -v
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /job HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: *\/*
>
* Request completely sent off
< HTTP/1.1 200 OK
HTTP/1.1 200 OK
< Content-Type: text/plain
Content-Type: text/plain
< X-Content-Type-Options: nosniff
X-Content-Type-Options: nosniff
< Date: Fri, 23 Oct 2026 16:19:25 GMT
Date: Fri, 23 Oct 2026 16:19:25 GMT
< Transfer-Encoding: chunked
Transfer-Encoding: chunked
<

{"id": 0, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 1, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 2, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 3, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 4, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 5, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 6, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 7, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 8, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 9, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 10, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 11, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 12, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 13, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 14, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 15, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 16, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 17, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 18, "user_ip": "127.0.0.1", "event": "testing_on_background"}
{"id": 19, "user_ip": "127.0.0.1", "event": "testing_on_background"}
* Connection #0 to host localhost left intact

*/

func longRunningProcess(logWriter *io.PipeWriter) {
	for i := 0; i < 20; i++ {
		fmt.Fprintf(logWriter, `{"id": %d, "user_ip": "127.0.0.1", "event": "testing_on_background"}`, i)
		fmt.Fprintln(logWriter)
		time.Sleep(time.Second)
	}
	logWriter.Close()
}

func processStreamer(logReader *io.PipeReader, w http.ResponseWriter, done chan struct{}) {
	buf := make([]byte, 500)
	f, flushSupported := w.(http.Flusher)

	defer logReader.Close()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Content-Type-Options", "nosniff") // what does it do?

	for {
		n, err := logReader.Read(buf)
		if err == io.EOF {
			break
		}
		w.Write(buf[:n])
		if flushSupported {
			f.Flush()
		}
	}
	done <- struct{}{}
}

func longRunningProcessHandler(w http.ResponseWriter, r *http.Request) {
	done := make(chan struct{})
	logReader, logWriter := io.Pipe()
	go longRunningProcess(logWriter)
	go processStreamer(logReader, w, done)

	<-done
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()
	mux.HandleFunc("GET /job", longRunningProcessHandler)

	log.Fatal(http.ListenAndServe(addr, mux))
}
