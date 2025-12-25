package main

import (
	"golang_per_day_29/proto/goods"
	protoimpl "golang_per_day_29/proto_impl"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	goods.RegisterGoodsServiceServer(s, &protoimpl.GoodsServiceServer{})
	s.Serve(lis)
}
