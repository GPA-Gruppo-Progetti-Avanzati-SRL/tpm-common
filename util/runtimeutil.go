package util

import "runtime"

func GetExecutingFunctionInfo() (string, int, string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.File, frame.Line, frame.Function
}
