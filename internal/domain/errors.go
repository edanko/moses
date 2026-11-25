package domain

import "errors"

var (
	ErrProject            = errors.New("project is required")
	ErrSection            = errors.New("section is required")
	ErrPositionNumber     = errors.New("position number is required")
	ErrProfileDesignation = errors.New("profile designation is invalid")
	ErrGrade              = errors.New("steel grade is required")
	ErrLength             = errors.New("length must be greater than zero")
	ErrQuantity           = errors.New("quantity must be greater than zero")
	ErrEndCut             = errors.New("both end cuts are required")
	ErrRemnantSource      = errors.New("remnant source is required")
	ErrRemnantConsumed    = errors.New("remnant is already consumed")
	ErrStockType          = errors.New("unknown stock type")
	ErrNoParts            = errors.New("nest has no parts")
)
