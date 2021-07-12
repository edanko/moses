package entities

type Bevel struct {
	AngleTs float64 `bson:"angle_ts,omitempty" json:"angle_ts,omitempty"`
	AngleOs float64 `bson:"angle_os,omitempty" json:"angle_os,omitempty"`
	DepthTs float64 `bson:"depth_ts,omitempty" json:"depth_ts,omitempty"`
	DepthOs float64 `bson:"depth_os,omitempty" json:"depth_os,omitempty"`
}
