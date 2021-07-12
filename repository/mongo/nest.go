package mongo

import (
	"context"
	"time"

	"github.com/edanko/moses/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NestRepository interface {
	GetOne(id string) (*entities.Nest, error)
	GetAll() ([]*entities.Nest, error)
	Get(project, dimension, quality string) ([]*entities.Nest, error)

	Create(e *entities.Nest) (*entities.Nest, error)
	Update(e *entities.Nest) (*entities.Nest, error)
	Delete(id string) error
	DeleteAll() error
}

type nestRepository struct {
	Collection *mongo.Collection
}

func NewNestRepo(collection *mongo.Collection) NestRepository {
	return &nestRepository{
		Collection: collection,
	}
}

func (n *nestRepository) GetOne(id string) (*entities.Nest, error) {
	nestID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor := n.Collection.FindOne(context.Background(), bson.M{"_id": nestID})

	var nest entities.Nest
	err = cursor.Decode(&nest)
	if err != nil {
		return nil, err
	}

	return &nest, nil
}

func (n *nestRepository) Create(e *entities.Nest) (*entities.Nest, error) {
	e.ID = primitive.NewObjectID()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	_, err := n.Collection.InsertOne(context.Background(), e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (n *nestRepository) GetAll() ([]*entities.Nest, error) {
	var nests []*entities.Nest

	cursor, err := n.Collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var user entities.Nest
		_ = cursor.Decode(&user)
		nests = append(nests, &user)
	}
	return nests, nil
}

func (n *nestRepository) Get(project, dimension, quality string) ([]*entities.Nest, error) {
	var nests []*entities.Nest

	cursor, err := n.Collection.Find(context.Background(), bson.M{"project": project, "dimension": dimension, "quality": quality})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var nest entities.Nest
		_ = cursor.Decode(&nest)

		nests = append(nests, &nest)
	}
	return nests, nil
}

func (n *nestRepository) Update(e *entities.Nest) (*entities.Nest, error) {
	e.UpdatedAt = time.Now()
	_, err := n.Collection.UpdateOne(context.Background(), bson.M{"_id": e.ID}, bson.M{"$set": e})
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (n *nestRepository) Delete(id string) error {
	nestID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = n.Collection.DeleteOne(context.Background(), bson.M{"_id": nestID})
	if err != nil {
		return err
	}
	return nil
}

func (n *nestRepository) DeleteAll() error {
	_, err := n.Collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}
