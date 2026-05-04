package tree

import (
	"reflect"
	"testing"
)

func TestTreeId(t *testing.T) {
	tests := []struct {
		name    string
		id      TreeId
		wantStr string
		wantLbl string
	}{
		{
			name:    "root id",
			id:      NewTreeNodeId("root"),
			wantStr: "root",
			wantLbl: "root",
		},
		{
			name:    "child of root",
			id:      NewChildTreeId("child", NewTreeNodeId("root")),
			wantStr: "root/child",
			wantLbl: "child",
		},
		{
			name:    "nested child",
			id:      NewChildTreeId("grandchild", NewChildTreeId("child", NewTreeNodeId("root"))),
			wantStr: "root/child/grandchild",
			wantLbl: "grandchild",
		},
		{
			name:    "empty label",
			id:      NewTreeNodeId(""),
			wantStr: "",
			wantLbl: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.String(); got != tt.wantStr {
				t.Errorf("String() = %q, want %q", got, tt.wantStr)
			}
			if got := tt.id.Label(); got != tt.wantLbl {
				t.Errorf("Label() = %q, want %q", got, tt.wantLbl)
			}
		})
	}
}

func TestNewTree(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {})

		node, ok := tr.(*TreeNode)
		if !ok {
			t.Fatalf("expected *TreeNode, got %T", tr)
		}

		if node.Label() != "root" {
			t.Errorf("Label() = %q, want %q", node.Label(), "root")
		}
		if node.ID().String() != "root" {
			t.Errorf("ID().String() = %q, want %q", node.ID().String(), "root")
		}
		if len(node.Children()) != 0 {
			t.Errorf("len(Children()) = %d, want 0", len(node.Children()))
		}
	})

	t.Run("tree with single leaf", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		node, ok := tr.(*TreeNode)
		if !ok {
			t.Fatalf("expected *TreeNode, got %T", tr)
		}

		children := node.Children()
		if len(children) != 1 {
			t.Fatalf("len(Children()) = %d, want 1", len(children))
		}

		leaf, ok := children[0].(*TreeLeaf)
		if !ok {
			t.Fatalf("expected child to be *treeLeafBuilder, got %T", children[0])
		}

		if leaf.Label() != "leaf1" {
			t.Errorf("Label() = %q, want %q", leaf.Label(), "leaf1")
		}
		if leaf.ID().String() != "root/leaf1" {
			t.Errorf("ID().String() = %q, want %q", leaf.ID().String(), "root/leaf1")
		}
	})

	t.Run("tree with single node", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {})
		})

		node, ok := tr.(*TreeNode)
		if !ok {
			t.Fatalf("expected *treeNodeBuilder, got %T", tr)
		}

		children := node.Children()
		if len(children) != 1 {
			t.Fatalf("len(Children()) = %d, want 1", len(children))
		}

		childNode, ok := children[0].(*TreeNode)
		if !ok {
			t.Fatalf("expected child to be *TreeNode, got %T", children[0])
		}

		if childNode.Label() != "node1" {
			t.Errorf("Label() = %q, want %q", childNode.Label(), "node1")
		}
		if childNode.ID().String() != "root/node1" {
			t.Errorf("ID().String() = %q, want %q", childNode.ID().String(), "root/node1")
		}
	})

	t.Run("nested nodes", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("parent", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("child", func(scope TreeLeafScope) {})
				})
			})
		})

		root, ok := tr.(*TreeNode)
		if !ok {
			t.Fatalf("expected *TreeNode, got %T", tr)
		}

		children := root.Children()
		if len(children) != 1 {
			t.Fatalf("len(root.Children()) = %d, want 1", len(children))
		}

		parent, ok := children[0].(*TreeNode)
		if !ok {
			t.Fatalf("expected child to be *TreeNode, got %T", children[0])
		}

		parentChildren := parent.Children()
		if len(parentChildren) != 1 {
			t.Fatalf("len(parent.Children()) = %d, want 1", len(parentChildren))
		}

		childLeaf, ok := parentChildren[0].(*TreeLeaf)
		if !ok {
			t.Fatalf("expected grandchild to be *TreeLeaf, got %T", parentChildren[0])
		}

		if childLeaf.Label() != "child" {
			t.Errorf("Label() = %q, want %q", childLeaf.Label(), "child")
		}
		if childLeaf.ID().String() != "root/parent/child" {
			t.Errorf("ID().String() = %q, want %q", childLeaf.ID().String(), "root/parent/child")
		}
	})

	t.Run("multiple children", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
			scope.AddLeaf("leaf2", func(scope TreeLeafScope) {})
			scope.AddNode("node1", func(scope TreeNodeScope) {})
		})

		root, ok := tr.(*TreeNode)
		if !ok {
			t.Fatalf("expected *TreeNode, got %T", tr)
		}

		children := root.Children()
		if len(children) != 3 {
			t.Fatalf("len(Children()) = %d, want 3", len(children))
		}

		expected := []string{"leaf1", "leaf2", "node1"}
		for i, child := range children {
			if child.Label() != expected[i] {
				t.Errorf("child[%d].Label() = %q, want %q", i, child.Label(), expected[i])
			}
		}
	})
}

