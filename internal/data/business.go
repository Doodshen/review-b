package data

import (
	"context"

	v1 "review-b/api/review/v1"
	"review-b/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type businessRepo struct {
	data *Data
	log  *log.Helper
}

// NewBusinessRepo .
func NewBusinessRepo(data *Data, logger log.Logger) biz.BusinessRepo {
	return &businessRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Reply 调用grpc服务，获取响应ID
func (r *businessRepo) Reply(ctx context.Context, param *biz.ReplyParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data] Reply param:%v", param)
	//之前是通过查询数据库来创建服务

	//现在需要通过rpc调用其他服务
	ret, err := r.data.rc.ReplyReview(ctx, &v1.ReplyReviewRequest{ //填写的是review中的v1 也就是review-service定义的v1 也就是调用这个服务的客户端
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	//写日志
	r.log.WithContext(ctx).Debugf("ReplyReview return ,ret :%v", err)
	if err != nil {
		return 0, err
	}
	return ret.GetReplyID(), nil
}
