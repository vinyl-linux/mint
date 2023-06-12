package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

// named is an interface used when solving unique 'things', such
// as types, enums, and values- this interface allows us to do
// clever things, like list all of the places where a unique name
// has been used in many places
type named interface {
	name() string
	pos() lexer.Position
}

// collision stores the number of places a particular `named` has
// been called in order to provide decent error messages
type collision struct {
	t         string
	name      string
	locations []lexer.Position
}

// String describes the specific collision in a way that can be used
// by errors and warnings to provide a decent message
func (c collision) String() string {
	cols := make([]string, len(c.locations))
	for i, l := range c.locations {
		cols[i] = l.String()
	}

	return fmt.Sprintf("%s %q has been defined in %d places: %s",
		c.t,
		c.name,
		len(c.locations),
		strings.Join(cols, "\n\t"),
	)
}

// collisionsErr holds a slice of collisions in order to implement the
// error.Error interface
type collisionsErr []collision

// Error returns an error message, thus fulfilling the error.Error
// interface
func (c collisionsErr) Error() string {
	out := make([]string, len(c))

	for i, col := range c {
		out[i] = col.String()
	}

	return strings.Join(out, "\n\n")
}

// namedSlice coerces a slice of types and enums into a slice of
// named types to aid the ast solver
func namedSlice(t []AnnotatedType, e []Enum) (out []named) {
	out = make([]named, len(t)+len(e))
	idx := 0

	for _, elem := range t {
		out[idx] = elem
		idx++
	}

	for _, elem := range e {
		out[idx] = elem
		idx++
	}

	return
}

// toCollisionError wraps creation of a collisionError by accepting a map
// of would-be unique names to the number of definitions for that name.
//
// Any time it finds multiple definitions, it appends to the underlying slice
// powering a collisionError. It then returns all of the collisions, providing
// multiple errors at once
func toCollisionError(t string, names map[string][]lexer.Position) error {
	collisions := make([]collision, 0)
	for name, locs := range names {
		if len(locs) > 1 {
			collisions = append(collisions, collision{
				t:         t,
				name:      name,
				locations: locs,
			})
		}
	}

	if len(collisions) != 0 {
		return collisionsErr(collisions)
	}

	return nil
}
