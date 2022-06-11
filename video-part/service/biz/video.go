package biz

import (
	"context"
	"yunyandz.com/tiktok/proto/pb"
)

type VideoServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (v VideoServiceServer) Feed(ctx context.Context, req *pb.FeedReq) (*pb.FeedRes, error) {
	//TODO implement me
	panic("implement me")
}

func (v VideoServiceServer) Publish(ctx context.Context, req *pb.PublishReq) (*pb.PublishRes, error) {

	//TODO implement me
	panic("implement me")
}

func (v VideoServiceServer) PublishList(ctx context.Context, req *pb.PublishListReq) (*pb.PublishListRes, error) {
	//TODO implement me
	panic("implement me")
}

func (v VideoServiceServer) Favorite(ctx context.Context, req *pb.FavoriteReq) (*pb.FavoriteRes, error) {
	//TODO implement me
	panic("implement me")
}

func (v VideoServiceServer) FavoriteList(ctx context.Context, req *pb.FavoriteListReq) (*pb.FavoriteListRes, error) {
	//TODO implement me
	panic("implement me")
}

func (v VideoServiceServer) Comment(ctx context.Context, req *pb.CommentReq) (*pb.CommentRes, error) {
	//TODO implement me
	panic("implement me")
}

func (v VideoServiceServer) CommentList(ctx context.Context, req *pb.CommentListReq) (*pb.CommentListRes, error) {
	//TODO implement me
	panic("implement me")
}
