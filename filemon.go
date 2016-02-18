package filemon

import (
	"gopkg.in/fsnotify.v1"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

type Watcher struct {
	watcher    *fsnotify.Watcher
	triggermap map[string]*triggered
	lock       sync.Mutex
	cb         OnFileEvCb
}

const (
	g_cooldown = (500 * time.Millisecond)
)

type EvType int

const (
	C_Create EvType = iota
	C_Modify
	C_Delete
	C_Rename
	C_Attrib
)

type OnFileEvCb func(*WatchEvent)

type WatchEvent struct {
	Fpath string
	Type  EvType
}

func NewWatcher(cb OnFileEvCb) *Watcher {
	w := new(Watcher)

	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		w = nil
	} else {
		w.triggermap = map[string]*triggered{}
		w.cb = cb
		go w.startWatching()
	}

	return w
}

// See Watch Example.
func (this *Watcher) Watch(fpath string) {
	if nil == this.watcher {
		return
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	fpath = filepath.Clean(fpath)

	err := this.watcher.Watch(fpath)
	if os.IsNotExist(err) {
		_, err = os.Create(fpath)
		if nil == err {
			// try to watch again for the last time
			err = this.watcher.Watch(fpath)
		}
	}
	if nil != err {
		fmt.Fprintln(os.Stderr, "Watch() err:", err)
		return
	}

}

func (this *Watcher) RemoveWatch(fpath string) {
	if nil == this.watcher {
		return
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.watcher.RemoveWatch(fpath); nil != err {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (this Watcher) WaitForKill() {
	onkill := make(chan os.Signal, 1)
	signal.Notify(onkill, os.Interrupt, os.Kill)
	<-onkill // wait for event
	fmt.Fprintln(os.Stderr, "\nkill triggered. exiting...")
}

func (this *Watcher) Close() {
	this.watcher.Close()
	this.watcher = nil
}
