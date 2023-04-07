package storage

import (
	"bnsportal/models"
	"context"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func (s *Storage) CreateNameIndex() error {
	model := []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "name", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bsonx.Doc{{Key: "id", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bsonx.Doc{{Key: "owner", Value: bsonx.Int32(1)}},
		},
	}
	_, err := mgm.Coll(&models.RegisteredNameInfo{}).Indexes().CreateMany(context.Background(), model)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateNameInfo(name *models.RegisteredNameInfo) error {
	name.Creating()
	_, err := s.mongo.Database(s.dbName).Collection(mgm.CollName(name)).InsertOne(nil, name)
	return err
}

func (s *Storage) GetNameInfoByID(id string) (*models.RegisteredNameInfo, error) {
	name := &models.RegisteredNameInfo{}
	result := s.mongo.Database(s.dbName).Collection(mgm.CollName(name)).FindOne(context.Background(), bson.M{
		"id": id,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	err := result.Decode(name)
	if err != nil {
		return nil, err
	}
	return name, nil
}

func (s *Storage) UpdateNameInfo(name *models.RegisteredNameInfo) error {
	name.Saving()
	_, err := s.mongo.Database(s.dbName).Collection(mgm.CollName(name)).UpdateOne(context.Background(), bson.M{
		"id": name.ID,
	}, bson.M{
		"$set": bson.M{
			"owner":      name.Owner,
			"updated_at": name.UpdatedAt,
		},
	})
	return err
}

func (s *Storage) GetAddressNames(address string) ([]models.RegisteredNameInfo, error) {
	var result []models.RegisteredNameInfo
	cursor, err := s.mongo.Database(s.dbName).Collection(mgm.CollName(&models.RegisteredNameInfo{})).Find(context.Background(), bson.M{
		"owner": bson.D{{"$regex", primitive.Regex{Pattern: address, Options: "i"}}},
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Storage) GetNames(limit, offset int64) ([]models.RegisteredNameInfo, error) {
	var result []models.RegisteredNameInfo
	cursor, err := s.mongo.Database(s.dbName).Collection(mgm.CollName(&models.RegisteredNameInfo{})).Find(context.Background(), bson.M{}, &options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
		Sort:  bson.M{"registered_at_block": -1},
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Storage) CheckNameAvailable(name string) (bool, error) {
	nameInfo := &models.RegisteredNameInfo{}
	result := s.mongo.Database(s.dbName).Collection(mgm.CollName(nameInfo)).FindOne(context.Background(), bson.M{
		"name": name,
	})
	if result.Err() != nil {
		if mongo.ErrNoDocuments == result.Err() {
			return true, nil
		}
	}
	err := result.Decode(nameInfo)
	if err != nil {
		return false, err
	}
	return false, nil
}

func (s *Storage) FilterName(limit, offset int64, filterMap map[string]interface{}) ([]models.RegisteredNameInfo, error) {
	var filter bson.M
	filter = make(bson.M)
	if len(filterMap) > 0 {
		for k, v := range filterMap {
			filter[k] = v
		}
	}
	result, err := s.mongo.Database(s.dbName).Collection(mgm.CollName(&models.RegisteredNameInfo{})).Find(nil, filter, options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	var names []models.RegisteredNameInfo
	err = result.All(nil, &names)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (s *Storage) GetNameInfo(name string) (*models.RegisteredNameInfo, error) {
	nameInfo := &models.RegisteredNameInfo{}
	result := s.mongo.Database(s.dbName).Collection(mgm.CollName(nameInfo)).FindOne(context.Background(), bson.M{
		"name": name,
	})
	if result.Err() != nil {
		if mongo.ErrNoDocuments == result.Err() {
			return nil, nil
		}
	}
	err := result.Decode(nameInfo)
	if err != nil {
		return nil, err
	}
	return nameInfo, nil
}
