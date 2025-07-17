package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompassFileToBase64(t *testing.T) {

	entries, err := os.ReadDir("D:\\EEG\\sub03-pangrichong-250709\\cam")
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		if v.IsDir() {
			p := filepath.Join("D:\\EEG\\sub03-pangrichong-250709\\cam", v.Name())
			println(p)
		}
	}

}
