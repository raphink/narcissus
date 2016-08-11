package narcissus

import (
	"errors"

	"honnef.co/go/augeas"
)

// ErrNodeNotFound is returned if requested node was not found
var ErrNodeNotFound = errors.New("Node not found")

// Narcissus is a Narcissus handler
type Narcissus struct {
	Augeas *augeas.Augeas
}

// New returns a new Narcissus handler with given Augeas handler
func New(aug *augeas.Augeas) *Narcissus {
	return &Narcissus{
		Augeas: aug,
	}
}
