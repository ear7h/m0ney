package log

import (
	"fmt"
	"runtime"
	"time"
	"io"
	"os"
)

const (
	OK = "Okay"
	WARNING = "Warning"
	ERROR = "Error"
	FATAL = "Fatal"
)

var OUT io.Writer

func init() {
	OUT = os.Stdout
}

func Enter(level string, v ...interface{}) {

	_, f, line, ok := runtime.Caller(1)


	if !ok {
		fmt.Fprintf(OUT, "%s\t| Level: %d | On _ of _ | %s\n",
			time.Now().Format(time.RFC1123),
			level,
			fmt.Sprint(v...),
		)
		return
	}

	fmt.Fprintf(OUT, "%s\t| Level: %d | On %d of %s | %+v\n",
		time.Now().Format(time.RFC1123),
		level,
		line, f,
		fmt.Sprint(v...),
	)
}