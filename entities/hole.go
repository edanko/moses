package entities

type Hole struct {
	Name   string             `bson:"name,omitempty" json:"name,omitempty"`
	Params map[string]float64 `bson:"params,omitempty" json:"params,omitempty"`
}
