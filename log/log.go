package log

import (
	"fmt"
	"runtime"
	"time"
)


func Enter(l int, v ...interface{}) {

	_, f, line, ok := runtime.Caller(1)

	if !ok {
		fmt.Printf("%s | Level: %d | On _ of _ | %+v\n",
			time.Now().Format(time.RFC1123),
			v,
		)
		return
	}

	fmt.Printf("%s | Level: %d | On %d of %s | %+v\n",
		time.Now().Format(time.RFC1123),
		l,
		line, f,
		v,
	)
}