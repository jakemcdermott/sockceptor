package main

import (
	"testing"
)

func TestCmdline(t *testing.T) {
	t.Run("Missing required node-id param results in error.", func(t *testing.T) {
		args := []string{"--noe-id", "node1", "--tcp-listener", "2222",
			"--tun-proxy", "service=tunnel", "remotenode=node3", "remoteservice=tunnel"}

		initCmdline()
		result := runCmdline(args)

		if result == nil {
			t.Errorf("Expected error but returned ok")
		}

		expected := "Unknown command: --noe-id\n"
		actual := result.Error()

		if actual != expected {
			t.Errorf("Expected %s but returned %s", expected, actual)
		}
	})
}
