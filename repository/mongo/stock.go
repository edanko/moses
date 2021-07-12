package mongo

import (
	"context"
	"time"

	"github.com/edanko/moses/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StockRepository interface {
	GetOne(dim, quality string) (*entities.Stock, error)
	GetAll() ([]*entities.Stock, error)

	Create(e *entities.Stock) (*entities.Stock, error)
	Update(e *entities.Stock) (*entities.Stock, error)
	Delete(id string) error
	DeleteAll() error
}

type stockRepository struct {
	Collection *mongo.Collection
}

func NewStockRepo(collection *mongo.Collection) StockRepository {
	return &stockRepository{
		Collection: collection,
	}
}

func (r *stockRepository) GetOne(dim, quality string) (*entities.Stock, error) {
	cursor := r.Collection.FindOne(context.Background(), bson.M{"dim": dim, "quality": quality})

	var stock entities.Stock
	err := cursor.Decode(&stock)
	if err != nil {
		return nil, err
	}

	return &stock, nil
}

func (r *stockRepository) Create(e *entities.Stock) (*entities.Stock, error) {
	e.ID = primitive.NewObjectID()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(context.Background(), e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *stockRepository) GetAll() ([]*entities.Stock, error) {
	var stocks []*entities.Stock

	cursor, err := r.Collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var stock entities.Stock
		_ = cursor.Decode(&stock)
		stocks = append(stocks, &stock)
	}
	return stocks, nil
}

func (r *stockRepository) Update(e *entities.Stock) (*entities.Stock, error) {
	e.UpdatedAt = time.Now()
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": e.ID}, bson.M{"$set": e})
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *stockRepository) Delete(id string) error {
	stockID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(context.Background(), bson.M{"_id": stockID})
	if err != nil {
		return err
	}
	return nil
}

func (r *stockRepository) DeleteAll() error {
	_, err := r.Collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}
