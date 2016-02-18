package filemon

import (
	"github.com/howeyc/fsnotify"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type triggered struct {
	cb         OnFileEvCb
	fpath      string
	fname      string
	lock       sync.Mutex
	busy       bool
	lastupdate time.Time
}

func (this *Watcher) startWatching() {
	for {
		select {
		case ev := <-this.watcher.Event:
			this.trigger(ev)
		case ev := <-this.watcher.Error:
			fmt.Fprintln(os.Stderr, ev.Error())
		}
	}
}

func (this *Watcher) trigger(ev *fsnotify.FileEvent) {
	this.lock.Lock()
	defer this.lock.Unlock()

	var v *triggered
	var ok bool
	if v, ok = this.triggermap[ev.Name]; !ok {
		v = &triggered{
			fpath: ev.Name,
			fname: filepath.Base(ev.Name),
			busy:  false,
		}
		this.triggermap[ev.Name] = v
	}
	go this.handlecb(v, ev)
}
func (this *triggered) canrun() (bret bool) {

	this.lock.Lock()
	defer this.lock.Unlock()
	if this.busy || time.Now().Sub(this.lastupdate) < g_cooldown {
		return
	} else {
		this.busy = true
		bret = true
	}

	return
}

func (this *triggered) setLastUpdate() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.lastupdate = time.Now()
	this.busy = false
}

func (this *Watcher) handlecb(v *triggered, ev *fsnotify.FileEvent) {

	if !v.canrun() {
		return
	}
	defer v.setLastUpdate()

	// in windows some events are sent 3 times
	c := time.Tick(g_cooldown)
	<-c

	var t EvType
	// execute
	if ev.IsModify() {
		t = C_Modify
	} else if ev.IsCreate() {
		t = C_Create
	} else if ev.IsDelete() {
		t = C_Delete
	} else if ev.IsRename() {
		t = C_Rename
	} else if ev.IsAttrib() {
		t = C_Attrib
	} else {
		fmt.Fprintln(os.Stderr, "unknown event")
		return
	}

	this.cb(&WatchEvent{
		Type:  t,
		Fpath: ev.Name,
	})

}
