package tree

import (
	"reflect"
	"testing"
)

func TestFoldIn(t *testing.T) {
	t.Run("single leaf tree", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		var visited []string
		result := FoldIn(tr, visited, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		want := []string{"leaf1", "root"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldIn order = %v, want %v", result, want)
		}
	})

	t.Run("single node tree", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {})
		})

		var visited []string
		result := FoldIn(tr, visited, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		want := []string{"node1", "root"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldIn order = %v, want %v", result, want)
		}
	})

	t.Run("multiple children - processes right to left", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
			scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
			scope.AddLeaf("leaf3", func(scope TreeLeafScope) {})
		})

		var visited []string
		result := FoldIn(tr, visited, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		want := []string{"leaf3", "leaf2", "leaf1", "root"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldIn order = %v, want %v", result, want)
		}
	})

	t.Run("nested nodes - depth first", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("level1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
					scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
				})
			})
		})

		var visited []string
		result := FoldIn(tr, visited, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		want := []string{"leaf2", "leaf1", "level1", "root"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldIn order = %v, want %v", result, want)
		}
	})

	t.Run("complex tree", func(t *testing.T) {
		// Structure:
		//       root
		//      /    \
		//    node1  node2
		//    /  \     \
		//  l1   l2    l3
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("l1", func(scope TreeLeafScope) {})
					scope.AddLeaf("l2", func(scope TreeLeafScope) {})
				})
			})
			scope.AddNode("node2", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("l3", func(scope TreeLeafScope) {})
				})
			})
		})

		var visited []string
		result := FoldIn(tr, visited, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		want := []string{"l3", "node2", "l2", "l1", "node1", "root"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldIn order = %v, want %v", result, want)
		}
	})

	t.Run("accumulates values correctly", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("a", func(scope TreeLeafScope) {})
			scope.AddLeaf("b", func(scope TreeLeafScope) {})
			scope.AddLeaf("c", func(scope TreeLeafScope) {})
		})

		result := FoldIn(tr, "", func(carry string, item Tree) string {
			return carry + item.Label()
		})

		want := "cbaroot"
		if result != want {
			t.Errorf("FoldIn result = %q, want %q", result, want)
		}
	})

	t.Run("sums numeric values", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("a", func(scope TreeLeafScope) {})
			scope.AddLeaf("b", func(scope TreeLeafScope) {})
			scope.AddLeaf("c", func(scope TreeLeafScope) {})
		})

		count := FoldIn(tr, 0, func(carry int, item Tree) int {
			return carry + 1
		})

		if count != 4 {
			t.Errorf("FoldIn count = %d, want 4", count)
		}
	})

	t.Run("deeply nested tree", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("n1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddNode("n2", func(scope TreeNodeScope) {
						scope.AddChild(func(scope TreeScope) {
							scope.AddNode("n3", func(scope TreeNodeScope) {
								scope.AddChild(func(scope TreeScope) {
									scope.AddLeaf("deep", func(scope TreeLeafScope) {})
								})
							})
						})
					})
				})
			})
		})

		var visited []string
		result := FoldIn(tr, visited, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		want := []string{"deep", "n3", "n2", "n1", "root"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldIn order = %v, want %v", result, want)
		}
	})

	t.Run("empty node with no children", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("empty", func(scope TreeNodeScope) {})
		})

		var visited []string
		result := FoldIn(tr, visited, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		want := []string{"empty", "root"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldIn order = %v, want %v", result, want)
		}
	})
}

