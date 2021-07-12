package entities

import (
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Nest struct {
	ID          primitive.ObjectID   `bson:"_id" json:"_id"`
	Project     string               `bson:"project" json:"project"`
	Launch      string               `bson:"launch" json:"launch"`
	Name        string               `bson:"name" json:"name"`
	Bar         *Bar                 `bson:"bar" json:"bar"`
	Profiles    []*Profile           `bson:"-" json:"-"`
	ProfilesIds []primitive.ObjectID `bson:"profiles_ids" json:"profiles_ids"`
	Machine     string               `bson:"machine" json:"machine"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type NestSlice []*Nest

func (bs NestSlice) Len() int { return len(bs) }
func (bs NestSlice) Less(i, j int) bool {
	return bs[i].Bar.Length < bs[j].Bar.Length
}
func (bs NestSlice) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func (n *Nest) GetRemnant() *Remnant {
	r := &Remnant{}

	r.Project = n.Project
	r.Dim = n.Bar.Dim
	r.Quality = n.Bar.Quality
	r.Length = n.Bar.Length - n.Bar.UsedLength
	r.Marking = n.Name + "R01"
	r.From = n.Name

	minRemnant, _ := strconv.ParseFloat(os.Getenv("MIN_REMNANT"), 64)
	if r.Length < minRemnant {
		return nil
	}

	return r
}

func (n *Nest) PartsLen() (res float64) {
	for _, p := range n.Profiles {
		res += p.Length
	}
	return
}

func (n *Nest) Scrap() (res float64) {
	for _, p := range n.Profiles {
		res += p.FullLength - p.Length
	}
	if r := n.GetRemnant(); r == nil {
		res += n.Bar.Length - n.Bar.UsedLength
	}
	return
}

func (n *Nest) Validate() error {
	if n.Name == "" {
		return ErrNestName
	}
	if n.Project == "" {
		return ErrProject
	}
	if n.Machine == "" {
		return ErrMachine
	}
	if len(n.Profiles) == 0 {
		return ErrProfiles
	}
	if n.Launch == "" {
		return ErrLaunch
	}

	return nil
}
