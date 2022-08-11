package build

import (
	"debug/buildinfo"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/kellegous/buildname"
)

type Info struct {
	Version     string      `json:"version"`
	Name        string      `json:"name"`
	Time        time.Time   `json:"time"`
	Environment Environment `json:"go"`
	Deps        []*Module   `json:"deps"`
}

func Read() (*Info, error) {
	bi, err := buildinfo.ReadFile(os.Args[0])
	if err != nil {
		return nil, err
	}

	settings := map[string]string{}
	for _, setting := range bi.Settings {
		settings[setting.Key] = setting.Value
	}

	version := settings["vcs.revision"]

	t, err := time.Parse(time.RFC3339, settings["vcs.time"])
	if err != nil {
		return nil, fmt.Errorf("vcs.time: %w", err)
	}

	return &Info{
		Version: version,
		Name:    buildname.FromVersion(version),
		Time:    t,
		Environment: Environment{
			Version: bi.GoVersion,
			Arch:    settings["GOARCH"],
			OS:      settings["GOOS"],
		},
		Deps: toModules(bi),
	}, nil
}

func toModules(bi *debug.BuildInfo) []*Module {
	mods := make([]*Module, 0, len(bi.Deps))
	for _, dep := range bi.Deps {
		mods = append(mods, &Module{
			Name:    dep.Path,
			Version: dep.Version,
		})
	}
	return mods
}
