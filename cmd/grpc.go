package main

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initRPC() *grpc.Server {
	s := grpc.NewServer(
		grpc.Creds(nil),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	reflection.Register(s)
	return s
}
