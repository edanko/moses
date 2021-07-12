package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Spacing struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Machine    string             `bson:"machine" json:"machine"`
	Dim        string             `bson:"dim" json:"dim"`
	Length     float64            `bson:"length" json:"length"`
	Name       string             `bson:"name" json:"name"`
	HasBevel   bool               `bson:"has_bevel" json:"has_bevel"`
	HasScallop bool               `bson:"has_scallop" json:"has_scallop"`
}

func (s *Spacing) Validate() error {
	if s.Machine == "" {
		return ErrMachine
	}
	if s.Dim == "" {
		return ErrDim
	}
	if s.Length <= 0 {
		return ErrLength
	}

	return nil
}
