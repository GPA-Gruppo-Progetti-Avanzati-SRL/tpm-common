package fmtstructure

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type MDWriter struct {
	numberOfLines int
	writer        io.Writer
}

func NewMDWriter(writer io.Writer, title string) *MDWriter {
	mdw := &MDWriter{writer: writer}

	mdw.write(fmt.Sprintf("## Readme %s\n\n", title))
	mdw.write(fmt.Sprintf("|%s|%s|%s|%s|", " * ", " Path ", " Tipo ", " Note "))
	mdw.write(fmt.Sprintf("|%s|%s|%s|%s|", "---", "------", "------", "------"))
	return mdw
}

func (mdw *MDWriter) WritePathValue(aPath string, aKind reflect.Kind, aType string, aValue string, isEmpty bool, aTagInfo TagInfo) {

	mdw.numberOfLines++

	if !(aTagInfo.OmitEmpty && isEmpty) {
		pathDepth := strings.Count(aPath, ".")

		if aKind == reflect.Array || aKind == reflect.Slice || aKind == reflect.Map || aKind == reflect.Struct {
			mdw.write(fmt.Sprintf("|%s|%s|%s|%s|", mdw.boldValue(fmt.Sprint(pathDepth)), mdw.boldValue(aPath), mdw.boldValue(aType), mdw.boldValue(aValue)))
		} else {
			mdw.write(fmt.Sprintf("|%d|%s|%s|%s|", pathDepth, aPath, aType, aValue))
		}
	}
}

func (mdw *MDWriter) boldValue(s string) string {
	if s == "" {
		return s
	}

	sb := strings.Builder{}
	sb.WriteString("**")
	sb.WriteString(s)
	sb.WriteString("**")

	return sb.String()
}

func (mdw *MDWriter) isStdOut() bool {
	return mdw.writer == nil
}

func (mdw *MDWriter) write(s string) error {

	if mdw.isStdOut() {
		fmt.Println(s)
	}

	io.WriteString(mdw.writer, s)
	io.WriteString(mdw.writer, "\n")
	return nil
}
