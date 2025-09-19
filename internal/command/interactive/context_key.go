package interactive

import (
	"fmt"
	"strings"
)

type contextKey string

const (
	// don't forget to add new value to allContextKeys
	ContextKeyPass             = contextKey("Pass")
	ContextKeyLastPassStrength = contextKey("LastPassStrength")
)

var allContextKeys = []contextKey{
	ContextKeyPass,
	ContextKeyLastPassStrength,
}

// Substitute takes target string, finds ck.Template occurences
// and substitute all of them with the subst value.
func (ck contextKey) Substitute(target, subst string) string {
	return strings.ReplaceAll(target, ck.Template(), subst)
}

func (ck contextKey) Template() string {
	return fmt.Sprintf("{{$%s}}", string(ck))
}
