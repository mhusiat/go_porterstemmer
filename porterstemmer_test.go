package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/mhusiat/go_porterstemmer/porterstemmer"
)

func TestTheWholeApp(t *testing.T) {
	vocFd, err := os.Open("voc.txt")
	if err != nil {
		t.Fatalf("Cannot open voc: %v", err)
	}
	defer vocFd.Close()
	voc := bufio.NewScanner(vocFd)

	outputFd, err := os.Open("output.txt")
	if err != nil {
		t.Fatalf("Cannot open output: %v", err)
	}
	defer outputFd.Close()
	out := bufio.NewScanner(outputFd)

	var miss, total int

	for lineNo := 1; ; lineNo++ {
		if !voc.Scan() || !out.Scan() {
			// Quit on any output end.
			t.Logf("Miss=%d, Total=%d, Perc=%.2f%%", miss, total, float64(miss)/float64(total)*100)
			return
		}
		word := voc.Text()
		result := porterstemmer.Stem(word)
		want := out.Text()

		if result != want {
			miss++
			t.Errorf("Line %3d %20s: want %s, but got %s.", lineNo, word, want, result)
		}
		total++
	}
}
