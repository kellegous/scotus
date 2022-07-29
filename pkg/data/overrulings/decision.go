package overrulings

type Decision struct {
	*Case
	Overruled []*Case
}
