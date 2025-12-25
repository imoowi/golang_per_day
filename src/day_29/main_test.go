package main

import (
	"context"
	"golang_per_day_29/proto/goods"
	"testing"

	"google.golang.org/grpc"
)

func TestGrpc(t *testing.T) {
	conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()
	client := goods.NewGoodsServiceClient(conn)
	res, _ := client.GetGoods(context.Background(), &goods.GetGoodsRequest{Id: 1})
	if res.GetId() != 1 {
		t.Error("grpc client error")
	}
}
