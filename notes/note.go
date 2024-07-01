package notes

type Note struct {
	NoteID       string   `json:"note_id"`
	Type         NoteType `json:"type"`
	Data         string   `json:"data"`
	Path         string   `json:"path"`
	PathElements []string `json:"path_elements"`
}
