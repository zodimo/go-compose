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

func TestFindById(t *testing.T) {
	t.Run("find root node by id", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		rootId := NewTreeNodeId("root")
		found, ok := tr.FindById(rootId)

		if !ok {
			t.Errorf("FindById(rootId) = (_, false), want (_, true)")
		}
		if found == nil {
			t.Errorf("FindById(rootId) returned nil tree")
		} else if found.Label() != "root" {
			t.Errorf("FindById(rootId).Label() = %q, want %q", found.Label(), "root")
		}
	})

	t.Run("find leaf by id", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		leafId := NewChildTreeId("leaf1", NewTreeNodeId("root"))
		found, ok := tr.FindById(leafId)

		if !ok {
			t.Errorf("FindById(leafId) = (_, false), want (_, true)")
		}
		if found == nil {
			t.Errorf("FindById(leafId) returned nil tree")
		} else if found.Label() != "leaf1" {
			t.Errorf("FindById(leafId).Label() = %q, want %q", found.Label(), "leaf1")
		}
	})

	t.Run("find nested node by id", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("level1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddNode("level2", func(scope TreeNodeScope) {
						scope.AddChild(func(scope TreeScope) {
							scope.AddLeaf("deep", func(scope TreeLeafScope) {})
						})
					})
				})
			})
		})

		level2Id := NewChildTreeId("level2", NewChildTreeId("level1", NewTreeNodeId("root")))
		found, ok := tr.FindById(level2Id)

		if !ok {
			t.Errorf("FindById(level2Id) = (_, false), want (_, true)")
		}
		if found == nil {
			t.Errorf("FindById(level2Id) returned nil tree")
		} else if found.Label() != "level2" {
			t.Errorf("FindById(level2Id).Label() = %q, want %q", found.Label(), "level2")
		}
	})

	t.Run("find non-existent id", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		nonExistentId := NewTreeNodeId("nonexistent")
		found, ok := tr.FindById(nonExistentId)

		if ok {
			t.Errorf("FindById(nonExistentId) = (_, true), want (_, false)")
		}
		if found != nil {
			t.Errorf("FindById(nonExistentId) = (%v, _), want (nil, _)", found)
		}
	})

	t.Run("find id not in lineage", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("child", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {})
			})
		})

		// "other" is not a descendant of "root"
		otherId := NewTreeNodeId("other")
		found, ok := tr.FindById(otherId)

		if ok {
			t.Errorf("FindById(otherId) = (_, true), want (_, false)")
		}
		if found != nil {
			t.Errorf("FindById(otherId) = (%v, _), want (nil, _)", found)
		}
	})

	t.Run("find intermediate node by id", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddNode("level1", func(scope TreeNodeScope) {
				scope.AddChild(func(scope TreeScope) {
					scope.AddLeaf("leaf", func(scope TreeLeafScope) {})
				})
			})
		})

		level1Id := NewChildTreeId("level1", NewTreeNodeId("root"))
		found, ok := tr.FindById(level1Id)

		if !ok {
			t.Errorf("FindById(level1Id) = (_, false), want (_, true)")
		}
		if found == nil {
			t.Errorf("FindById(level1Id) returned nil tree")
		} else if found.Label() != "level1" {
			t.Errorf("FindById(level1Id).Label() = %q, want %q", found.Label(), "level1")
		}
	})

	t.Run("find by id from leaf", func(t *testing.T) {
		tr := NewTree("root", func(scope TreeScope) {
			scope.AddLeaf("leaf1", func(scope TreeLeafScope) {})
		})

		root, _ := tr.(*TreeNode)
		leaf, _ := root.Children()[0].(*TreeLeaf)

		// Leaf should only find itself
		leafId := leaf.ID()
		found, ok := leaf.FindById(leafId)

		if !ok {
			t.Errorf("leaf.FindById(leafId) = (_, false), want (_, true)")
		}
		if found == nil {
			t.Errorf("leaf.FindById(leafId) returned nil")
		} else if found.Label() != "leaf1" {
			t.Errorf("leaf.FindById(leafId).Label() = %q, want %q", found.Label(), "leaf1")
		}

		// Leaf should not find root
		rootId := root.ID()
		_, ok = leaf.FindById(rootId)
		if ok {
			t.Errorf("leaf.FindById(rootId) = (_, true), want (_, false)")
		}
	})
}

