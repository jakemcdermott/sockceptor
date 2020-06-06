package main

import (
	"fmt"
	_ "github.com/ghjm/sockceptor/pkg/backends"
	"github.com/ghjm/sockceptor/pkg/cmdline"
	_ "github.com/ghjm/sockceptor/pkg/controlsock"
	"github.com/ghjm/sockceptor/pkg/debug"
	"github.com/ghjm/sockceptor/pkg/netceptor"
	_ "github.com/ghjm/sockceptor/pkg/services"
	_ "github.com/ghjm/sockceptor/pkg/workceptor"
	"os"
	"time"
)

var nodeID string

type nodeIDCfg struct {
	NodeID string `description:"Node ID" barevalue:"yes" required:"yes"`
}

func (cfg nodeIDCfg) Prepare() error {
	nodeID = cfg.NodeID
	netceptor.MainInstance = netceptor.New(cfg.NodeID)
	return nil
}

type debugCfg struct{}

func (cfg debugCfg) Prepare() error {
	debug.Enable = true
	return nil
}

type traceCfg struct{}

func (cfg traceCfg) Prepare() error {
	debug.Trace = true
	return nil
}

type nullBackendCfg struct{}

func (cfg nullBackendCfg) Run() error {
	// This is a null backend that doesn't do anything
	netceptor.AddBackend()
	return nil
}

func initCmdline() {
	cmdline.AddConfigType("node-id", "Network node ID of this instance", nodeIDCfg{}, true, nil)
	cmdline.AddConfigType("debug", "Enables debug output", debugCfg{}, false, nil)
	cmdline.AddConfigType("trace", "Enables packet tracing output", traceCfg{}, false, nil)
	cmdline.AddConfigType("local-only", "Run a self-contained node with no backends", nullBackendCfg{}, false, nil)
}

func runCmdline(args []string) error {
	return cmdline.ParseAndRun(args)
}

func main() {
	initCmdline()
	err := runCmdline(os.Args[1:])

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if nodeID == "" {
		println("Must specify a node ID")
		os.Exit(1)
	}

	// Fancy footwork to set an error exitcode if we're immediately exiting at startup
	done := make(chan struct{})
	go func() {
		netceptor.BackendWait()
		close(done)
	}()
	select {
	case <-done:
		if netceptor.BackendCount() > 0 {
			fmt.Printf("All backends have failed. Exiting.\n")
			os.Exit(1)
		} else {
			fmt.Printf("Nothing to do - no backends were specified.\n")
			fmt.Printf("Run %s --help for command line instructions.\n", os.Args[0])
			os.Exit(1)
		}
	case <-time.After(100 * time.Millisecond):
	}
	debug.Printf("Initialization complete\n")
	<-done
}
