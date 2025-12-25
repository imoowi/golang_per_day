package protoimpl

import (
	"context"
	"golang_per_day_29/proto/goods"
)

type GoodsServiceServer struct {
	goods.UnimplementedGoodsServiceServer
}

func (s *GoodsServiceServer) GetGoods(ctx context.Context, req *goods.GetGoodsRequest) (*goods.GetGoodsResponse, error) {
	return &goods.GetGoodsResponse{
		Id:    req.Id,
		Name:  "Codee君",
		Price: 99.99,
		Stock: 1000,
	}, nil
}
