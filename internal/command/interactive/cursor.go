package interactive

type cursor int16

func (c *cursor) up() {
	if *c > 0 {
		*c--
	}
}

func (c *cursor) down(threshold int) {
	if *c < cursor(threshold) {
		*c++
	}
}

func (c *cursor) setStart(opts []option) {
	var zero cursor
	*c = zero
	for i := range opts[:len(opts)-1] {
		// skip all unselectables
		if opts[i].optType == OptTypeUnselectableString {
			*c = cursor(i + 1)
			continue
		}
		return
	}
}
