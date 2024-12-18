package notes

import (
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/RaghavSood/collectibles/clogger"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var log = clogger.NewLogger("notes")

func RenderNote(content string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(content))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return wrapForTailwind(string(markdown.Render(doc, renderer)))
}

func ReadNoteFile(file fs.File) (string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ReadNotes(pointers []NotePointer) []Note {
	var notes []Note
	for _, pointer := range pointers {
		path, err := DeriveNotePath(pointer.NoteType, pointer.PathElements...)
		if err != nil {
			log.Debug().
				Err(err).
				Str("noteType", string(pointer.NoteType)).
				Strs("pathElements", pointer.PathElements).
				Msg("Failed to derive note path")

			continue
		}

		noteFile, err := NotesFS.Open(path)
		if err != nil {
			log.Debug().
				Err(err).
				Str("path", path).
				Msg("Failed to open note file")
			continue
		}

		content, err := ReadNoteFile(noteFile)
		if err != nil {
			log.Debug().
				Err(err).
				Str("path", path).
				Msg("Failed to read note file")
			continue
		}

		log.Info().
			Str("path", path).
			Str("noteType", string(pointer.NoteType)).
			Msg("Read note")

		renderedContent := RenderNote(content)

		note := Note{
			NoteID:       fmt.Sprintf("%s/%s", pointer.NoteType, strings.Join(pointer.PathElements, ":")),
			Type:         pointer.NoteType,
			Data:         renderedContent,
			Path:         path,
			PathElements: pointer.PathElements,
		}

		notes = append(notes, note)
	}

	return notes
}
