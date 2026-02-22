package navigation

import (
	"testing"
	"time"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

func newMockTypedMutableValue[T any](value T) state.MutableValueTyped[T] {
	return state.MustRemember[T](compose.NewComposer(), time.Now().String(), func() T {
		return value
	})
}

func TestNavController(t *testing.T) {
	// Initialize NavController with empty stack
	backStack := newMockTypedMutableValue([]BackStackEntry{})
	nc := NewNavController(backStack)

	// Test Initial State
	if nc.CurrentEntry() != nil {
		t.Error("Expected nil CurrentEntry for empty stack")
	}

	// Test Navigate
	nc.Navigate("home")
	if nc.CurrentEntry().Route != "home" {
		t.Errorf("Expected route 'home', got '%s'", nc.CurrentEntry().Route)
	}

	nc.Navigate("details")
	if nc.CurrentEntry().Route != "details" {
		t.Errorf("Expected route 'details', got '%s'", nc.CurrentEntry().Route)
	}

	// Test PopBackStack
	popped := nc.PopBackStack()
	if !popped {
		t.Error("Expected PopBackStack to return true")
	}
	if nc.CurrentEntry().Route != "home" {
		t.Errorf("Expected route 'home' after pop, got '%s'", nc.CurrentEntry().Route)
	}

	// Test PopBackStack to empty
	popped = nc.PopBackStack()
	if !popped {
		t.Error("Expected PopBackStack to return true for last item")
	}
	if nc.CurrentEntry() != nil {
		t.Error("Expected nil CurrentEntry after popping last item")
	}

	// Test PopBackStack on empty
	popped = nc.PopBackStack()
	if popped {
		t.Error("Expected PopBackStack to return false for empty stack")
	}
}

// Simple test for builder
func TestNavGraphBuilder(t *testing.T) {
	builder := NewNavGraphBuilder()

	dummy := func(c api.Composer) api.Composer { return c }

	builder.Composable("route1", dummy)

	if _, ok := builder.destinations["route1"]; !ok {
		t.Error("Expected route1 to be present")
	}
}

// Test route pattern parsing
func TestParseRoutePattern(t *testing.T) {
	tests := []struct {
		pattern           string
		expectedPathArgs  []string
		expectedQueryArgs []string
	}{
		{"home", nil, nil},
		{"details/{id}", []string{"id"}, nil},
		{"user/{userId}/post/{postId}", []string{"userId", "postId"}, nil},
		{"search?query={query}", nil, []string{"query"}},
		{"details/{id}?filter={filter}", []string{"id"}, []string{"filter"}},
		{"edit/{id}?name={name}&active={active}", []string{"id"}, []string{"name", "active"}},
	}

	for _, tt := range tests {
		pathArgs, queryArgs := parseRoutePattern(tt.pattern)

		if len(pathArgs) != len(tt.expectedPathArgs) {
			t.Errorf("Pattern %q: expected %d path args, got %d", tt.pattern, len(tt.expectedPathArgs), len(pathArgs))
			continue
		}
		for i, expected := range tt.expectedPathArgs {
			if pathArgs[i] != expected {
				t.Errorf("Pattern %q: path arg %d expected %q, got %q", tt.pattern, i, expected, pathArgs[i])
			}
		}

		if len(queryArgs) != len(tt.expectedQueryArgs) {
			t.Errorf("Pattern %q: expected %d query args, got %d", tt.pattern, len(tt.expectedQueryArgs), len(queryArgs))
			continue
		}
		for i, expected := range tt.expectedQueryArgs {
			if queryArgs[i] != expected {
				t.Errorf("Pattern %q: query arg %d expected %q, got %q", tt.pattern, i, expected, queryArgs[i])
			}
		}
	}
}

// Test route matching with path arguments
func TestMatchRoutePathArgs(t *testing.T) {
	tests := []struct {
		pattern      string
		route        string
		shouldMatch  bool
		expectedArgs NavArguments
	}{
		{"home", "home", true, NavArguments{}},
		{"home", "settings", false, nil},
		{"details/{id}", "details/123", true, NavArguments{"id": "123"}},
		{"details/{id}", "details", false, nil},
		{"user/{userId}/post/{postId}", "user/42/post/99", true, NavArguments{"userId": "42", "postId": "99"}},
	}

	for _, tt := range tests {
		args, matched := matchRoute(tt.pattern, tt.route)

		if matched != tt.shouldMatch {
			t.Errorf("Pattern %q, Route %q: expected match=%v, got %v", tt.pattern, tt.route, tt.shouldMatch, matched)
			continue
		}

		if matched {
			for key, expected := range tt.expectedArgs {
				if got, ok := args[key]; !ok || got != expected {
					t.Errorf("Pattern %q, Route %q: arg %q expected %q, got %q", tt.pattern, tt.route, key, expected, got)
				}
			}
		}
	}
}

// Test route matching with query arguments
func TestMatchRouteQueryArgs(t *testing.T) {
	tests := []struct {
		pattern      string
		route        string
		expectedArgs NavArguments
	}{
		{"search?query={query}", "search?query=hello", NavArguments{"query": "hello"}},
		{"details/{id}?filter={filter}", "details/123?filter=active", NavArguments{"id": "123", "filter": "active"}},
		{"details/{id}?filter={filter}", "details/123", NavArguments{"id": "123"}}, // Query arg is optional
	}

	for _, tt := range tests {
		args, matched := matchRoute(tt.pattern, tt.route)

		if !matched {
			t.Errorf("Pattern %q, Route %q: expected match, got no match", tt.pattern, tt.route)
			continue
		}

		for key, expected := range tt.expectedArgs {
			if got, ok := args[key]; !ok || got != expected {
				t.Errorf("Pattern %q, Route %q: arg %q expected %q, got %v", tt.pattern, tt.route, key, expected, got)
			}
		}
	}
}

// Test NavArguments type-safe getters
func TestNavArgumentsGetters(t *testing.T) {
	args := NavArguments{
		"stringVal": "hello",
		"intVal":    42,
		"intStr":    "123",
		"boolVal":   true,
		"boolStr":   "true",
	}

	// Test GetString
	if s, ok := args.GetString("stringVal"); !ok || s != "hello" {
		t.Errorf("GetString: expected 'hello', got '%s', ok=%v", s, ok)
	}
	if _, ok := args.GetString("missing"); ok {
		t.Error("GetString: expected ok=false for missing key")
	}

	// Test GetInt with int value
	if i, ok := args.GetInt("intVal"); !ok || i != 42 {
		t.Errorf("GetInt: expected 42, got %d, ok=%v", i, ok)
	}

	// Test GetInt with string value (should parse)
	if i, ok := args.GetInt("intStr"); !ok || i != 123 {
		t.Errorf("GetInt from string: expected 123, got %d, ok=%v", i, ok)
	}

	// Test GetBool with bool value
	if b, ok := args.GetBool("boolVal"); !ok || !b {
		t.Errorf("GetBool: expected true, got %v, ok=%v", b, ok)
	}

	// Test GetBool with string value (should parse)
	if b, ok := args.GetBool("boolStr"); !ok || !b {
		t.Errorf("GetBool from string: expected true, got %v, ok=%v", b, ok)
	}
}

// Test build route
func TestBuildRoute(t *testing.T) {
	tests := []struct {
		pattern  string
		args     NavArguments
		expected string
	}{
		{"home", NavArguments{}, "home"},
		{"details/{id}", NavArguments{"id": "123"}, "details/123"},
		{"user/{userId}/post/{postId}", NavArguments{"userId": "42", "postId": "99"}, "user/42/post/99"},
		{"search?query={query}", NavArguments{"query": "hello"}, "search?query=hello"},
		{"details/{id}?filter={filter}", NavArguments{"id": "123", "filter": "active"}, "details/123?filter=active"},
	}

	for _, tt := range tests {
		result := buildRoute(tt.pattern, tt.args)
		if result != tt.expected {
			t.Errorf("buildRoute(%q, %v): expected %q, got %q", tt.pattern, tt.args, tt.expected, result)
		}
	}
}

// Test NavGraphBuilder.findDestination
func TestNavGraphBuilderFindDestination(t *testing.T) {
	builder := NewNavGraphBuilder()

	// Register destinations
	builder.ComposableWithArgs("home", func(entry *BackStackEntry) api.Composable {
		return func(c api.Composer) api.Composer { return c }
	})
	builder.ComposableWithArgs("details/{id}", func(entry *BackStackEntry) api.Composable {
		return func(c api.Composer) api.Composer { return c }
	})

	// Test exact match
	_, args, ok := builder.findDestination("home")
	if !ok {
		t.Error("Expected to find 'home' destination")
	}
	if len(args) != 0 {
		t.Errorf("Expected no args for 'home', got %v", args)
	}

	// Test pattern match with args
	_, args, ok = builder.findDestination("details/123")
	if !ok {
		t.Error("Expected to find 'details/123' destination")
	}
	if id, ok := args.GetString("id"); !ok || id != "123" {
		t.Errorf("Expected id='123', got '%s'", id)
	}

	// Test no match
	_, _, ok = builder.findDestination("unknown/route")
	if ok {
		t.Error("Expected no match for 'unknown/route'")
	}
}

// Test backward compatibility - Composable still works
func TestBackwardCompatibility(t *testing.T) {
	builder := NewNavGraphBuilder()

	dummy := func(c api.Composer) api.Composer { return c }

	// Use legacy Composable method
	builder.Composable("home", dummy)

	// Should still be findable
	_, args, ok := builder.findDestination("home")
	if !ok {
		t.Error("Expected to find 'home' destination using legacy Composable")
	}
	if len(args) != 0 {
		t.Errorf("Expected no args for legacy route, got %v", args)
	}
}

// Test URL encoding/decoding of special characters
func TestSpecialCharactersInArguments(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		args    NavArguments
		encoded string
		decoded NavArguments
	}{
		{
			name:    "space in path",
			pattern: "search/{query}",
			args:    NavArguments{"query": "hello world"},
			encoded: "search/hello%20world",
			decoded: NavArguments{"query": "hello world"},
		},
		{
			name:    "slash in path (encoded)",
			pattern: "file/{path}",
			args:    NavArguments{"path": "dir/subdir/file.txt"},
			encoded: "file/dir%2Fsubdir%2Ffile.txt",
			decoded: NavArguments{"path": "dir/subdir/file.txt"},
		},
		{
			name:    "ampersand in query",
			pattern: "search?q={query}",
			args:    NavArguments{"query": "foo&bar"},
			encoded: "search?q=foo%26bar",
			decoded: NavArguments{"query": "foo&bar"},
		},
		{
			name:    "unicode characters",
			pattern: "user/{name}",
			args:    NavArguments{"name": "日本語"},
			encoded: "user/%E6%97%A5%E6%9C%AC%E8%AA%9E",
			decoded: NavArguments{"name": "日本語"},
		},
		{
			name:    "special URL characters",
			pattern: "data/{value}",
			args:    NavArguments{"value": "a=b&c=d"},
			encoded: "data/a=b&c=d", // Note: = and & are valid in path segments
			decoded: NavArguments{"value": "a=b&c=d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test buildRoute encodes correctly
			encoded := buildRoute(tt.pattern, tt.args)
			if encoded != tt.encoded {
				t.Errorf("buildRoute: expected %q, got %q", tt.encoded, encoded)
			}

			// Test matchRoute decodes correctly
			args, matched := matchRoute(tt.pattern, tt.encoded)
			if !matched {
				t.Errorf("matchRoute: expected match for %q", tt.encoded)
				return
			}

			for key, expected := range tt.decoded {
				if got, ok := args.GetString(key); !ok || got != expected {
					t.Errorf("matchRoute: arg %q expected %q, got %q", key, expected, got)
				}
			}
		})
	}
}

// Test round-trip: buildRoute then matchRoute preserves values
func TestRoundTripEncoding(t *testing.T) {
	testValues := []string{
		"hello world",
		"path/to/file",
		"query=value&other=123",
		"名前",
		"special!@#$%^*()",
		"encoded%20already",
	}

	pattern := "test/{value}"

	for _, value := range testValues {
		t.Run(value, func(t *testing.T) {
			// Build route with value
			args := NavArguments{"value": value}
			route := buildRoute(pattern, args)

			// Match and extract
			extracted, ok := matchRoute(pattern, route)
			if !ok {
				t.Errorf("Failed to match route %q", route)
				return
			}

			got, _ := extracted.GetString("value")
			if got != value {
				t.Errorf("Round-trip failed: input %q, got %q", value, got)
			}
		})
	}
}
