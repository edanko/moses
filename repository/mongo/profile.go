package mongo

import (
	"context"
	"time"

	"github.com/edanko/moses/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileRepository interface {
	GetOne(id string) (*entities.Profile, error)
	GetAll() ([]*entities.Profile, error)
	Get(project, dimension, quality string) ([]*entities.Profile, error)

	Create(e *entities.Profile) (*entities.Profile, error)
	Update(e *entities.Profile) (*entities.Profile, error)
	Delete(id string) error
	DeleteAll() error
}

type profileRepository struct {
	Collection *mongo.Collection
}

func NewProfileRepo(collection *mongo.Collection) ProfileRepository {
	return &profileRepository{
		Collection: collection,
	}
}

func (p *profileRepository) GetOne(id string) (*entities.Profile, error) {
	profileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor := p.Collection.FindOne(context.Background(), bson.M{"_id": profileID})

	var profile entities.Profile
	err = cursor.Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (p *profileRepository) Create(e *entities.Profile) (*entities.Profile, error) {
	e.ID = primitive.NewObjectID()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()

	_, err := p.Collection.InsertOne(context.Background(), e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (p *profileRepository) GetAll() ([]*entities.Profile, error) {
	var profiles []*entities.Profile

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"name": 1})

	cursor, err := p.Collection.Find(context.Background(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var user entities.Profile
		_ = cursor.Decode(&user)
		profiles = append(profiles, &user)
	}
	return profiles, nil
}

func (p *profileRepository) Get(project, dim, quality string) ([]*entities.Profile, error) {
	var profiles []*entities.Profile

	cursor, err := p.Collection.Find(context.Background(), bson.M{"project": project, "dim": dim, "quality": quality})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var profile entities.Profile
		_ = cursor.Decode(&profile)

		profiles = append(profiles, &profile)
	}
	return profiles, nil
}

func (p *profileRepository) Update(e *entities.Profile) (*entities.Profile, error) {
	e.UpdatedAt = time.Now()
	_, err := p.Collection.UpdateOne(context.Background(), bson.M{"_id": e.ID}, bson.M{"$set": e})
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (p *profileRepository) Delete(id string) error {
	profileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = p.Collection.DeleteOne(context.Background(), bson.M{"_id": profileID})
	if err != nil {
		return err
	}
	return nil
}

func (p *profileRepository) DeleteAll() error {
	_, err := p.Collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}
