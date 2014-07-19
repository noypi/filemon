package filemon_test

func ExampleWatcher() {

	// create a new watcher
	w := filemon.NewWatcher(func(ev *filemon.WatchEvent) {
		fmt.Println(ev)
	})

	// watch the current path
	w.Watch("./")

	// wait for a ctrl+c
	w.WaitForKill()

}
