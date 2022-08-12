package web

import (
	"github.com/kellegous/scotus/pkg/build"
	"github.com/kellegous/scotus/pkg/data"
)

type Data struct {
	Build *build.Info
	Model *data.Model
}
