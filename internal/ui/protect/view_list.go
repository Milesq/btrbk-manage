package protect

import (
	"strings"

	"milesq.dev/btrbk-manage/internal/utils"
)

func (m Model) viewList(b *strings.Builder) {
	if len(m.backups) == 0 {
		b.WriteString("No snapshot groups found.\n")
		return
	}

	for i, g := range m.backups {
		if i == m.cursor {
			b.WriteString(focusedStyle.Render("> "))
		} else {
			b.WriteString("  ")
		}

		if g.IsProtected {
			b.WriteString(focusedStyle.Render("★ "))
		} else if g.IsTrashed {
			b.WriteString(trashStyle.Render("🗑 "))
		} else {
			b.WriteString("  ")
		}

		timestampStyle := emptyStyle
		if g.IsProtected {
			timestampStyle = unpersistedStyle
		}

		b.WriteString(timestampStyle.Render(utils.PrettifyDate(g.Timestamp)))

		if g.ProtectionNote.Note != "" || g.ProtectionNote.Reason != "" || len(g.ProtectionNote.Tags) > 0 {
			b.WriteString(" ")
			b.WriteString(blurredStyle.Render("—"))
			if g.ProtectionNote.Note != "" {
				note := g.ProtectionNote.Note
				maxNoteLen := m.maxNoteLen()
				if len(note) > maxNoteLen {
					note = note[:maxNoteLen] + "..."
				}
				b.WriteString(" ")
				b.WriteString(blurredStyle.Render(note))
			}
			if g.ProtectionNote.Reason != "" {
				b.WriteString(" ")
				b.WriteString(blurredStyle.Render("[" + g.ProtectionNote.Reason + "]"))
			}
			if len(g.ProtectionNote.Tags) > 0 {
				b.WriteString(" ")
				b.WriteString(blurredStyle.Render(strings.Join(g.ProtectionNote.Tags, ", ")))
			}
		}

		b.WriteRune('\n')
	}
}

func (m Model) maxNoteLen() int {
	// Fixed width: cursor(2) + icon(2) + date(16) + separator(3) + buffer(10)
	const fixedWidth = 33
	available := m.width - fixedWidth
	return utils.MinMax(10, available, 100)
}
