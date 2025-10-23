package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	public := "public"
	port := 8080

	fmt.Printf("Serving %s at http://localhost:%d\n", public, port)

	fs := http.FileServer(http.Dir(public))
	http.Handle("/", fs)

	// OPTIONAL: auto rebuild every time a file changes
	if os.Getenv("WATCH") == "1" {
		go watchAndRebuild()
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// optional — simple periodic rebuild (basic version)
func watchAndRebuild() {
	last := time.Now()
	for {
		time.Sleep(2 * time.Second)
		latest := newestMod("content")
		if latest.After(last) {
			fmt.Println("Detected changes → rebuilding…")
			cmd := exec.Command("go", "run", "./cmd/build")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			_ = cmd.Run()
			last = latest
		}
	}
}

func newestMod(dir string) time.Time {
	var newest time.Time
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
		return nil
	})
	return newest
}
