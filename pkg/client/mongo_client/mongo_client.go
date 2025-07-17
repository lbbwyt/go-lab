package mongo_client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient 封装MongoDB客户端
type MongoDBClient struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewMongoDBClient 创建新的MongoDB客户端
func NewMongoDBClient(uri, dbName, collName string) (*MongoDBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("连接MongoDB失败: %v", err)
	}

	// 验证连接
	err = client.Ping(ctx, nil)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("验证MongoDB连接失败: %v", err)
	}

	db := client.Database(dbName)
	coll := db.Collection(collName)

	return &MongoDBClient{
		client:     client,
		database:   db,
		collection: coll,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

// Close 关闭连接
func (m *MongoDBClient) Close() error {
	if m.cancel != nil {
		m.cancel()
	}
	return m.client.Disconnect(m.ctx)
}

// InsertOne 插入单个文档
func (m *MongoDBClient) InsertOne(document interface{}) (primitive.ObjectID, error) {
	result, err := m.collection.InsertOne(m.ctx, document)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("插入文档失败: %v", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}

	return primitive.NilObjectID, errors.New("无法获取插入文档的ID")
}

// InsertMany 插入多个文档
func (m *MongoDBClient) InsertMany(documents []interface{}) ([]primitive.ObjectID, error) {
	result, err := m.collection.InsertMany(m.ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("批量插入文档失败: %v", err)
	}

	ids := make([]primitive.ObjectID, 0, len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		if oid, ok := id.(primitive.ObjectID); ok {
			ids = append(ids, oid)
		}
	}

	return ids, nil
}

// FindOne 查询单个文档
func (m *MongoDBClient) FindOne(filter interface{}, result interface{}) error {
	err := m.collection.FindOne(m.ctx, filter).Decode(result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("未找到匹配的文档")
		}
		return fmt.Errorf("查询文档失败: %v", err)
	}
	return nil
}

// FindMany 查询多个文档
func (m *MongoDBClient) FindMany(filter interface{}, opts ...*options.FindOptions) ([]bson.M, error) {
	cursor, err := m.collection.Find(m.ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("查询文档失败: %v", err)
	}
	defer cursor.Close(m.ctx)

	var results []bson.M
	if err = cursor.All(m.ctx, &results); err != nil {
		return nil, fmt.Errorf("解析查询结果失败: %v", err)
	}

	if len(results) == 0 {
		return nil, errors.New("未找到匹配的文档")
	}

	return results, nil
}

// UpdateOne 更新单个文档
func (m *MongoDBClient) UpdateOne(filter, update interface{}) (int64, error) {
	result, err := m.collection.UpdateOne(m.ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("更新文档失败: %v", err)
	}
	return result.ModifiedCount, nil
}

// UpdateMany 更新多个文档
func (m *MongoDBClient) UpdateMany(filter, update interface{}) (int64, error) {
	result, err := m.collection.UpdateMany(m.ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("批量更新文档失败: %v", err)
	}
	return result.ModifiedCount, nil
}

// DeleteOne 删除单个文档
func (m *MongoDBClient) DeleteOne(filter interface{}) (int64, error) {
	result, err := m.collection.DeleteOne(m.ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("删除文档失败: %v", err)
	}
	return result.DeletedCount, nil
}

// DeleteMany 删除多个文档
func (m *MongoDBClient) DeleteMany(filter interface{}) (int64, error) {
	result, err := m.collection.DeleteMany(m.ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("批量删除文档失败: %v", err)
	}
	return result.DeletedCount, nil
}

// Count 统计文档数量
func (m *MongoDBClient) Count(filter interface{}) (int64, error) {
	count, err := m.collection.CountDocuments(m.ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计文档数量失败: %v", err)
	}
	return count, nil
}

// GetCollection 获取集合实例
func (m *MongoDBClient) GetCollection() *mongo.Collection {
	return m.collection
}
