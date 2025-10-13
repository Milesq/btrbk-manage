package protect

import (
	"fmt"
	"io"
)

func (m Model) writeHelpMessage(w io.Writer) {
	dot := focusedStyle.Render(" • ")

	fmt.Fprint(
		w,
		"\nSpace ",
		blurredStyle.Render("to un/protect backup"),
		dot,

		"Enter ",
		blurredStyle.Render("to edit note"),
		dot,

		"d ",
		blurredStyle.Render("to delete"),
		dot,

		"q ",
		blurredStyle.Render("to quit"),
	)

	fmt.Fprint(w, "\n\n")

	tStyle := emptyStyle
	if m.trashMode {
		tStyle = activeFilterStyle
	}

	mStyle := emptyStyle
	if m.listProtectedOnly {
		mStyle = activeFilterStyle
	}

	fmt.Fprint(
		w,
		tStyle.Render("t"),
		blurredStyle.Render(" to see trash"),
		dot,

		mStyle.Render("m"),
		blurredStyle.Render(" to list manual only"),
	)

	if m.trashMode {
		fmt.Fprint(
			w,
			"\n\n",
			"D ",
			blurredStyle.Render("to remove all trashed backups"),
		)
	}
}