func TestFoldOut(t *testing.T) {
	t.Run("single leaf tree", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		var visited []string
		result := FoldOut(tr, visited, func(item Tree, carry []string) []string {
			return append(carry, item.Label())
		})

		want := []string{"root", "leaf1"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldOut order = %v, want %v", result, want)
		}
	})

	t.Run("single node tree", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {})
		})

		var visited []string
		result := FoldOut(tr, visited, func(item Tree, carry []string) []string {
			return append(carry, item.Label())
		})

		want := []string{"root", "node1"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldOut order = %v, want %v", result, want)
		}
	})

	t.Run("multiple children - processes left to right", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
			scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
			scope.AddLeaf("leaf3", func(scope TreeLeafScope) {})
		})

		var visited []string
		result := FoldOut(tr, visited, func(item Tree, carry []string) []string {
			return append(carry, item.Label())
		})

		want := []string{"root", "leaf1", "leaf2", "leaf3"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldOut order = %v, want %v", result, want)
		}
	})

	t.Run("nested nodes - breadth first", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("level1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
					scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
				})
			})
		})

		var visited []string
		result := FoldOut(tr, visited, func(item Tree, carry []string) []string {
			return append(carry, item.Label())
		})

		want := []string{"root", "level1", "leaf1", "leaf2"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldOut order = %v, want %v", result, want)
		}
	})

	t.Run("complex tree", func(t *testing.T) {
		// Structure:
		//       root
		//      /    \
		//    node1  node2
		//    /  \     \
		//  l1   l2    l3
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("l1", func(scope TreeLeafScope) {})
					scope.AddLeaf("l2", func(scope TreeLeafScope) {})
				})
			})
			scope.AddNode("node2", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("l3", func(scope TreeLeafScope) {})
				})
			})
		})

		var visited []string
		result := FoldOut(tr, visited, func(item Tree, carry []string) []string {
			return append(carry, item.Label())
		})

		want := []string{"root", "node1", "l1", "l2", "node2", "l3"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldOut order = %v, want %v", result, want)
		}
	})

	t.Run("accumulates values correctly", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("a", func(scope TreeLeafScope) {})
			scope.AddLeaf("b", func(scope TreeLeafScope) {})
			scope.AddLeaf("c", func(scope TreeLeafScope) {})
		})

		result := FoldOut(tr, "", func(item Tree, carry string) string {
			return carry + item.Label()
		})

		want := "rootabc"
		if result != want {
			t.Errorf("FoldOut result = %q, want %q", result, want)
		}
	})

	t.Run("sums numeric values", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("a", func(scope TreeLeafScope) {})
			scope.AddLeaf("b", func(scope TreeLeafScope) {})
			scope.AddLeaf("c", func(scope TreeLeafScope) {})
		})

		count := FoldOut(tr, 0, func(item Tree, carry int) int {
			return carry + 1
		})

		if count != 4 {
			t.Errorf("FoldOut count = %d, want 4", count)
		}
	})

	t.Run("deeply nested tree", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("n1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddNode("n2", func(scope TreeNodeScope) {
						scope.AddChild(func(scope TreeScope) {
							scope.AddNode("n3", func(scope TreeNodeScope) {
								scope.AddChild(func(scope TreeScope) {
									scope.AddLeaf("deep", func(scope TreeLeafScope) {})
								})
							})
						})
					})
				})
			})
		})

		var visited []string
		result := FoldOut(tr, visited, func(item Tree, carry []string) []string {
			return append(carry, item.Label())
		})

		want := []string{"root", "n1", "n2", "n3", "deep"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldOut order = %v, want %v", result, want)
		}
	})

	t.Run("empty node with no children", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("empty", func(scope TreeNodeScope) {})
		})

		var visited []string
		result := FoldOut(tr, visited, func(item Tree, carry []string) []string {
			return append(carry, item.Label())
		})

		want := []string{"root", "empty"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldOut order = %v, want %v", result, want)
		}
	})
}