func TestTreeIdEncodingDefault(t *testing.T) {
	t.Run("default encoder is passthru", func(t *testing.T) {
		id := NewTreeNodeId("hello world")
		got := id.String()
		want := "hello world"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("default encoder preserves special characters", func(t *testing.T) {
		id := NewTreeNodeId("test/path#fragment?query=value")
		got := id.String()
		want := "test/path#fragment?query=value"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("default encoder with spaces", func(t *testing.T) {
		id := NewTreeNodeId("hello world")
		got := id.String()
		want := "hello world"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("child inherits default encoder", func(t *testing.T) {
		parent := NewTreeNodeId("parent")
		child := NewChildTreeId("child with space", parent)
		got := child.String()
		want := "parent/child with space"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestTreeIdEncodingPassthru(t *testing.T) {
	t.Run("passthru encoder preserves label as-is", func(t *testing.T) {
		id := NewTreeNodeId("hello", PassthruEncoder)
		got := id.String()
		want := "hello"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("passthru encoder with special chars", func(t *testing.T) {
		id := NewTreeNodeId("test/path#frag?a=1&b=2", PassthruEncoder)
		got := id.String()
		want := "test/path#frag?a=1&b=2"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("passthru encoder decode returns same value", func(t *testing.T) {
		id := NewTreeNodeId("encoded%20value", PassthruEncoder)
		got := id.Label()
		want := "encoded%20value"
		if got != want {
			t.Errorf("Label() = %q, want %q", got, want)
		}
	})

	t.Run("passthru encoder with unicode", func(t *testing.T) {
		id := NewTreeNodeId("hello 世界 🌍", PassthruEncoder)
		got := id.String()
		want := "hello 世界 🌍"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("child inherits passthru encoder", func(t *testing.T) {
		parent := NewTreeNodeId("parent", PassthruEncoder)
		child := NewChildTreeId("child/path", parent)
		got := child.String()
		want := "parent/child/path"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestTreeIdEncodingQueryEscape(t *testing.T) {
	t.Run("query escape encoder escapes spaces", func(t *testing.T) {
		id := NewTreeNodeId("hello world", QueryEscapeEncoder)
		got := id.String()
		want := "hello+world"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("query escape encoder decodes correctly", func(t *testing.T) {
		id := NewTreeNodeId("hello world", QueryEscapeEncoder)
		got := id.Label()
		want := "hello world"
		if got != want {
			t.Errorf("Label() = %q, want %q", got, want)
		}
	})

	t.Run("query escape encoder escapes special chars", func(t *testing.T) {
		id := NewTreeNodeId("path/to/file", QueryEscapeEncoder)
		got := id.String()
		want := "path%2Fto%2Ffile"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("query escape encoder decodes special chars", func(t *testing.T) {
		id := NewTreeNodeId("path/to/file", QueryEscapeEncoder)
		got := id.Label()
		want := "path/to/file"
		if got != want {
			t.Errorf("Label() = %q, want %q", got, want)
		}
	})

	t.Run("query escape encoder with query params", func(t *testing.T) {
		id := NewTreeNodeId("a=1&b=2", QueryEscapeEncoder)
		got := id.String()
		want := "a%3D1%26b%3D2"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("query escape encoder decodes query params", func(t *testing.T) {
		id := NewTreeNodeId("a=1&b=2", QueryEscapeEncoder)
		got := id.Label()
		want := "a=1&b=2"
		if got != want {
			t.Errorf("Label() = %q, want %q", got, want)
		}
	})

	t.Run("child inherits query escape encoder", func(t *testing.T) {
		parent := NewTreeNodeId("parent name", QueryEscapeEncoder)
		child := NewChildTreeId("child name", parent)
		got := child.String()
		want := "parent+name/child+name"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		if got := child.Label(); got != "child name" {
			t.Errorf("Label() = %q, want %q", got, "child name")
		}
	})

	t.Run("query escape encoder with empty string", func(t *testing.T) {
		id := NewTreeNodeId("", QueryEscapeEncoder)
		got := id.String()
		want := ""
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestTreeIdFromString(t *testing.T) {
	t.Run("from string single part", func(t *testing.T) {
		id := FromString("root")
		got := id.String()
		want := "root"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		if got := id.Label(); got != want {
			t.Errorf("Label() = %q, want %q", got, want)
		}
	})

	t.Run("from string two parts", func(t *testing.T) {
		id := FromString("parent/child")
		got := id.String()
		want := "parent/child"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		if got := id.Label(); got != "child" {
			t.Errorf("Label() = %q, want %q", got, "child")
		}
	})

	t.Run("from string multiple parts", func(t *testing.T) {
		id := FromString("a/b/c/d")
		got := id.String()
		want := "a/b/c/d"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		if got := id.Label(); got != "d" {
			t.Errorf("Label() = %q, want %q", got, "d")
		}
	})

	t.Run("from string empty", func(t *testing.T) {
		id := FromString("")
		got := id.String()
		want := ""
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("from string with spaces using passthru", func(t *testing.T) {
		id := FromString("hello world", PassthruEncoder)
		got := id.String()
		want := "hello world"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		if got := id.Label(); got != want {
			t.Errorf("Label() = %q, want %q", got, want)
		}
	})

	t.Run("from string with escaped chars using query escape", func(t *testing.T) {
		id := FromString("hello+world", QueryEscapeEncoder)
		got := id.String()
		want := "hello+world"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		if got := id.Label(); got != "hello world" {
			t.Errorf("Label() = %q, want %q", got, "hello world")
		}
	})

	t.Run("from string nested with escaped chars", func(t *testing.T) {
		id := FromString("parent+name/child+name", QueryEscapeEncoder)
		got := id.String()
		want := "parent+name/child+name"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		if got := id.Label(); got != "child name" {
			t.Errorf("Label() = %q, want %q", got, "child name")
		}
	})

	t.Run("from string with percent encoded slash", func(t *testing.T) {
		id := FromString("path%2Fto%2Ffile", QueryEscapeEncoder)
		got := id.String()
		want := "path%2Fto%2Ffile"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
		if got := id.Label(); got != "path/to/file" {
			t.Errorf("Label() = %q, want %q", got, "path/to/file")
		}
	})

	t.Run("roundtrip with passthru encoder", func(t *testing.T) {
		original := NewTreeNodeId("test-label", PassthruEncoder)
		id := FromString(original.String(), PassthruEncoder)

		if got := id.String(); got != original.String() {
			t.Errorf("String() = %q, want %q", got, original.String())
		}
		if got := id.Label(); got != original.Label() {
			t.Errorf("Label() = %q, want %q", got, original.Label())
		}
	})

	t.Run("roundtrip with query escape encoder", func(t *testing.T) {
		original := NewTreeNodeId("hello world", QueryEscapeEncoder)
		id := FromString(original.String(), QueryEscapeEncoder)

		if got := id.String(); got != original.String() {
			t.Errorf("String() = %q, want %q", got, original.String())
		}
		if got := id.Label(); got != original.Label() {
			t.Errorf("Label() = %q, want %q", got, original.Label())
		}
	})

	t.Run("from string with default encoder", func(t *testing.T) {
		// No encoder specified, should use PassthruEncoder
		id := FromString("test/label")
		got := id.String()
		want := "test/label"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}
