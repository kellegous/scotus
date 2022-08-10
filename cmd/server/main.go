package main

import (
	"context"
	"flag"

	"github.com/kellegous/scotus/pkg/logging"
	"github.com/kellegous/scotus/pkg/web"
	"go.uber.org/zap"
)

type Flags struct {
	DataDir   string
	ResetData bool
	HTTPAddr  string
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.DataDir,
		"data-dir",
		"data",
		"the directory where the data is stashed")

	fs.BoolVar(
		&f.ResetData,
		"reset-data",
		false,
		"whether to nuke the data directory")

	fs.StringVar(
		&f.HTTPAddr,
		"http.addr",
		":8080",
		"the address where the server will run")
}

func main() {
	var flags Flags
	flags.Register(flag.CommandLine)
	flag.Parse()

	lg := logging.MustSetup()

	ctx := context.Background()

	lg.Info("server has started",
		zap.String("http.addr", flags.HTTPAddr),
		zap.Bool("reset-data", flags.ResetData),
		zap.String("data-dir", flags.DataDir))

	if err := web.ListenAndServe(ctx, flags.HTTPAddr); err != nil {
		lg.Fatal("unable to run http server",
			zap.Error(err))
		return
	}
}
