package repository

import (
	"context"
	"errors"
	"go-framework/internal/model"
	"go-framework/util/xlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBRepository 表示数据访问对象 【MongoDB】
type MongoDBRepository struct {
	Model model.MongoDBModelImpl
	Log   *xlog.Log
}

// Options 表示查询选项
type Options struct {
	Page             int64
	PageSize         int64
	Pipeline         interface{}
	AggregateOptions *options.AggregateOptions
	FindOptions      *options.FindOptions
	FindOneOptions   *options.FindOneOptions
}

// NewMongoDBRepository 创建一个MongoDBRepository实例
func NewMongoDBRepository(model model.MongoDBModelImpl, log *xlog.Log) *MongoDBRepository {
	return &MongoDBRepository{Model: model, Log: log}
}

// QueryAll 查询符合条件的所有记录
func (r *MongoDBRepository) QueryAll(ctx context.Context, condition interface{}, dest interface{}, opts ...*Options) error {
	var err error
	var cursor *mongo.Cursor
	if len(opts) == 0 {
		cursor, err = r.Model.Model().Find(ctx, condition)
	} else {
		opt := opts[0]
		if opt.Page > 0 && opt.PageSize > 0 {
			opt.FindOptions.SetSkip((opt.Page - 1) * opt.PageSize)
			opt.FindOptions.SetLimit(opt.PageSize)
		}
		cursor, err = r.Model.Model().Find(ctx, condition, opt.FindOptions)
	}

	// 执行查询
	if err != nil {
		return err
	}

	// 将查询结果解码到目标对象中
	if err = cursor.All(ctx, dest); err != nil {
		return err
	}

	return nil
}

// QueryOne 查询符合条件的一条记录
func (r *MongoDBRepository) QueryOne(ctx context.Context, condition interface{}, dest interface{}, opts ...*Options) error {
	var err error
	if len(opts) == 0 {
		err = r.Model.Model().FindOne(ctx, condition).Decode(dest)
	} else {
		opt := opts[0]

		err = r.Model.Model().FindOne(ctx, condition, opt.FindOneOptions).Decode(dest)
	}

	// 执行查询
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	return nil
}

// Create 创建一条记录
func (r *MongoDBRepository) Create(ctx context.Context, document interface{}) (interface{}, error) {
	result, err := r.Model.Model().InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

// Update 更新符合条件的所有记录
func (r *MongoDBRepository) Update(ctx context.Context, condition interface{}, document interface{}) (interface{}, error) {

	result, err := r.Model.Model().UpdateMany(ctx, condition, document)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete 删除符合条件的所有记录
func (r *MongoDBRepository) Delete(ctx context.Context, condition interface{}) (interface{}, error) {

	result, err := r.Model.Model().DeleteMany(ctx, condition)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Count 统计符合条件的记录数
func (r *MongoDBRepository) Count(ctx context.Context, condition interface{}) (int64, error) {
	return r.Model.Model().CountDocuments(ctx, condition)
}