func TestFoldInVsFoldOut(t *testing.T) {
	t.Run("same result for commutative operations", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("a", func(scope TreeLeafScope) {})
			scope.AddNode("b", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("c", func(scope TreeLeafScope) {})
				})
			})
		})

		countIn := FoldIn(tr, 0, func(carry int, item Tree) int {
			return carry + 1
		})

		countOut := FoldOut(tr, 0, func(item Tree, carry int) int {
			return carry + 1
		})

		if countIn != countOut {
			t.Errorf("counts differ: FoldIn=%d, FoldOut=%d", countIn, countOut)
		}

		if countIn != 4 {
			t.Errorf("count = %d, want 4", countIn)
		}
	})

	t.Run("different results for non-commutative operations", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("a", func(scope TreeLeafScope) {})
			scope.AddLeaf("b", func(scope TreeLeafScope) {})
		})

		resultIn := FoldIn(tr, "", func(carry string, item Tree) string {
			return carry + item.Label()
		})

		resultOut := FoldOut(tr, "", func(item Tree, carry string) string {
			return carry + item.Label()
		})

		if resultIn == resultOut {
			t.Errorf("results should differ for non-commutative op, both = %q", resultIn)
		}

		if resultIn != "baroot" {
			t.Errorf("FoldIn result = %q, want %q", resultIn, "baroot")
		}

		if resultOut != "rootab" {
			t.Errorf("FoldOut result = %q, want %q", resultOut, "rootab")
		}
	})
}

func TestFoldEdgeCases(t *testing.T) {
	t.Run("FoldIn with empty string accumulator", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf", func(scope TreeLeafScope) {})
		})

		result := FoldIn(tr, "", func(carry string, item Tree) string {
			if carry == "" {
				return item.Label()
			}
			return carry + "," + item.Label()
		})

		want := "leaf,root"
		if result != want {
			t.Errorf("FoldIn result = %q, want %q", result, want)
		}
	})

	t.Run("FoldOut with empty string accumulator", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf", func(scope TreeLeafScope) {})
		})

		result := FoldOut(tr, "", func(item Tree, carry string) string {
			if carry == "" {
				return item.Label()
			}
			return carry + "," + item.Label()
		})

		want := "root,leaf"
		if result != want {
			t.Errorf("FoldOut result = %q, want %q", result, want)
		}
	})

	t.Run("FoldIn with struct accumulator", func(t *testing.T) {
		type Stats struct {
			Nodes int
			Leaves int
		}

		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
				})
			})
			scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
		})

		result := FoldIn(tr, Stats{}, func(carry Stats, item Tree) Stats {
			_, isLeaf := item.(*TreeLeaf)
			if isLeaf {
				carry.Leaves++
			} else {
				carry.Nodes++
			}
			return carry
		})

		if result.Leaves != 2 {
			t.Errorf("Leaves = %d, want 2", result.Leaves)
		}
		if result.Nodes != 2 {
			t.Errorf("Nodes = %d, want 2", result.Nodes)
		}
	})

	t.Run("FoldOut with struct accumulator", func(t *testing.T) {
		type Stats struct {
			Nodes int
			Leaves int
		}

		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
				})
			})
			scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
		})

		result := FoldOut(tr, Stats{}, func(item Tree, carry Stats) Stats {
			_, isLeaf := item.(*TreeLeaf)
			if isLeaf {
				carry.Leaves++
			} else {
				carry.Nodes++
			}
			return carry
		})

		if result.Leaves != 2 {
			t.Errorf("Leaves = %d, want 2", result.Leaves)
		}
		if result.Nodes != 2 {
			t.Errorf("Nodes = %d, want 2", result.Nodes)
		}
	})

	t.Run("FoldIn stops on condition", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("target", func(scope TreeLeafScope) {})
			scope.AddLeaf("other", func(scope TreeLeafScope) {})
		})

		found := FoldIn(tr, false, func(carry bool, item Tree) bool {
			if item.Label() == "target" {
				return true
			}
			return carry
		})

		if !found {
			t.Error("Should have found 'target' node")
		}
	})

	t.Run("FoldOut stops on condition", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("target", func(scope TreeLeafScope) {})
			scope.AddLeaf("other", func(scope TreeLeafScope) {})
		})

		found := FoldOut(tr, false, func(item Tree, carry bool) bool {
			if item.Label() == "target" {
				return true
			}
			return carry
		})

		if !found {
			t.Error("Should have found 'target' node")
		}
	})

	t.Run("FoldIn with map accumulator", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("n1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("l1", func(scope TreeLeafScope) {})
				})
			})
		})

		result := FoldIn(tr, map[string]int{}, func(carry map[string]int, item Tree) map[string]int {
			carry[item.Label()] = len(item.ID().String())
			return carry
		})

		if len(result) != 3 {
			t.Errorf("map length = %d, want 3", len(result))
		}

		if _, ok := result["root"]; !ok {
			t.Error("missing 'root' entry")
		}
		if _, ok := result["n1"]; !ok {
			t.Error("missing 'n1' entry")
		}
		if _, ok := result["l1"]; !ok {
			t.Error("missing 'l1' entry")
		}
	})

	t.Run("FoldOut with map accumulator", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("n1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("l1", func(scope TreeLeafScope) {})
				})
			})
		})

		result := FoldOut(tr, map[string]int{}, func(item Tree, carry map[string]int) map[string]int {
			carry[item.Label()] = len(item.ID().String())
			return carry
		})

		if len(result) != 3 {
			t.Errorf("map length = %d, want 3", len(result))
		}

		if _, ok := result["root"]; !ok {
			t.Error("missing 'root' entry")
		}
		if _, ok := result["n1"]; !ok {
			t.Error("missing 'n1' entry")
		}
		if _, ok := result["l1"]; !ok {
			t.Error("missing 'l1' entry")
		}
	})

	t.Run("FoldIn with mixed node and leaf children", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
			scope.AddNode("node1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
				})
			})
			scope.AddLeaf("leaf3", func(scope TreeLeafScope) {})
		})

		var visited []string
		result := FoldIn(tr, visited, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		want := []string{"leaf3", "leaf2", "node1", "leaf1", "root"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldIn order = %v, want %v", result, want)
		}
	})

	t.Run("FoldOut with mixed node and leaf children", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
			scope.AddNode("node1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
				})
			})
			scope.AddLeaf("leaf3", func(scope TreeLeafScope) {})
		})

		var visited []string
		result := FoldOut(tr, visited, func(item Tree, carry []string) []string {
			return append(carry, item.Label())
		})

		want := []string{"root", "leaf1", "node1", "leaf2", "leaf3"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("FoldOut order = %v, want %v", result, want)
		}
	})
}

