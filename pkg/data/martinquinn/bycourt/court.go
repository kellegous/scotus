package bycourt

type Court struct {
	Year  int
	Stats []*Stats
}

func (c *Court) Median() float64 {
	// just take the first median because none of them differ
	// all that much when there are replacements in a term
	return c.Stats[0].MedianJusticeScore
}
