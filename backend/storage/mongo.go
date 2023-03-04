package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	dexDB  *mongo.Database
	ctx    context.Context
	client *mongo.Client
}

const mongoURI string = "mongodb://localhost:27017"

func NewMongoStorage() *MongoStorage {
	return &MongoStorage{}
}

func (s *MongoStorage) DisconnectClient() {
	s.client.Disconnect(s.ctx)
}

func (s *MongoStorage) ConnectMongoDB() {
	s.ctx = context.Background()
	client, err := mongo.Connect(s.ctx, options.Client().ApplyURI(mongoURI))
	s.client = client
	if err != nil {
		log.Fatal(err)
	}
	s.dexDB = client.Database("dex_2")
}

func (s *MongoStorage) InsertOne(data interface{}, collectionString string) (*mongo.InsertOneResult, error) {
	collection := s.dexDB.Collection(collectionString)
	res, err := collection.InsertOne(s.ctx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *MongoStorage) InsertMany(data []interface{}, collectionString string) (*mongo.InsertManyResult, error) {
	collection := s.dexDB.Collection(collectionString)
	res, err := collection.InsertMany(s.ctx, []interface{}{data})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *MongoStorage) Find(search *bson.D, collectionString string) ([]map[string]interface{}, error) {
	collection := s.dexDB.Collection(collectionString)
	cursor, err := collection.Find(s.ctx, search)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	var data []map[string]interface{}
	if err = cursor.All(s.ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *MongoStorage) FindOne(search *bson.M, collectionString string) (map[string]interface{}, error) {
	collection := s.dexDB.Collection(collectionString)
	var data bson.M
	if err := collection.FindOne(s.ctx, search).Decode(&data); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

func (s *MongoStorage) Update(search *bson.M, update bson.M, collectionString string) (*mongo.UpdateResult, error) {
	collection := s.dexDB.Collection(collectionString)
	updated, err := collection.UpdateOne(s.ctx, search, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
