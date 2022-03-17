package main

import (
	"testing"
	"os"
	"github.com/stretchr/testify/require"
	"github.com/cappuccinotm/flangc/app/cmd"
	"path"
	"github.com/stretchr/testify/assert"
	"log"
	"io"
	"bytes"
)

func TestAll(t *testing.T) {
	files, err := os.ReadDir("../_example")
	require.NoError(t, err)

	for _, dirEntry := range files {
		dirEntry := dirEntry
		fInfo, err := dirEntry.Info()
		require.NoError(t, err, "Failed to get file info for %s", dirEntry.Name())
		if fInfo.IsDir() {
			continue
		}
		t.Run(dirEntry.Name(), func(t *testing.T) {
			log.Printf("[INFO] parsing file %s", dirEntry.Name())

			fLoc := path.Join("../_example", dirEntry.Name())

			f, err := os.Open(fLoc)
			require.NoError(t, err)

			b, err := io.ReadAll(f)
			require.NoError(t, err)

			if bytes.HasPrefix(b, []byte("//skip")) {
				t.Skip()
			}

			require.NoError(t, f.Close())

			err = (&cmd.Run{FileLocation: fLoc, FailOnError: true}).
				Execute([]string{})
			assert.NoError(t, err, "Failed to compile %s", dirEntry.Name())
		})
	}
}
