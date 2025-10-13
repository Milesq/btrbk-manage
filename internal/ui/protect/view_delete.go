package protect

import (
	"strings"

	"milesq.dev/btrbk-manage/internal/utils"
)

func (m Model) viewDeleteConfirmation(b *strings.Builder) {
	backup := m.selected

	b.WriteString(warningStyle.Render("⚠ DELETE BACKUP"))
	b.WriteString("\n\n")

	b.WriteString("Timestamp: ")
	b.WriteString(utils.PrettifyDate(backup.Timestamp))
	b.WriteString("\n")

	statusText := "unprotected"
	if backup.IsProtected {
		statusText = "protected"
	} else if backup.IsTrashed {
		statusText = "trashed"
	}
	b.WriteString("Status: ")
	b.WriteString(statusText)
	b.WriteString("\n")

	if backup.ProtectionNote.Note != "" || backup.ProtectionNote.Reason != "" || len(backup.ProtectionNote.Tags) > 0 {
		b.WriteString("\n")
		if backup.ProtectionNote.Note != "" {
			b.WriteString("Note: ")
			b.WriteString(backup.ProtectionNote.Note)
			b.WriteString("\n")
		}
		if backup.ProtectionNote.Reason != "" {
			b.WriteString("Reason: ")
			b.WriteString(backup.ProtectionNote.Reason)
			b.WriteString("\n")
		}
		if len(backup.ProtectionNote.Tags) > 0 {
			b.WriteString("Tags: ")
			b.WriteString(strings.Join(backup.ProtectionNote.Tags, ", "))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	b.WriteString(warningStyle.Render("This action cannot be undone!"))
	b.WriteString("\n\n")

	b.WriteString("Delete this backup? ")
	b.WriteString(blurredStyle.Render("Y / N\n"))
}
