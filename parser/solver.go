package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type named interface {
	name() string
	pos() lexer.Position
}

type collision struct {
	t         string
	name      string
	locations []lexer.Position
}

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

type collisionsErr struct {
	collisions []collision
}

func (c collisionsErr) Error() string {
	out := make([]string, len(c.collisions))

	for i, col := range c.collisions {
		out[i] = col.String()
	}

	return strings.Join(out, "\n\n")
}

func namedSlice(t []annotatedType, e []Enum) (out []named) {
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

func merge(in []*AST) (out *AST, err error) {
	names := make(map[string][]lexer.Position)
	intermediateOut := new(AST)

	for _, a := range in {
		if a == nil {
			continue
		}

		intermediateOut.Types = append(intermediateOut.Types, a.Types...)
		intermediateOut.Enums = append(intermediateOut.Enums, a.Enums...)

		for _, t := range namedSlice(a.Types, a.Enums) {
			if _, ok := names[t.name()]; !ok {
				names[t.name()] = make([]lexer.Position, 0)
			}

			names[t.name()] = append(names[t.name()], t.pos())
		}
	}

	err = toCollisionError("type", names)
	if err != nil {
		return
	}

	return intermediateOut, nil
}

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
		return collisionsErr{collisions}
	}

	return nil
}
