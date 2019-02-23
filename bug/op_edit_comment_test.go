package bug

import (
	"testing"
	"time"

	"github.com/MichaelMure/git-bug/identity"
	"gotest.tools/assert"
)

func TestEdit(t *testing.T) {
	snapshot := Snapshot{}

	var rene = identity.NewIdentity("René Descartes", "rene@descartes.fr")

	unix := time.Now().Unix()

	create := NewCreateOp(rene, unix, "title", "create", nil)
	create.Apply(&snapshot)

	hash1, err := create.Hash()
	if err != nil {
		t.Fatal(err)
	}

	comment := NewAddCommentOp(rene, unix, "comment", nil)
	comment.Apply(&snapshot)

	hash2, err := comment.Hash()
	if err != nil {
		t.Fatal(err)
	}

	edit := NewEditCommentOp(rene, unix, hash1, "create edited", nil)
	edit.Apply(&snapshot)

	assert.Equal(t, len(snapshot.Timeline), 2)
	assert.Equal(t, len(snapshot.Timeline[0].(*CreateTimelineItem).History), 2)
	assert.Equal(t, len(snapshot.Timeline[1].(*AddCommentTimelineItem).History), 1)
	assert.Equal(t, snapshot.Comments[0].Message, "create edited")
	assert.Equal(t, snapshot.Comments[1].Message, "comment")

	edit2 := NewEditCommentOp(rene, unix, hash2, "comment edited", nil)
	edit2.Apply(&snapshot)

	assert.Equal(t, len(snapshot.Timeline), 2)
	assert.Equal(t, len(snapshot.Timeline[0].(*CreateTimelineItem).History), 2)
	assert.Equal(t, len(snapshot.Timeline[1].(*AddCommentTimelineItem).History), 2)
	assert.Equal(t, snapshot.Comments[0].Message, "create edited")
	assert.Equal(t, snapshot.Comments[1].Message, "comment edited")
}
