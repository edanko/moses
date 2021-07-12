package mongo

import (
	"context"
	"time"

	"github.com/edanko/moses/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RemnantRepository interface {
	GetOne(id string) (*entities.Remnant, error)
	GetAll() ([]*entities.Remnant, error)
	GetNotUsed(project, dim, quality string) ([]*entities.Remnant, error)
	GetAllNotUsed() ([]*entities.Remnant, error)

	Create(e *entities.Remnant) (*entities.Remnant, error)
	Update(e *entities.Remnant) (*entities.Remnant, error)
	Delete(id string) error
	DeleteAll() error
}

type remnantRepository struct {
	Collection *mongo.Collection
}

func NewRemnantRepo(collection *mongo.Collection) RemnantRepository {
	return &remnantRepository{
		Collection: collection,
	}
}

func (r *remnantRepository) GetOne(id string) (*entities.Remnant, error) {
	remnantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor := r.Collection.FindOne(context.Background(), bson.M{"_id": remnantID})

	var remnant entities.Remnant
	err = cursor.Decode(&remnant)
	if err != nil {
		return nil, err
	}

	return &remnant, nil
}

func (r *remnantRepository) Create(e *entities.Remnant) (*entities.Remnant, error) {
	e.ID = primitive.NewObjectID()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(context.Background(), e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *remnantRepository) GetAll() ([]*entities.Remnant, error) {
	var remnants []*entities.Remnant

	cursor, err := r.Collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var remnant entities.Remnant
		_ = cursor.Decode(&remnant)
		remnants = append(remnants, &remnant)
	}
	return remnants, nil
}

func (r *remnantRepository) GetNotUsed(project, dim, quality string) ([]*entities.Remnant, error) {
	var remnants []*entities.Remnant

	cursor, err := r.Collection.Find(context.Background(), bson.M{"project": project, "dim": dim, "quality": quality, "used": false})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var user entities.Remnant
		_ = cursor.Decode(&user)

		remnants = append(remnants, &user)
	}
	return remnants, nil
}

func (r *remnantRepository) GetAllNotUsed() ([]*entities.Remnant, error) {
	var remnants []*entities.Remnant

	cursor, err := r.Collection.Find(context.Background(), bson.M{"used": false})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var remnant entities.Remnant
		_ = cursor.Decode(&remnant)

		remnants = append(remnants, &remnant)
	}
	return remnants, nil
}

func (r *remnantRepository) Update(e *entities.Remnant) (*entities.Remnant, error) {
	e.UpdatedAt = time.Now()
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": e.ID}, bson.M{"$set": e})
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *remnantRepository) Delete(id string) error {
	remnantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(context.Background(), bson.M{"_id": remnantID})
	if err != nil {
		return err
	}
	return nil
}

func (r *remnantRepository) DeleteAll() error {
	_, err := r.Collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}
