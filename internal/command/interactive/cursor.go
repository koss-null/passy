package interactive

type cursor uint16

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

func (c *cursor) setStart() {
	var zero cursor
	c = &zero
}
