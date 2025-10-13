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
				if len(note) > 30 {
					note = note[:30] + "..."
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
