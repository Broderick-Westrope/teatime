package tui

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestOverlayCenter(t *testing.T) {
	tt := map[string]struct {
		bg      string
		overlay string
		want    string
	}{
		"simple": {
			bg: heredoc.Doc(`
				Facere enim neque consectetur soluta tenetur ducimus omnis. Voluptatibus accusantium maiores quia eaque velit nesciunt hic saepe tenetur.
				Amet quidem reprehenderit ex. Error illum sit est expedita sapiente neque. Laborum vero necessitatibus similique suscipit nam.
				Tempore occaecati eligendi accusamus eos similique harum impedit. Quas nam molestiae architecto quam.
				Accusamus pariatur facilis ea nostrum exercitationem quam. Sit ipsam aperiam aspernatur hic fugit officia inventore.
				Reiciendis doloribus ut eius id. Repellendus eum enim. Reprehenderit veritatis nulla molestiae nulla veniam.
				Nemo animi nisi blanditiis. Eligendi tempora laudantium assumenda nam.
			`),
			overlay: "*********\n*****",
			want: heredoc.Doc(`
				Facere enim neque consectetur soluta tenetur ducimus omnis. Voluptatibus accusantium maiores quia eaque velit nesciunt hic saepe tenetur.
				Amet quidem reprehenderit ex. Error illum sit est expedita sapiente neque. Laborum vero necessitatibus similique suscipit nam.
				Tempore occaecati eligendi accusamus eos similique harum impedit*********m molestiae architecto quam.
				Accusamus pariatur facilis ea nostrum exercitationem quam. Sit i*****aperiam aspernatur hic fugit officia inventore.
				Reiciendis doloribus ut eius id. Repellendus eum enim. Reprehenderit veritatis nulla molestiae nulla veniam.
				Nemo animi nisi blanditiis. Eligendi tempora laudantium assumenda nam.
			`),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := OverlayCenter(tc.bg, tc.overlay)
			assert.Equal(t, tc.want, result)
		})
	}
}

func TestOverlay(t *testing.T) {
	tt := map[string]struct {
		bg      string
		overlay string
		row     int
		col     int
		want    string
	}{
		"single line; start": {
			bg:      "Nostrum libero modi velit neque dolores.",
			overlay: "*********",
			row:     0,
			col:     0,
			want:    "*********ibero modi velit neque dolores.",
		},
		"single line; middle": {
			bg:      "Nostrum libero modi velit neque dolores.",
			overlay: "*********",
			row:     0,
			col:     10,
			want:    "Nostrum li********* velit neque dolores.",
		},
		"single line; end beyond background": {
			bg:      "Nostrum libero modi velit neque dolores.",
			overlay: "*********",
			row:     0,
			col:     35,
			want:    "Nostrum libero modi velit neque dol*********",
		},
		"single line; lipgloss styled": {
			bg: "Nostrum libero modi velit neque dolores.",
			overlay: lipgloss.NewStyle().PaddingLeft(2).
				Underline(true).Foreground(lipgloss.Color("1")).Render("*****"),
			row:  0,
			col:  5,
			want: "Nostr  *****ro modi velit neque dolores.",
		},
		"single line; manual escape code": {
			bg:      "Nostrum libero modi velit neque dolores.",
			overlay: "\x1b[31m*****\x1b[0m",
			row:     0,
			col:     5,
			want:    "Nostr\u001B[31m*****\u001B[0mbero modi velit neque dolores.",
		},
		"multi-line background; overlay middle line": {
			bg:      "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
			overlay: "*****",
			row:     2,
			col:     0,
			want:    "Line 1\nLine 2\n*****3\nLine 4\nLine 5",
		},
		"multi-line overlay; beyond background": {
			bg:      "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
			overlay: "*******\n*******",
			row:     1,
			col:     5,
			want:    "Line 1\nLine *******\nLine *******\nLine 4\nLine 5",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := Overlay(tc.bg, tc.overlay, tc.row, tc.col)
			assert.Equal(t, tc.want, result)
		})
	}
}
