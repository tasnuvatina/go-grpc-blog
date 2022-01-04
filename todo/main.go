package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/spf13/viper"

	"github.com/tasnuvatina/grpc-blog/todo/services/todo"
	"github.com/tasnuvatina/grpc-blog/todo/services/user"
	"github.com/tasnuvatina/grpc-blog/todo/services/blog"

	tpb "github.com/tasnuvatina/grpc-blog/proto/todo"
	upb "github.com/tasnuvatina/grpc-blog/proto/user"
	bpb "github.com/tasnuvatina/grpc-blog/proto/blog"

	tc "github.com/tasnuvatina/grpc-blog/todo/core/todo"
	uc "github.com/tasnuvatina/grpc-blog/todo/core/user"
	bc "github.com/tasnuvatina/grpc-blog/todo/core/blog"
	"github.com/tasnuvatina/grpc-blog/todo/storage/postgres"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("todo/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error loading configuration: %v", err)
	}

	grpcServer := grpc.NewServer()
	store, err := newDBFromConfig(config)
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}
	// todo store service
	cs := tc.NewCoreSvc(store)
	s := todo.NewTodoServer(cs)

	//user store service
	ucs := uc.NewUserCoreSvc(store)
	us := user.NewTodoServer(ucs)

	//blog store service
	bcs := bc.NewBlogCoreSvc(store)
	bs := blog.NewBlogServer(bcs)


	tpb.RegisterTodoServiceServer(grpcServer, s)
	upb.RegisterTodoServiceServer(grpcServer,us)
	bpb.RegisterBlogServiceServer(grpcServer,bs)

	
	host, port := config.GetString("server.host"), config.GetString("server.port")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	log.Printf("Server is starting on: http://%s:%s\n", host, port)

	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func newDBFromConfig(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}



