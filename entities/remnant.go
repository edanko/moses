package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Remnant struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Project   string             `bson:"project" json:"project"`
	From      string             `bson:"from" json:"from"`
	Dim       string             `bson:"dim" json:"dim"`
	Quality   string             `bson:"quality" json:"quality"`
	Length    float64            `bson:"length" json:"length"`
	Marking   string             `bson:"marking" json:"marking"`
	UsedIn    string             `bson:"used_in,omitempty" json:"used_in,omitempty"`
	Used      bool               `bson:"used" json:"used"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (r *Remnant) ToBar() *Bar {
	if r.Used {
		return nil
	}

	b := &Bar{}
	b.Dim = r.Dim
	b.Quality = r.Quality
	b.Length = r.Length
	b.IsRemnant = true
	b.RemnantID = r.ID

	return b
}

func (r *Remnant) Validate() error {
	if r.Project == "" {
		return ErrProject
	}
	if r.From == "" {
		return ErrFrom
	}
	if r.Dim == "" {
		return ErrDim
	}
	if r.Quality == "" {
		return ErrQuality
	}
	if r.Length <= 0 {
		return ErrLength
	}
	return nil
}
