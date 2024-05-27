package class

func GetNotesText(notes []*Comment) string {
	s := ""
	for _, n := range notes {
		s += n.Text
	}
	return s
}
