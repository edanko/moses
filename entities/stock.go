package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Dim       string             `bson:"dim" json:"dim"`
	Quality   string             `bson:"quality" json:"quality"`
	Length    float64            `bson:"length" json:"length"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (s *Stock) ToBar() *Bar {
	b := &Bar{}

	b.Dim = s.Dim
	b.Quality = s.Quality
	b.Length = s.Length

	return b
}

func (s *Stock) Validate() error {
	if s.Dim == "" {
		return ErrDim
	}
	if s.Quality == "" {
		return ErrQuality
	}
	if s.Length <= 0 {
		return ErrLength
	}

	return nil
}
