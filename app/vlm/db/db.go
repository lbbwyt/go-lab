package db

import (
	"go-lab/app/vlm/model"
	"go-lab/pkg/client/mongo_client"
	"go.mongodb.org/mongo-driver/bson"
)

type MultiModelDataDao struct {
	db *mongo_client.MongoDBClient
}

func NewMultiModelDataDao(d *mongo_client.MongoDBClient) *MultiModelDataDao {
	return &MultiModelDataDao{db: d}
}

func (d *MultiModelDataDao) FindOneByTmp(tmp int64) (*model.MultiModelData, error) {
	filter := bson.M{"timestep": tmp}
	m := &model.MultiModelData{}
	err := d.db.FindOne(filter, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (d *MultiModelDataDao) UpdateLabel(tmp int64, labels string) (int64, error) {
	filter := bson.M{"timestep": tmp}
	update := bson.M{"$set": bson.M{"label": labels}}
	return d.db.UpdateMany(filter, update)
}
