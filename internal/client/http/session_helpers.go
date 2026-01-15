package http

import (
	"os"
	"path/filepath"
)

func tokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".postflow", "token")
}

func saveToken(t string) {
	p := tokenPath()
	_ = os.MkdirAll(filepath.Dir(p), 0700)
	_ = os.WriteFile(p, []byte(t), 0600)
}

func loadToken() string {
	b, _ := os.ReadFile(tokenPath())
	return string(b)
}
