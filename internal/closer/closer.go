package closer

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	SignalSet = []os.Signal{
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGABRT,
	}
)

var (
	// ExitCodeOK is a successfull exit code.
	ExitCodeOK = 0
	// ExitCodeErr is a failure exit code.
	ExitCodeErr = 1
	// ExitSignals is the active list of signals to watch for.
	ExitSignals = SignalSet
)

type closer struct {
	signalChanel chan os.Signal
	exitChanel   chan int
	cleanups     []func()
	mutex        sync.Mutex
	closeOnce    sync.Once
}

var c = newCloser()

func newCloser() *closer {
	c := &closer{
		signalChanel: make(chan os.Signal, 1),
		exitChanel:   make(chan int),
	}

	signal.Notify(c.signalChanel, ExitSignals...)

	// start waiting
	go c.wait()
	return c
}

func (c *closer) wait() {
	select {
	case <-c.signalChanel:
		for _, fn := range c.cleanups {
			fn()
		}
		c.exitChanel <- ExitCodeOK
	}
}

func (c *closer) closeErr() {
	c.closeOnce.Do(func() {
		c.exitChanel <- ExitCodeErr
	})
}

func Bind(cleanup func()) {
	c.mutex.Lock()
	// store in the reverse order
	s := make([]func(), 0, 1+len(c.cleanups))
	s = append(s, cleanup)
	c.cleanups = append(s, c.cleanups...)
	c.mutex.Unlock()
}

func Checked(target func() error) int {
	go func() {
		defer func() {
			// check if there was a panic
			if x := recover(); x != nil {
				c.closeErr()
			}
		}()
		if err := target(); err != nil {
			// close with an error
			c.closeErr()
		}
	}()

	return <-c.exitChanel
}
