package util

import (
	"regexp"
	"strings"
)

// utils used by generators to handle references to packages and folders.

var HierarchicalNamePathForm = regexp.MustCompile("^(./)?[a-z\\-A-Z_0-9]+(/[a-z\\-A-Z_0-9]+)*$")
var HierarchicalNameDotForm = regexp.MustCompile("^[a-z\\-A-Z_0-9]+(\\.[a-z\\-A-Z_0-9]+)*$")

type PackageName string

func (pkgName PackageName) IsPathForm() bool {
	return HierarchicalNamePathForm.Match([]byte(pkgName))
}

func (pkgName PackageName) IsDottedForm() bool {
	return HierarchicalNameDotForm.Match([]byte(pkgName))
}

func (pkgName PackageName) Slice() []string {

	if pkgName.IsPathForm() {
		return strings.Split(string(pkgName), "/")
	}

	if pkgName.IsDottedForm() {
		return strings.Split(string(pkgName), ".")
	}

	return nil
}

func NameWellFormed(n string) bool {
	return FieldNameWellFormed(n)
}

func FieldNameWellFormed(n string) bool {

	if n == "" {
		return false
	}
	return true
}
