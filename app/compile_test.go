package main

import (
	"testing"
	"os"
	"github.com/stretchr/testify/require"
	"github.com/cappuccinotm/flangc/app/cmd"
	"path"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	files, err := os.ReadDir("../_example")
	require.NoError(t, err)

	for _, dirEntry := range files {
		fInfo, err := dirEntry.Info()
		require.NoError(t, err, "Failed to get file info for %s", dirEntry.Name())

		if fInfo.IsDir() {
			continue
		}

		err = (&cmd.Run{FileLocation: path.Join("../_example", dirEntry.Name())}).
			Execute([]string{})
		assert.NoError(t, err, "Failed to compile %s", dirEntry.Name())
	}
}
