package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"github.com/tasnuvatina/grpc-blog/cms/handler"
	bpb "github.com/tasnuvatina/grpc-blog/proto/blog"
	tpb "github.com/tasnuvatina/grpc-blog/proto/todo"
	upb "github.com/tasnuvatina/grpc-blog/proto/user"
	"google.golang.org/grpc"
)

func main() {
	// getting data from evn/config file
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("cms/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error loading configuration: %v", err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	store := sessions.NewCookieStore([]byte(config.GetString("session.secret")))

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", config.GetString("todo.host"), config.GetString("todo.port")), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	tc := tpb.NewTodoServiceClient(conn)
	uc := upb.NewTodoServiceClient(conn)
	bc := bpb.NewBlogServiceClient(conn)
	r := handler.New(decoder, store, tc, uc, bc)

	host, port := config.GetString("server.host"), config.GetString("server.port")
	log.Println("Cleint starting ........")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal(err)
	}
}
