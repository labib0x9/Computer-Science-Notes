package main

/*
	gRPC - Basics (server, client) ? Why am I complexing HERE... IDK ????

	go get google.golang.org/grpc
	go get google.golang.org/protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


	brew install protobuf

	// Compile
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	proto/todo.proto
*/

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	todo "test/proto"
)

type Item struct {
	Id    int32
	Title string
}

var items []Item

type Server struct {
	todo.UnimplementedTodoServer
}

func (s *Server) CreateTodo(ctx context.Context, in *todo.TodoItem) (*todo.TodoItem, error) {
	newItem := Item{
		Id:    int32(len(items) + 1),
		Title: in.Title,
	}

	items = append(items, newItem)
	return &todo.TodoItem{Id: newItem.Id, Title: newItem.Title}, nil
}

func (s *Server) ReadTodo(ctx context.Context, in *todo.VoidNoParameter) (*todo.TodoItems, error) {
	var todoList []*todo.TodoItem
	for _, item := range items {
		todoList = append(todoList, &todo.TodoItem{
			Id:    item.Id,
			Title: item.Title,
		})
	}
	return &todo.TodoItems{Item: todoList}, nil
}

func (s *Server) ReadTodoStream(in *todo.VoidNoParameter, stream todo.Todo_ReadTodoStreamServer) error {
	for _, item := range items {
		if err := stream.Send(&todo.TodoItem{
			Id:    item.Id,
			Title: item.Title,
		}); err != nil {
			return err
		}
	}
	return nil
}

func server() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	todo.RegisterTodoServer(grpcServer, &Server{})

	fmt.Println("gRPC server listening on :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func client(title string) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	c := todo.NewTodoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. CreateTodo RPC
	fmt.Println("--- Creating Todos ---")
	item1, err := c.CreateTodo(ctx, &todo.TodoItem{Title: title})
	if err != nil {
		log.Fatalf("could not create todo: %v", err)
	}
	fmt.Printf("Created: ID=%d, Title=%s\n", item1.GetId(), item1.GetTitle())

	// 2. ReadTodo (Unary RPC)
	fmt.Println("\n--- Reading All Todos (Unary) ---")
	todos, err := c.ReadTodo(ctx, &todo.VoidNoParameter{})
	if err != nil {
		log.Fatalf("could not read todos: %v", err)
	}
	for _, item := range todos.GetItem() {
		fmt.Printf("Todo: ID=%d, Title=%s\n", item.GetId(), item.GetTitle())
	}

	// 3. ReadTodoStream (Server Streaming RPC)
	fmt.Println("\n--- Reading Todos Stream ---")
	stream, err := c.ReadTodoStream(ctx, &todo.VoidNoParameter{})
	if err != nil {
		log.Fatalf("could not open stream: %v", err)
	}

	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving stream: %v", err)
		}
		fmt.Printf("Streamed Todo: ID=%d, Title=%s\n", item.GetId(), item.GetTitle())
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("go run main.go server|client")
		return
	}

	switch args[0] {
	case "client":
		client(args[1])
	case "server":
		server()
	}
}
