package util

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

	for i := 0; i <= len(ph); i++ {
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
