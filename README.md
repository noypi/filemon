[![Build Status](https://travis-ci.org/noypi/filemon.svg?branch=master)](https://travis-ci.org/noypi/filemon)
[![GoDoc](https://godoc.org/github.com/noypi/filemon?status.png)](http://godoc.org/github.com/noypi/filemon)

###Example:
```go
// create a new watcher
w := filemon.NewWatcher(func(ev *filemon.WatchEvent) {
    fmt.Println(ev)
})

// watch the current path
w.Watch("./")

// wait for a ctrl+c
w.WaitForKill()
```

###Structs
```go
type WatchEvent struct {
    Fpath string
    Type  EvType
}

type EvType int
const (
    C_Create EvType = iota
    C_Modify
    C_Delete
    C_Rename
    C_Attrib
)
```
