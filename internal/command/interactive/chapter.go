package interactive

type chapter string

const (
	ChapterMain          = chapter("Main Menu")
	ChapterUnimplemented = chapter("Under Construction")
	ChapterFinal         = chapter("Quitting")
)

func (c *chapter) isFinal() bool {
	return *c == ChapterFinal
}

func (c *chapter) change(nextChap chapter) bool {
	// TODO: should add some state machine here
	*c = nextChap
	return true
}
