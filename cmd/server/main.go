package main

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"os/signal"

	"github.com/kellegous/scotus/pkg/build"
	"github.com/kellegous/scotus/pkg/logging"
	"github.com/kellegous/scotus/pkg/web"

	"go.uber.org/zap"
)

type Flags struct {
	DataDir   string
	ResetData bool
	HTTP      struct {
		Addr      string
		AssetsDir string
	}
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
		&f.HTTP.Addr,
		"http.addr",
		":8080",
		"the address where the server will run")

	fs.StringVar(
		&f.HTTP.AssetsDir,
		"http.assets-dir",
		"",
		"where to load web assets from")
}

func startWebpackWatch(
	ctx context.Context,
	root string,
) error {
	c := exec.CommandContext(ctx, "npx", "webpack", "watch", "--mode=development")
	c.Dir = root
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}

func startHTTPServer(
	ctx context.Context,
	addr string,
	assetsDir string,
	data *web.Data,
) chan error {
	ch := make(chan error)

	go func() {
		ch <- web.ListenAndServe(ctx, addr, assetsDir, data)
	}()

	return ch
}

func main() {
	var flags Flags
	flags.Register(flag.CommandLine)
	flag.Parse()

	lg := logging.MustSetup()

	b, err := build.Read()
	if err != nil {
		lg.Fatal("unable to read build info",
			zap.Error(err))
	}

	ctx, done := signal.NotifyContext(
		context.Background(),
		os.Interrupt)
	defer done()

	lg.Info("server has started",
		zap.String("http.addr", flags.HTTP.Addr),
		zap.String("http.assets-dir", flags.HTTP.AssetsDir),
		zap.Bool("reset-data", flags.ResetData),
		zap.String("data-dir", flags.DataDir),
		zap.String("version", b.Version),
		zap.String("name", b.Name))

	if err := startWebpackWatch(ctx, "."); err != nil {
		lg.Fatal("could not start webpack watcher",
			zap.Error(err))
	}

	ch := startHTTPServer(
		ctx,
		flags.HTTP.Addr,
		flags.HTTP.AssetsDir,
		&web.Data{Build: b},
	)

	select {
	case err := <-ch:
		lg.Fatal("http server error",
			zap.Error(err))
	case <-ctx.Done():
		break
	}
}
