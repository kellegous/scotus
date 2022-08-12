package data

import (
	"context"

	"github.com/kellegous/scotus/pkg/async"
	"github.com/kellegous/scotus/pkg/data/martinquinn/bycourt"
	"github.com/kellegous/scotus/pkg/data/option"
	"github.com/kellegous/scotus/pkg/data/overrulings"
	"github.com/kellegous/scotus/pkg/data/scotusdb"
)

type Model struct {
	Overrulings       []*overrulings.Decision
	MartinQuinnByYear []*bycourt.Court
	SCOTUSDBCases     []*scotusdb.Term
}

func LoadModel(
	ctx context.Context,
	dataDir string,
) (*Model, error) {
	fa := async.Run(
		func() ([]*overrulings.Decision, error) {
			return overrulings.Read(ctx, option.WithDataDir(dataDir))
		},
		nil)

	fb := async.Run(
		func() ([]*bycourt.Court, error) {
			return bycourt.Read(ctx, option.WithDataDir(dataDir))
		},
		nil)

	fc := async.Run(
		func() ([]*scotusdb.Term, error) {
			return scotusdb.Read(ctx, scotusdb.WithDataDir(dataDir))
		},
		nil)

	m := &Model{}
	var err error
	m.Overrulings, err = fa.Resolve()
	if err != nil {
		return nil, err
	}

	m.MartinQuinnByYear, err = fb.Resolve()
	if err != nil {
		return nil, err
	}

	m.SCOTUSDBCases, err = fc.Resolve()
	if err != nil {
		return nil, err
	}

	return m, nil
}
