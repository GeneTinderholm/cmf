package this

import (
	"os"
	"path"
	"runtime"
)

// Dir returns the path *with* the path separator at the end
func Dir() string {
	_, f, _, _ := runtime.Caller(1)
	return path.Dir(f) + string(os.PathSeparator)
}

func File() string {
	_, f, _, _ := runtime.Caller(1)
	return f
}
