package main

import (
	"log"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/layout/tree"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3/divider"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			func(c api.Composer) api.Composer {
				// Title
				m3text.HeadlineMedium("Tree Demo", ftext.WithModifier(padding.All(16)))(c)

				divider.Divider()(c)

				// Two panes: Declarative and Data-Driven
				row.Row(
					func(c api.Composer) api.Composer {
						// Pane 1: Declarative
						column.Column(
							func(c api.Composer) api.Composer {
								m3text.TitleMedium("Declarative", ftext.WithModifier(padding.All(8)))(c)
								c.Key("declarative", DeclarativeTreeDemo())(c)
								return c
							},
							column.WithModifier(weight.Weight(1).Then(size.FillMaxHeight())),
						)(c)

						// Divider
						divider.Divider(
							divider.WithModifier(size.Width(1).Then(size.FillMaxHeight())),
						)(c)

						// Pane 2: Data-Driven
						column.Column(
							func(c api.Composer) api.Composer {
								m3text.TitleMedium("Data-Driven", ftext.WithModifier(padding.All(8)))(c)
								c.Key("datadriven", DataDrivenTreeDemo())(c)
								return c
							},
							column.WithModifier(weight.Weight(1).Then(size.FillMaxHeight())),
						)(c)

						return c
					},
					row.WithModifier(size.FillMax()),
				)(c)

				return c
			},
			column.WithModifier(size.FillMax().Then(background.Background(graphics.ColorWhite))),
		)(c)
	}
}

func DeclarativeTreeDemo() api.Composable {
	return func(c api.Composer) api.Composer {
		state := tree.RememberTreeState(c)

		return tree.Tree(
			state,
			func(scope tree.TreeScope) {
				scope.Branch("dir1", m3text.BodyMedium("Directory 1"), func(s tree.TreeScope) {
					s.Node("file1", m3text.BodyMedium("File 1.txt"))
					s.Node("file2", m3text.BodyMedium("File 2.txt"))
					s.Branch("subdir1", m3text.BodyMedium("Sub-directory 1"),
						func(s2 tree.TreeScope) {
							s2.Node("file3", m3text.BodyMedium("File 3.txt"))
						},
					)
				})
				scope.Branch("dir2", m3text.BodyMedium("Directory 2"), func(s tree.TreeScope) {
					s.Node("file4", m3text.BodyMedium("File 4.txt"))
				})
				scope.Node("rootFile", m3text.BodyMedium("Root File.txt"))
			},
			tree.WithOnBranchOpened(func(id any) {
				log.Printf("[Declarative] Branch opened: %v", id)
			}),
			tree.WithOnBranchClosed(func(id any) {
				log.Printf("[Declarative] Branch closed: %v", id)
			}),
			tree.WithOnSelected(func(id any) {
				log.Printf("[Declarative] Node selected: %v", id)
			}),
			tree.WithOnUnselected(func(id any) {
				log.Printf("[Declarative] Node unselected: %v", id)
			}),
		)(c)
	}
}

// Data-Driven Mock Data
type FileSystemItem struct {
	Name     string
	IsFolder bool
	Children []string // IDs of children
}

var fsData = map[string]FileSystemItem{
	"root":   {Name: "/", IsFolder: true, Children: []string{"etc", "home", "usr", "readme"}},
	"etc":    {Name: "etc", IsFolder: true, Children: []string{"hosts", "passwd"}},
	"home":   {Name: "home", IsFolder: true, Children: []string{"user"}},
	"usr":    {Name: "usr", IsFolder: true, Children: []string{"bin", "lib"}},
	"user":   {Name: "jaco", IsFolder: true, Children: []string{"docs", "pics"}},
	"docs":   {Name: "Documents", IsFolder: true, Children: []string{}},
	"pics":   {Name: "Pictures", IsFolder: true, Children: []string{}},
	"bin":    {Name: "bin", IsFolder: true, Children: []string{}},
	"lib":    {Name: "lib", IsFolder: true, Children: []string{}},
	"readme": {Name: "README.md", IsFolder: false},
	"hosts":  {Name: "hosts", IsFolder: false},
	"passwd": {Name: "passwd", IsFolder: false},
}

func DataDrivenTreeDemo() api.Composable {
	return func(c api.Composer) api.Composer {
		state := tree.RememberTreeState(c)

		return tree.TreeFromData(
			state,
			[]any{"root"}, // Root IDs
			func(id any) []any {
				strID := id.(string)
				item := fsData[strID]
				if !item.IsFolder {
					return nil
				}
				children := make([]any, len(item.Children))
				for i, child := range item.Children {
					children[i] = child
				}
				return children
			},
			func(id any) bool {
				strID := id.(string)
				return fsData[strID].IsFolder
			},
			func(id any) api.Composable {
				strID := id.(string)
				name := fsData[strID].Name
				// Simple icon logic
				icon := "📄"
				if fsData[strID].IsFolder {
					icon = "📁"
				}
				return row.Row(
					func(c api.Composer) api.Composer {
						m3text.BodyMedium(icon)(c)
						spacer.Width(8)(c)
						m3text.BodyMedium(name)(c)
						return c
					},
				)
			},
			tree.WithOnBranchOpened(func(id any) {
				log.Printf("[DataDriven] Branch opened: %v", id)
			}),
			tree.WithOnBranchClosed(func(id any) {
				log.Printf("[DataDriven] Branch closed: %v", id)
			}),
			tree.WithOnSelected(func(id any) {
				log.Printf("[DataDriven] Node selected: %v", id)
			}),
			tree.WithOnUnselected(func(id any) {
				log.Printf("[DataDriven] Node unselected: %v", id)
			}),
		)(c)
	}
}
