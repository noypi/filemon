package main

import (
	"fmt"
	"github.com/noypi/filemon"
)

func main() {
	w := filemon.NewWatcher(func(ev *filemon.WatchEvent) {
		fmt.Println(ev)
	})
	w.Watch("./")

	w.WaitForKill()
}