func TestFoldTypeSafety(t *testing.T) {
	t.Run("FoldIn preserves type", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("a", func(scope TreeLeafScope) {})
			scope.AddLeaf("b", func(scope TreeLeafScope) {})
		})

		result := FoldIn(tr, []string{}, func(carry []string, item Tree) []string {
			return append(carry, item.Label())
		})

		if len(result) != 3 {
			t.Errorf("len(result) = %d, want 3", len(result))
		}
	})

	t.Run("FoldOut preserves type", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("a", func(scope TreeLeafScope) {})
			scope.AddLeaf("b", func(scope TreeLeafScope) {})
		})

		result := FoldOut(tr, 0, func(item Tree, carry int) int {
			return carry + len(item.Label())
		})

		if result != 6 {
			t.Errorf("result = %d, want 6", result)
		}
	})

	t.Run("FoldIn with custom type", func(t *testing.T) {
		type Path struct {
			Nodes []string
		}

		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("n1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("l1", func(scope TreeLeafScope) {})
				})
			})
		})

		result := FoldIn(tr, Path{}, func(carry Path, item Tree) Path {
			carry.Nodes = append(carry.Nodes, item.Label())
			return carry
		})

		if len(result.Nodes) != 3 {
			t.Errorf("len(Nodes) = %d, want 3", len(result.Nodes))
		}
	})
}
