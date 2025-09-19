package interactive

type chapter string

const (
	ChapterMain              = chapter("Main Menu")
	ChapterUnimplemented     = chapter("Under Construction")
	ChapterPasswordGen       = chapter("Choose pathword strength")
	ChapterShowGeneratedPass = chapter("Generated password:")
	ChapterFinal             = chapter("Quitting")
)

// ifNil returns true if chapter is not set.
func (c *chapter) isNil() bool {
	return c == nil || *c == ""
}

// isFinal returns true if this chapter should lead to the app quitting.
func (c *chapter) isFinal() bool {
	return *c == ChapterFinal
}

// change switches the chapter to the nextChap,
// returns true if the transfer is legal and transition was made.
func (c *chapter) change(nextChap chapter) bool {
	// TODO: should add some state machine here to check if the change is legal
	*c = nextChap
	return true
}
