package entities

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrDim        = errors.New("dimension can't be empty")
	ErrQuality    = errors.New("quality can't be empty")
	ErrLength     = errors.New("length must be greater than zero")
	ErrFullLength = errors.New("full length must be greater than zero")
	ErrUsedLength = errors.New("used length can't be greater than bar length")
	ErrNestName   = errors.New("nest name can't be empty")
	ErrProject    = errors.New("project can't be empty")
	ErrFrom       = errors.New("remnant can't be without parent")
	ErrSection    = errors.New("section can't be empty")
	ErrPosNo      = errors.New("posno can't be empty")
	ErrMachine    = errors.New("machine can't be empty")
	ErrProfiles   = errors.New("nest must contain at least one profile")
	ErrLaunch     = errors.New("nest must contain launch")
	ErrQuantity   = errors.New("quantity can't be zero")
	ErrEnd        = errors.New("profile must contain endcut")
	ErrParts      = errors.New("can't nest without parts")
	ErrBars       = errors.New("can't nest without bars")
)
