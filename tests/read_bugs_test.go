package tests

import (
	"testing"

	"github.com/MichaelMure/git-bug/bug"
	"github.com/MichaelMure/git-bug/util/test"
)

func TestReadBugs(t *testing.T) {
	repo := test.CreateFilledRepo(15)
	bugs := bug.ReadAllLocalBugs(repo)
	for b := range bugs {
		if b.Err != nil {
			t.Fatal(b.Err)
		}
	}
}

func benchmarkReadBugs(bugNumber int, t *testing.B) {
	repo := test.CreateFilledRepo(bugNumber)
	t.ResetTimer()

	for n := 0; n < t.N; n++ {
		bugs := bug.ReadAllLocalBugs(repo)
		for b := range bugs {
			if b.Err != nil {
				t.Fatal(b.Err)
			}
		}
	}
}

func BenchmarkReadBugs5(b *testing.B)   { benchmarkReadBugs(5, b) }
func BenchmarkReadBugs25(b *testing.B)  { benchmarkReadBugs(25, b) }
func BenchmarkReadBugs150(b *testing.B) { benchmarkReadBugs(150, b) }
