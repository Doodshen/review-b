package data

import (
	"context"
	v1 "review-b/api/review/v1"
	"review-b/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

// ProviderSet is data providers.新增一个构造函数要写入wire中
var ProviderSet = wire.NewSet(NewReviewService, NewData, NewBusinessRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	//通过嵌入一个gRPC的客户端，然后通过这个client去调用review-service
	rc  v1.ReviewClient
	log *log.Helper
}

// NewData 初始化data结构体
func NewData(c *conf.Data, rc v1.ReviewClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		rc:  rc,
		log: log.NewHelper(logger),
	}, cleanup, nil
}

// NewReviewService()  创建一个Review-service的gRPC客户端
func NewReviewService() v1.ReviewClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9092"),
		//添加中间键
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	)
	if err != nil {
		panic(err)
	}
	return v1.NewReviewClient(conn)
}