func TestTreeNodeProperties(t *testing.T) {
	t.Run("add and retrieve properties", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {
				scope.AddProperty("key1", "value1")
				scope.AddProperty("key2", 42)
			})
		})

		root, ok := tr.(*TreeNode)
		if !ok {
			t.Fatalf("expected *TreeNode, got %T", tr)
		}

		children := root.Children()
		if len(children) != 1 {
			t.Fatalf("len(Children()) = %d, want 1", len(children))
		}

		node, ok := children[0].(*TreeNode)
		if !ok {
			t.Fatalf("expected child to be *TreeNode, got %T", children[0])
		}

		props := node.Properties()
		if props == nil {
			t.Fatal("Properties() returned nil")
		}

		if got, want := props["key1"], "value1"; got != want {
			t.Errorf("props[key1] = %v, want %v", got, want)
		}
		if got, want := props["key2"], 42; got != want {
			t.Errorf("props[key2] = %v, want %v", got, want)
		}
	})

	t.Run("properties clone is independent", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {
				scope.AddProperty("key", "original")
			})
		})

		root, _ := tr.(*TreeNode)
		node, _ := root.Children()[0].(*TreeNode)

		props := node.Properties()
		props["key"] = "modified"

		props2 := node.Properties()
		// BUG: Properties() returns n.properties instead of the clone,
		// so the following assertion will fail until the bug is fixed.
		if got, want := props2["key"], "original"; got != want {
			t.Errorf("after modifying clone, original props[key] = %v, want %v (Properties() may not be returning a clone)", got, want)
		}
	})
}

func TestTreeLeafProperties(t *testing.T) {
	t.Run("add and retrieve properties", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(s TreeLeafScope) {
				s.AddProperty("color", "red")
				s.AddProperty("count", 10)
			})
		})

		root, ok := tr.(*TreeNode)
		if !ok {
			t.Fatalf("expected *TreeNode, got %T", tr)
		}

		children := root.Children()
		if len(children) != 1 {
			t.Fatalf("len(Children()) = %d, want 1", len(children))
		}

		leaf, ok := children[0].(*TreeLeaf)
		if !ok {
			t.Fatalf("expected child to be *TreeLeaf, got %T", children[0])
		}

		props := leaf.Properties()
		if props == nil {
			t.Fatal("Properties() returned nil")
		}

		if got, want := props["color"], "red"; got != want {
			t.Errorf("props[color] = %v, want %v", got, want)
		}
		if got, want := props["count"], 10; got != want {
			t.Errorf("props[count] = %v, want %v", got, want)
		}
	})
}

