package cli

import "github.com/helmedeiros/amp/internal/music"

// style wraps a string, typically in ANSI color codes.
type style func(string) string

func identity(s string) string { return s }

// Theme controls how status fields are styled. Keeping styling behind a theme
// lets the renderer stay pure and lets the command pick plain vs colored output
// based on the environment.
type Theme struct {
	label style
	title style
	state func(music.PlayerState, string) string
}

// PlainTheme renders without any styling.
var PlainTheme = Theme{
	label: identity,
	title: identity,
	state: func(_ music.PlayerState, s string) string { return s },
}

// ansi returns a style that wraps text in the given SGR code and a reset.
func ansi(code string) style {
	return func(s string) string { return "\x1b[" + code + "m" + s + "\x1b[0m" }
}

// ColorTheme styles output for a terminal: the state word colored by state
// (green playing, yellow paused, grey stopped), dim labels, and a bold title.
func ColorTheme() Theme {
	dim, bold := ansi("2"), ansi("1")
	green, yellow, grey := ansi("32"), ansi("33"), ansi("90")

	return Theme{
		label: dim,
		title: bold,
		state: func(st music.PlayerState, s string) string {
			switch st {
			case music.Playing:
				return green(s)
			case music.Paused:
				return yellow(s)
			default:
				return grey(s)
			}
		},
	}
}
