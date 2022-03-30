package ghostscript

import (
	"errors"
	"fmt"
	"os"
)

const (
	GhostScriptWorkFolderEnvVarName = "GHOSTSCRIPT_WORKFOLDER"
	DefaultGhostScriptWorkFolder    = "/tmp/ghs"

	WhichGhostScriptCommandEnvVarName = "GHOSTSCRIPT_CMD"
	DefaultGhostscriptCommand         = "/opt/homebrew/bin/gs"
)

type GsConfig struct {
	workFolderMountPoint string
	whichGs              string
}

type GsOption func(cfg *GsConfig)

func WithWorkFolder(p string) GsOption {
	return func(cfg *GsConfig) {
		cfg.workFolderMountPoint = p
	}
}

func WithWhichGs(e string) GsOption {
	return func(cfg *GsConfig) {
		cfg.whichGs = e
	}
}

type Ghostscripter struct {
	cfg GsConfig
}

func (gs *Ghostscripter) WorkFolder() string {
	return gs.cfg.workFolderMountPoint
}

func NewGhostscripter(opts ...GsOption) (Ghostscripter, error) {

	cfg := GsConfig{}
	for _, o := range opts {
		o(&cfg)
	}

	gs := Ghostscripter{cfg: cfg}

	if cfg.workFolderMountPoint == "" {
		envValue := os.Getenv(GhostScriptWorkFolderEnvVarName)
		if envValue != "" {
			gs.cfg.workFolderMountPoint = envValue
		} else {
			gs.cfg.workFolderMountPoint = DefaultGhostScriptWorkFolder
		}
	}

	if err := ghostscriptExists(gs.cfg.workFolderMountPoint, true, false); err != nil {
		return gs, err
	}

	if cfg.whichGs == "" {
		envValue := os.Getenv(WhichGhostScriptCommandEnvVarName)
		if envValue != "" {
			gs.cfg.whichGs = envValue
		} else {
			gs.cfg.whichGs = DefaultGhostscriptCommand
		}
	}

	if err := ghostscriptExists(gs.cfg.whichGs, false, true); err != nil {
		return gs, err
	}

	return gs, nil
}

func ghostscriptExists(v string, dir bool, exe bool) error {

	var err error
	var fi os.FileInfo
	if fi, err = os.Stat(v); err == nil {
		if dir && !fi.IsDir() {
			return fmt.Errorf("path %s exists but is not a directory", v)
		}

		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("path %s doesn't exist", v)
	}

	return err
}
