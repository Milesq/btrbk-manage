package protect

import (
	"fmt"
	"io"
)

func writeHelpMessage(w io.Writer) {
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
		blurredStyle.Render("to quit\n"),
	)
}
