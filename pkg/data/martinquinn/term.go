package martinquinn

type Term struct {
	Year     int        `json:"year"`
	Justices []*Justice `json:"justices"`
}