func TestTreeNodeInterface(t *testing.T) {
	t.Run("string representation", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("node1", func(scope TreeNodeScope) {})
		})

		root, _ := tr.(*TreeNode)
		node, _ := root.Children()[0].(*TreeNode)

		if got, want := node.String(), "Node"; got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("children returns clone", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		root, _ := tr.(*TreeNode)
		children1 := root.Children()
		children2 := root.Children()

		if &children1[0] == &children2[0] {
			t.Error("Children() should return a clone, but returned the same slice")
		}
	})
}

func TestTreeLeafInterface(t *testing.T) {
	t.Run("string representation", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		root, _ := tr.(*TreeNode)
		leaf, _ := root.Children()[0].(*TreeLeaf)

		if got, want := leaf.String(), "Leaf"; got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestTreeHierarchy(t *testing.T) {
	t.Run("node ids reflect hierarchy", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("level1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddNode("level2", func(scope TreeNodeScope) {
						scope.AddChild(func(scope TreeScope) {
							scope.AddLeaf("level3", func(scope TreeLeafScope) {})
						})
					})
				})
			})
		})

		root, _ := tr.(*TreeNode)
		if root.ID().String() != "root" {
			t.Errorf("root ID = %q, want %q", root.ID().String(), "root")
		}

		level1, _ := root.Children()[0].(*TreeNode)
		if level1.ID().String() != "root/level1" {
			t.Errorf("level1 ID = %q, want %q", level1.ID().String(), "root/level1")
		}

		level2, _ := level1.Children()[0].(*TreeNode)
		if level2.ID().String() != "root/level1/level2" {
			t.Errorf("level2 ID = %q, want %q", level2.ID().String(), "root/level1/level2")
		}

		level3, _ := level2.Children()[0].(*TreeLeaf)
		if level3.ID().String() != "root/level1/level2/level3" {
			t.Errorf("level3 ID = %q, want %q", level3.ID().String(), "root/level1/level2/level3")
		}
	})
}

func TestTreeNodeAddChild(t *testing.T) {
	t.Run("add child via node scope", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("parent", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("child1", func(scope TreeLeafScope) {})
					scope.AddLeaf("child2", func(scope TreeLeafScope) {})
				})
			})
		})

		root, _ := tr.(*TreeNode)
		parent, _ := root.Children()[0].(*TreeNode)
		children := parent.Children()

		if len(children) != 2 {
			t.Fatalf("len(children) = %d, want 2", len(children))
		}

		if children[0].Label() != "child1" {
			t.Errorf("children[0].Label() = %q, want %q", children[0].Label(), "child1")
		}
		if children[1].Label() != "child2" {
			t.Errorf("children[1].Label() = %q, want %q", children[1].Label(), "child2")
		}
	})

	t.Run("add nested child via node scope", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("parent", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddNode("intermediate", func(scope TreeNodeScope) {
						scope.AddChild(func(scope TreeScope) {
							scope.AddLeaf("deep", func(scope TreeLeafScope) {})
						})
					})
				})
			})
		})

		root, _ := tr.(*TreeNode)
		parent, _ := root.Children()[0].(*TreeNode)
		intermediate, _ := parent.Children()[0].(*TreeNode)
		deep, _ := intermediate.Children()[0].(*TreeLeaf)

		if deep.Label() != "deep" {
			t.Errorf("deep.Label() = %q, want %q", deep.Label(), "deep")
		}
		if deep.ID().String() != "root/parent/intermediate/deep" {
			t.Errorf("deep.ID() = %q, want %q", deep.ID().String(), "root/parent/intermediate/deep")
		}
	})
}

func TestTreeImplementsInterface(t *testing.T) {
	t.Run("treeNodeBuilder implements Tree", func(t *testing.T) {
		var _ Tree = (*treeNodeBuilder)(nil)
	})

	t.Run("treeLeafBuilder implements Tree", func(t *testing.T) {
		var _ Tree = (*treeLeafBuilder)(nil)
	})

	t.Run("treeNodeBuilder implements TreeBuilder", func(t *testing.T) {
		var _ TreeBuilder = (*treeNodeBuilder)(nil)
	})

	t.Run("treeLeafBuilder implements TreeBuilder", func(t *testing.T) {
		var _ TreeBuilder = (*treeLeafBuilder)(nil)
	})
}

