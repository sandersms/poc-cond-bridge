// SPDX-License-Identifier: Apache-2.0

// main is the main package of the conditional bridge application
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/opiproject/opi-spdk-bridge/pkg/utils"
	"github.com/sandersms/Protos/Cond-bridge/pkg/inventory"

	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/redis"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	var grpcPort int
	flag.IntVar(&grpcPort, "grpc_port", 50051, "The gRPC server port")

	var httpPort int
	flag.IntVar(&httpPort, "http_port", 8082, "The HTTP server port")

	var spdkAddress string
	flag.StringVar(&spdkAddress, "spdk_addr", "/var/tmp/spdk.sock", "Points to SPDK unix socket/tcp socket to interact with")

	var tlsFiles string
	flag.StringVar(&tlsFiles, "tls", "", "TLS files in server_cert:server_key:ca_cert format.")

	var redisAddress string
	flag.StringVar(&redisAddress, "redis_addr", "127.0.0.1:6379", "Redis address in ip_address:port format")

	flag.Parse()

	// Create KV store for persistence
	options := redis.DefaultOptions
	options.Codec = utils.ProtoCodec{}
	options.Address = redisAddress
	store, err := redis.NewClient(options)
	if err != nil {
		log.Panic(err)
	}
	defer func(store gokv.Store) {
		err := store.Close()
		if err != nil {
			log.Panic(err)
		}
	}(store)

	go runGatewayServer(grpcPort, httpPort)
	runGrpcServer(grpcPort)
}

func runGrpcServer(grpcPort int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	inventory.RegisterInventorytoGrpc(s)
	reflection.Register(s)

	log.Printf("gRPC Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runGatewayServer(grpcPort int, httpPort int) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("localhost:%d", grpcPort)
	inventory.RegisterInventorytoGateway(ctx, mux, endpoint, opts)
	

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Printf("HTTP Server listening at %v", httpPort)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", httpPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Panic("cannot start HTTP gateway server")
	}
}
