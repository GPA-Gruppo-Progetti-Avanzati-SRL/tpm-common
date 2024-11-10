package fileutil

import (
	"os"
	"path/filepath"
	"strings"
)

func ListPathHierarchy(startingPath string, upWard bool) []string {

	if startingPath == "/" {
		return []string{startingPath}
	}

	if strings.HasSuffix(startingPath, "/") {
		startingPath = startingPath[0 : len(startingPath)-1]
	}

	if startingPath == "" || startingPath[0:1] == "." || startingPath[0:1] != "/" {
		wd, _ := os.Getwd()
		startingPath = filepath.Join(wd, startingPath)
	}

	pathSegs := strings.Split(startingPath, "/")
	if pathSegs[0] == "" {
		pathSegs[0] = "/"
	}

	resultPaths := make([]string, 0, len(pathSegs))
	startLoop := len(pathSegs)
	stepLoop := -1
	if !upWard {
		startLoop = 1
		stepLoop = 1
	}
	for i := startLoop; (upWard && i > 0) || (!upWard && i <= len(pathSegs)); i += stepLoop {
		p := filepath.Join(pathSegs[0:i]...)
		resultPaths = append(resultPaths, p)
	}
	return resultPaths
}

func FindFileInClosestDirectory(startingFolder, fileName string) string {
	ph := ListPathHierarchy(startingFolder, true)

	// Fix. Loop was going out of range in ph.
	for i := 0; i < len(ph); i++ {
		fp := filepath.Join(ph[i], fileName)
		_, err := os.Stat(fp)
		if err == nil {
			return fp
		}
	}

	return ""
}

func FindGoModFolder(startingFolder string) string {
	ph := FindFileInClosestDirectory(startingFolder, "go.mod")

	if ph != "" {
		return filepath.Dir(ph)
	}

	return ""
}

func ResolvePath(p string) (string, bool) {

	var okGoMod bool
	var okUserHome bool
	if p != "" && p[:1] == "~" {
		goModPath := FindGoModFolder(".")
		userHomePath, err := os.UserHomeDir()
		if err != nil {
			return "", false
		}

		if userHomePath != "" {
			userHomePath, okUserHome = resolvePathInHomeFolder(userHomePath, p[1:])
		}

		if goModPath != "" {
			goModPath, okGoMod = resolvePathInHomeFolder(goModPath, p[1:])
		}

		if okUserHome && !okGoMod {
			return userHomePath, okUserHome
		}

		return goModPath, okGoMod
	}

	return p, FileExists(p)
}

func resolvePathInHomeFolder(home string, p string) (string, bool) {
	exists := false
	retp := filepath.Join(home, p)
	if FileExists(retp) {
		exists = true
	}

	return retp, exists
}
