package diff

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	fText "github.com/zodimo/go-compose/compose/foundation/text"
	iconbutton "github.com/zodimo/go-compose/compose/material3/iconbutton"
	m3Text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/compose/ui"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type ChangeType int

const (
	ChangeEqual ChangeType = iota
	ChangeDiff
)

type DiffSection struct {
	Type     ChangeType
	Original string
	Modified string
}

func CalculateDiff(original, modified string) []DiffSection {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(original, modified)
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, c)

	var sections []DiffSection
	for i := 0; i < len(diffs); {
		if diffs[i].Type == diffmatchpatch.DiffEqual {
			sections = append(sections, DiffSection{
				Type:     ChangeEqual,
				Original: diffs[i].Text,
				Modified: diffs[i].Text,
			})
			i++
		} else {
			// Combine adjacent deletes and inserts into a single section
			var orig, mod strings.Builder
			for i < len(diffs) && diffs[i].Type != diffmatchpatch.DiffEqual {
				if diffs[i].Type == diffmatchpatch.DiffDelete {
					orig.WriteString(diffs[i].Text)
				} else if diffs[i].Type == diffmatchpatch.DiffInsert {
					mod.WriteString(diffs[i].Text)
				}
				i++
			}
			sections = append(sections, DiffSection{
				Type:     ChangeDiff,
				Original: orig.String(),
				Modified: mod.String(),
			})
		}
	}
	return sections
}

// DiffViewer renders the diff between original and modified strings
func DiffViewer(original, modified string, onApplyLeftToRight, onApplyRightToLeft func(sectionIndex int, section DiffSection)) api.Composable {
	sections := CalculateDiff(original, modified)

	return func(c api.Composer) api.Composer {
		c.StartBlock("DiffViewer")
		defer c.EndBlock()

		column.Column(
			c.Sequence(
				c.Range(len(sections), func(i int) api.Composable {
					return row.Row(
						c.Sequence(
							// Left Side
							row.Row(
								c.Sequence(renderSection(i, sections[i], true, onApplyLeftToRight)),
								row.WithModifier(weight.Weight(1)),
							),
							// Right Side
							row.Row(
								c.Sequence(renderSection(i, sections[i], false, onApplyRightToLeft)),
								row.WithModifier(weight.Weight(1)),
							),
						),
					)
				}),
			),
		)(c)

		return c
	}
}

func renderSection(index int, section DiffSection, isLeft bool, onApply func(int, DiffSection)) api.Composable {
	return func(c api.Composer) api.Composer {
		c.StartBlock("DiffSection")
		defer c.EndBlock()

		return column.Column(
			c.Sequence(func(c api.Composer) api.Composer {
				textStr := section.Original
				if !isLeft {
					textStr = section.Modified
				}

				if textStr == "" && section.Type == ChangeDiff {
					// Placeholder for deleted/empty blocks
					m3Text.Text("---", fText.WithModifier(
						background.Background(graphics.Color(0xFFC8C8C8)),
					))(c)
					return c
				}

				lines := strings.Split(strings.TrimSuffix(textStr, "\n"), "\n")
				var mod ui.Modifier = ui.EmptyModifier

				if section.Type == ChangeDiff {
					if isLeft {
						mod = background.Background(graphics.Color(0xFFFFC8C8)) // Faint red
					} else {
						mod = background.Background(graphics.Color(0xFFC8FFC8)) // Faint green
					}
				}

				c.Range(len(lines), func(j int) api.Composable {
					return m3Text.Text(lines[j], fText.WithModifier(mod))
				})(c)

				// Inline button for changes
				if section.Type == ChangeDiff {
					if isLeft {
						iconbutton.Standard(func() { onApply(index, section) }, icons.HardwareKeyboardArrowRight, "Apply Left to Right")(c)
					} else {
						iconbutton.Standard(func() { onApply(index, section) }, icons.HardwareKeyboardArrowLeft, "Apply Right to Left")(c)
					}
				}

				return c
			}),
		)(c)
	}
}
