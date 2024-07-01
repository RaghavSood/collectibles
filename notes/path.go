package notes

import (
	"fmt"
)

type NoteType string

const (
	Creator NoteType = "creator"
	Series  NoteType = "series"
	Item    NoteType = "item"
)

func DeriveNotePath(noteType NoteType, elements ...string) (string, error) {
	var noteID string
	if len(elements) != 1 {
		return "", fmt.Errorf("invalid number of elements: %d", len(elements))
	}

	filePath := fmt.Sprintf("%s/%s.md", noteType, elements[0])
	return filePath, nil
}