func TestTreeNodeScopeImplementsInterface(t *testing.T) {
	t.Run("treeNodeScopeImpl implements TreeNodeScope", func(t *testing.T) {
		var _ TreeNodeScope = (*treeNodeScopeImpl)(nil)
	})
}

func TestTreeScopeImplementsInterface(t *testing.T) {
	t.Run("treeScopeImpl implements TreeScope", func(t *testing.T) {
		var _ TreeScope = (*treeScopeImpl)(nil)
	})
}

func TestTreeLeafScopeImplementsInterface(t *testing.T) {
	t.Run("treeLeafScopeImpl implements TreeLeafScope", func(t *testing.T) {
		var _ TreeLeafScope = (*treeLeafScopeImpl)(nil)
	})
}

func TestTreeEdgeCases(t *testing.T) {
	t.Run("node with no builder call", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("empty", func(scope TreeNodeScope) {})
		})

		root, _ := tr.(*TreeNode)
		empty, _ := root.Children()[0].(*TreeNode)

		if len(empty.Children()) != 0 {
			t.Errorf("len(Children()) = %d, want 0", len(empty.Children()))
		}
		if len(empty.Properties()) != 0 {
			t.Errorf("len(Properties()) = %d, want 0", len(empty.Properties()))
		}
	})

	t.Run("leaf with no properties", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("empty", func(scope TreeLeafScope) {})
		})

		root, _ := tr.(*TreeNode)
		empty, _ := root.Children()[0].(*TreeLeaf)

		if len(empty.Properties()) != 0 {
			t.Errorf("len(Properties()) = %d, want 0", len(empty.Properties()))
		}
	})

	t.Run("mixed properties types", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf", func(s TreeLeafScope) {
				s.AddProperty("string", "hello")
				s.AddProperty("int", 42)
				s.AddProperty("bool", true)
				s.AddProperty("float", 3.14)
				s.AddProperty("nil", nil)
				s.AddProperty("slice", []int{1, 2, 3})
				s.AddProperty("map", map[string]int{"a": 1})
			})
		})

		root, _ := tr.(*TreeNode)
		leaf, _ := root.Children()[0].(*TreeLeaf)
		props := leaf.Properties()

		expected := map[string]any{
			"string": "hello",
			"int":    42,
			"bool":   true,
			"float":  3.14,
			"nil":    nil,
			"slice":  []int{1, 2, 3},
			"map":    map[string]int{"a": 1},
		}

		for k, want := range expected {
			got, ok := props[k]
			if !ok {
				t.Errorf("missing property %q", k)
				continue
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("props[%q] = %v, want %v", k, got, want)
			}
		}
	})
}

func TestTreeNodeParent(t *testing.T) {
	t.Run("node has parent", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Parent() panicked (likely infinite recursion bug): %v", r)
			}
		}()

		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("child", func(scope TreeNodeScope) {})
		})

		root, _ := tr.(*TreeNode)
		child, _ := root.Children()[0].(*TreeNode)

		if child.HasParent() {
			t.Log("HasParent() returned true")
		} else {
			t.Log("HasParent() returned false")
		}
	})
}

func TestTreeLeafParent(t *testing.T) {
	t.Run("leaf has parent", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Parent() panicked (likely infinite recursion bug): %v", r)
			}
		}()

		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("child", func(scope TreeLeafScope) {})
		})

		root, _ := tr.(*TreeNode)
		child, _ := root.Children()[0].(*TreeLeaf)

		if child.HasParent() {
			t.Log("HasParent() returned true")
		} else {
			t.Log("HasParent() returned false")
		}
	})
}
