package random_bugs

import (
	"math/rand"
	"strings"
	"time"

	"github.com/MichaelMure/git-bug/bug"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/repository"
	"github.com/icrowley/fake"
)

type opsGenerator func(bug.Interface, identity.Interface)

type Options struct {
	BugNumber    int
	PersonNumber int
	MinOp        int
	MaxOp        int
}

func DefaultOptions() Options {
	return Options{
		BugNumber:    15,
		PersonNumber: 5,
		MinOp:        3,
		MaxOp:        20,
	}
}

func CommitRandomBugs(repo repository.ClockedRepo, opts Options) {
	CommitRandomBugsWithSeed(repo, opts, time.Now().UnixNano())
}

func CommitRandomBugsWithSeed(repo repository.ClockedRepo, opts Options, seed int64) {
	bugs := GenerateRandomBugsWithSeed(opts, seed)

	for _, b := range bugs {
		err := b.Commit(repo)
		if err != nil {
			panic(err)
		}
	}
}

func GenerateRandomBugs(opts Options) []*bug.Bug {
	return GenerateRandomBugsWithSeed(opts, time.Now().UnixNano())
}

func GenerateRandomBugsWithSeed(opts Options, seed int64) []*bug.Bug {
	rand.Seed(seed)
	fake.Seed(seed)

	opsGenerators := []opsGenerator{
		comment,
		comment,
		title,
		labels,
		open,
		close,
	}

	result := make([]*bug.Bug, opts.BugNumber)

	for i := 0; i < opts.BugNumber; i++ {
		addedLabels = []string{}

		b, _, err := bug.Create(
			randomPerson(opts.PersonNumber),
			time.Now().Unix(),
			fake.Sentence(),
			paragraphs(),
		)

		if err != nil {
			panic(err)
		}

		nOps := opts.MinOp

		if opts.MaxOp > opts.MinOp {
			nOps += rand.Intn(opts.MaxOp - opts.MinOp)
		}

		for j := 0; j < nOps; j++ {
			index := rand.Intn(len(opsGenerators))
			opsGenerators[index](b, randomPerson(opts.PersonNumber))
		}

		result[i] = b
	}

	return result
}

func GenerateRandomOperationPacks(packNumber int, opNumber int) []*bug.OperationPack {
	return GenerateRandomOperationPacksWithSeed(packNumber, opNumber, time.Now().UnixNano())
}

func GenerateRandomOperationPacksWithSeed(packNumber int, opNumber int, seed int64) []*bug.OperationPack {
	// Note: this is a bit crude, only generate a Create + Comments

	rand.Seed(seed)
	fake.Seed(seed)

	result := make([]*bug.OperationPack, packNumber)

	for i := 0; i < packNumber; i++ {
		opp := &bug.OperationPack{}

		var op bug.Operation

		op = bug.NewCreateOp(
			randomPerson(5),
			time.Now().Unix(),
			fake.Sentence(),
			paragraphs(),
			nil,
		)

		opp.Append(op)

		for j := 0; j < opNumber-1; j++ {
			op = bug.NewAddCommentOp(
				randomPerson(5),
				time.Now().Unix(),
				paragraphs(),
				nil,
			)
			opp.Append(op)
		}

		result[i] = opp
	}

	return result
}

func person() identity.Interface {
	return identity.NewBare(fake.FullName(), fake.EmailAddress())
}

var persons []identity.Interface

func randomPerson(personNumber int) identity.Interface {
	if len(persons) == 0 {
		persons = make([]identity.Interface, personNumber)
		for i := range persons {
			persons[i] = person()
		}
	}

	index := rand.Intn(personNumber)
	return persons[index]
}

func paragraphs() string {
	p := fake.Paragraphs()
	return strings.Replace(p, "\t", "\n\n", -1)
}

func comment(b bug.Interface, p identity.Interface) {
	_, _ = bug.AddComment(b, p, time.Now().Unix(), paragraphs())
}

func title(b bug.Interface, p identity.Interface) {
	_, _ = bug.SetTitle(b, p, time.Now().Unix(), fake.Sentence())
}

func open(b bug.Interface, p identity.Interface) {
	_, _ = bug.Open(b, p, time.Now().Unix())
}

func close(b bug.Interface, p identity.Interface) {
	_, _ = bug.Close(b, p, time.Now().Unix())
}

var addedLabels []string

func labels(b bug.Interface, p identity.Interface) {
	var removed []string
	nbRemoved := rand.Intn(3)
	for nbRemoved > 0 && len(addedLabels) > 0 {
		index := rand.Intn(len(addedLabels))
		removed = append(removed, addedLabels[index])
		addedLabels[index] = addedLabels[len(addedLabels)-1]
		addedLabels = addedLabels[:len(addedLabels)-1]
		nbRemoved--
	}

	var added []string
	nbAdded := rand.Intn(3)
	for i := 0; i < nbAdded; i++ {
		label := fake.Word()
		added = append(added, label)
		addedLabels = append(addedLabels, label)
	}

	// ignore error
	// if the randomisation produce no changes, no op
	// is added to the bug
	_, _, _ = bug.ChangeLabels(b, p, time.Now().Unix(), added, removed)
}
