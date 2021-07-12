package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Project    string             `bson:"project" json:"project"`
	Section    string             `bson:"section" json:"section"`
	PosNo      string             `bson:"pos_no" json:"pos_no"`
	Quality    string             `bson:"quality" json:"quality"`
	Dim        string             `bson:"dim" json:"dim"`
	Length     float64            `bson:"length" json:"length"`
	FullLength float64            `bson:"-" json:"-"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	L          *End               `bson:"l" json:"l"`
	R          *End               `bson:"r" json:"r"`
	TraceBevel *Bevel             `bson:"trace_bevel,omitempty" json:"trace_bevel,omitempty"`
	Holes      []*Hole            `bson:"holes,omitempty" json:"holes,omitempty"`
	Source     string             `bson:"source,omitempty" json:"source,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	//IsBended  bool      `bson:"is_bended,omitempty" json:"is_bended,omitempty"`
	//BendingCurve [][]float64 `bson:"bending_curve,omitempty" json:"bending_curve,omitempty"`
}

type ProfileSlice []*Profile

func (ps ProfileSlice) Len() int { return len(ps) }
func (ps ProfileSlice) Less(i, j int) bool {
	return ps[i].Length > ps[j].Length
}
func (ps ProfileSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (p *Profile) Validate() error {
	if p.Project == "" {
		return ErrProject
	}
	if p.Section == "" {
		return ErrSection
	}
	if p.PosNo == "" {
		return ErrPosNo
	}
	if p.Quality == "" {
		return ErrQuality
	}
	if p.Dim == "" {
		return ErrDim
	}
	if p.Length <= 0 {
		return ErrLength
	}
	if p.Quantity <= 0 {
		return ErrQuantity
	}
	if p.L == nil || p.R == nil {
		return ErrEnd
	}

	return nil
}

func (p *Profile) InvertHolesX() {
	for _, h := range p.Holes {
		h.Params["X"] = p.Length - h.Params["X"]
	}
}
