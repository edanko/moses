package mongo

import (
	"context"

	"github.com/edanko/moses/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpacingRepository interface {
	GetAll(machine string) ([]*entities.Spacing, error)
	GetOne(machine, dim string, e *entities.End) (*entities.Spacing, error)
	Create(e *entities.Spacing) (*entities.Spacing, error)
	Update(e *entities.Spacing) (*entities.Spacing, error)
	Delete(id string) error
	DeleteAll() error
}

type spacingRepository struct {
	Collection *mongo.Collection
}

func NewSpacingRepo(collection *mongo.Collection) SpacingRepository {
	return &spacingRepository{
		Collection: collection,
	}
}

func (r *spacingRepository) Create(e *entities.Spacing) (*entities.Spacing, error) {
	e.ID = primitive.NewObjectID()

	_, err := r.Collection.InsertOne(context.Background(), e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *spacingRepository) GetAll(machine string) ([]*entities.Spacing, error) {
	var spacings []*entities.Spacing

	cursor, err := r.Collection.Find(context.Background(), bson.M{"machine": machine})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var spacing entities.Spacing
		_ = cursor.Decode(&spacing)
		spacings = append(spacings, &spacing)
	}
	return spacings, nil
}

func (r *spacingRepository) GetOne(machine, dim string, e *entities.End) (*entities.Spacing, error) {
	hasBevel := e.WebBevel != nil || e.FlangeBevel != nil
	hasScallop := e.Scallop != nil

	cursor := r.Collection.FindOne(context.Background(), bson.M{"machine": machine, "dim": dim, "name": e.Name, "has_bevel": hasBevel, "has_scallop": hasScallop})

	var spacing entities.Spacing
	err := cursor.Decode(&spacing)
	if err != nil {
		return nil, err
	}

	return &spacing, nil
}

func (r *spacingRepository) Update(e *entities.Spacing) (*entities.Spacing, error) {
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": e.ID}, bson.M{"$set": e})
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *spacingRepository) Delete(id string) error {
	spacingID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(context.Background(), bson.M{"_id": spacingID})
	if err != nil {
		return err
	}
	return nil
}

func (r *spacingRepository) DeleteAll() error {
	_, err := r.Collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}
