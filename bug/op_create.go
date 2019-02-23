package bug

import (
	"fmt"
	"strings"

	"github.com/MichaelMure/git-bug/identity"

	"github.com/MichaelMure/git-bug/util/git"
	"github.com/MichaelMure/git-bug/util/text"
)

var _ Operation = &CreateOperation{}

// CreateOperation define the initial creation of a bug
type CreateOperation struct {
	OpBase
	Title   string     `json:"title"`
	Message string     `json:"message"`
	Files   []git.Hash `json:"files"`
}

func (op *CreateOperation) base() *OpBase {
	return &op.OpBase
}

func (op *CreateOperation) Hash() (git.Hash, error) {
	return hashOperation(op)
}

func (op *CreateOperation) Apply(snapshot *Snapshot) {
	snapshot.Title = op.Title

	comment := Comment{
		Message:  op.Message,
		Author:   op.Author,
		UnixTime: Timestamp(op.UnixTime),
	}

	snapshot.Comments = []Comment{comment}
	snapshot.Author = op.Author
	snapshot.CreatedAt = op.Time()

	hash, err := op.Hash()
	if err != nil {
		// Should never error unless a programming error happened
		// (covered in OpBase.Validate())
		panic(err)
	}

	snapshot.Timeline = []TimelineItem{
		&CreateTimelineItem{
			CommentTimelineItem: NewCommentTimelineItem(hash, comment),
		},
	}
}

func (op *CreateOperation) GetFiles() []git.Hash {
	return op.Files
}

func (op *CreateOperation) Validate() error {
	if err := opBaseValidate(op, CreateOp); err != nil {
		return err
	}

	if text.Empty(op.Title) {
		return fmt.Errorf("title is empty")
	}

	if strings.Contains(op.Title, "\n") {
		return fmt.Errorf("title should be a single line")
	}

	if !text.Safe(op.Title) {
		return fmt.Errorf("title is not fully printable")
	}

	if !text.Safe(op.Message) {
		return fmt.Errorf("message is not fully printable")
	}

	return nil
}

// Sign post method for gqlgen
func (op *CreateOperation) IsAuthored() {}

func NewCreateOp(author identity.Interface, unixTime int64, title, message string, files []git.Hash) *CreateOperation {
	return &CreateOperation{
		OpBase:  newOpBase(CreateOp, author, unixTime),
		Title:   title,
		Message: message,
		Files:   files,
	}
}

// CreateTimelineItem replace a Create operation in the Timeline and hold its edition history
type CreateTimelineItem struct {
	CommentTimelineItem
}

// Convenience function to apply the operation
func Create(author identity.Interface, unixTime int64, title, message string) (*Bug, *CreateOperation, error) {
	return CreateWithFiles(author, unixTime, title, message, nil)
}

func CreateWithFiles(author identity.Interface, unixTime int64, title, message string, files []git.Hash) (*Bug, *CreateOperation, error) {
	newBug := NewBug()
	createOp := NewCreateOp(author, unixTime, title, message, files)

	if err := createOp.Validate(); err != nil {
		return nil, createOp, err
	}

	newBug.Append(createOp)

	return newBug, createOp, nil
}
