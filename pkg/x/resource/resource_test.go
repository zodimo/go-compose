package resource

import (
	"errors"
	"testing"
)

func TestResourceFactories(t *testing.T) {
	t.Run("ResourceSuccess", func(t *testing.T) {
		r := ResourceSuccess("test")
		if r.IsLoading() {
			t.Errorf("expected not loading")
		}
		if r.Error() != nil {
			t.Errorf("expected no error, got %v", r.Error())
		}
		if !r.HasData() {
			t.Errorf("expected data")
		}
		if r.Data() != "test" {
			t.Errorf("expected 'test', got %v", r.Data())
		}
	})

	t.Run("ResourceLoading", func(t *testing.T) {
		r := ResourceLoading[string]()
		if !r.IsLoading() {
			t.Errorf("expected loading")
		}
		if r.Error() != nil {
			t.Errorf("expected no error, got %v", r.Error())
		}
		if r.HasData() {
			t.Errorf("expected no data")
		}

		defer func() {
			if recover() == nil {
				t.Errorf("expected panic when calling Data() on loading resource")
			}
		}()
		r.Data()
	})

	t.Run("ResourceError", func(t *testing.T) {
		err := errors.New("fail")
		r := ResourceError[string](err)
		if r.IsLoading() {
			t.Errorf("expected not loading")
		}
		if r.Error() != err {
			t.Errorf("expected error %v, got %v", err, r.Error())
		}
		if r.HasData() {
			t.Errorf("expected no data")
		}

		defer func() {
			if recover() == nil {
				t.Errorf("expected panic when calling Data() on error resource")
			}
		}()
		r.Data()
	})
}

func TestMapData(t *testing.T) {
	t.Run("Success to Success", func(t *testing.T) {
		r := ResourceSuccess("test")
		mapped := MapData(r, func(s string) int { return len(s) })
		if mapped.Data() != 4 {
			t.Errorf("expected 4, got %d", mapped.Data())
		}
	})
}

func TestFlatMapData(t *testing.T) {
	t.Run("Success to Success", func(t *testing.T) {
		r := ResourceSuccess("test")
		mapped := FlatMapData(r, func(s string) Resource[int] {
			return ResourceSuccess(len(s))
		})
		if !mapped.HasData() {
			t.Errorf("expected data")
		}
		if mapped.Data() != 4 {
			t.Errorf("expected 4, got %d", mapped.Data())
		}
	})

	t.Run("Success to Loading", func(t *testing.T) {
		r := ResourceSuccess("test")
		mapped := FlatMapData(r, func(s string) Resource[int] {
			return ResourceLoading[int]()
		})

		if !mapped.IsLoading() {
			t.Errorf("expected loading")
		}
		if mapped.HasData() {
			t.Errorf("expected no data")
		}

		defer func() {
			if recover() == nil {
				t.Errorf("Should panic because it is loading, not because of bug")
			}
		}()
		mapped.Data()
	})
}

func TestMapDataCaching(t *testing.T) {
	r := ResourceSuccess("test")
	count := 0
	mapped := MapData(r, func(s string) int {
		count++
		return len(s)
	})

	// First call triggers mapping
	if got := mapped.Data(); got != 4 {
		t.Errorf("expected 4, got %d", got)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Second call should return cached value without incrementing count
	if got := mapped.Data(); got != 4 {
		t.Errorf("expected 4, got %d", got)
	}
	if count != 1 {
		t.Errorf("Mapping function called more than once, count: %d", count)
	}
}

func TestMapDataCachingWithPointer(t *testing.T) {
	r := ResourceSuccess("test")
	count := 0
	mapped := MapData(r, func(s string) *int {
		count++
		val := len(s)
		return &val
	})

	// First call triggers mapping
	res1 := mapped.Data()
	if *res1 != 4 {
		t.Errorf("expected 4, got %d", *res1)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Second call should return cached value (same pointer) without incrementing count
	res2 := mapped.Data()
	if *res2 != 4 {
		t.Errorf("expected 4, got %d", *res2)
	}
	if res1 != res2 {
		t.Errorf("Should return the exact same pointer")
	}
	if count != 1 {
		t.Errorf("Mapping function called more than once, count: %d", count)
	}
}

func TestMatch(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := ResourceSuccess("success")
		result := Match(r,
			func() string { return "loading" },
			func(err error) string { return "error" },
			func(data string) string { return data },
		)
		if result != "success" {
			t.Errorf("expected 'success', got %s", result)
		}
	})

	t.Run("Loading", func(t *testing.T) {
		r := ResourceLoading[string]()
		result := Match(r,
			func() string { return "loading" },
			func(err error) string { return "error" },
			func(data string) string { return data },
		)
		if result != "loading" {
			t.Errorf("expected 'loading', got %s", result)
		}
	})

	t.Run("Error", func(t *testing.T) {
		r := ResourceError[string](errors.New("fail"))
		result := Match(r,
			func() string { return "loading" },
			func(err error) string { return "error: " + err.Error() },
			func(data string) string { return data },
		)
		if result != "error: fail" {
			t.Errorf("expected 'error: fail', got %s", result)
		}
	})
}

func TestMatchPartial(t *testing.T) {
	t.Run("Loading with no options returns None", func(t *testing.T) {
		r := ResourceLoading[string]()
		result := MatchPartial[string, string](r)
		if result.IsSome() {
			t.Errorf("expected None, got Some(%v)", result.OrElse(""))
		}
	})

	t.Run("Success with no options returns None", func(t *testing.T) {
		r := ResourceSuccess("data")
		result := MatchPartial[string, string](r)
		if result.IsSome() {
			t.Errorf("expected None, got Some(%v)", result.OrElse(""))
		}
	})

	t.Run("Error with no options returns None", func(t *testing.T) {
		r := ResourceError[string](errors.New("fail"))
		result := MatchPartial[string, string](r)
		if result.IsSome() {
			t.Errorf("expected None, got Some(%v)", result.OrElse(""))
		}
	})

	t.Run("Loading with OnLoading handler", func(t *testing.T) {
		r := ResourceLoading[string]()
		result := MatchPartial[string, string](r,
			WithOnLoading[string, string](func() string { return "loading" }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "loading" {
			t.Errorf("expected 'loading', got %s", result.OrElse(""))
		}
	})

	t.Run("Success with OnData handler", func(t *testing.T) {
		r := ResourceSuccess("hello")
		result := MatchPartial[string, string](r,
			WithOnData[string, string](func(data string) string { return "got: " + data }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "got: hello" {
			t.Errorf("expected 'got: hello', got %s", result.OrElse(""))
		}
	})

	t.Run("Error with OnError handler", func(t *testing.T) {
		r := ResourceError[string](errors.New("fail"))
		result := MatchPartial[string, string](r,
			WithOnError[string, string](func(err error) string { return "error: " + err.Error() }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "error: fail" {
			t.Errorf("expected 'error: fail', got %s", result.OrElse(""))
		}
	})

	t.Run("Loading uses default when no OnLoading", func(t *testing.T) {
		r := ResourceLoading[string]()
		result := MatchPartial[string, string](r,
			WithDefault[string, string](func() string { return "default" }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "default" {
			t.Errorf("expected 'default', got %s", result.OrElse(""))
		}
	})

	t.Run("Success uses default when no OnData", func(t *testing.T) {
		r := ResourceSuccess("data")
		result := MatchPartial[string, string](r,
			WithDefault[string, string](func() string { return "default" }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "default" {
			t.Errorf("expected 'default', got %s", result.OrElse(""))
		}
	})

	t.Run("Error uses default when no OnError", func(t *testing.T) {
		r := ResourceError[string](errors.New("fail"))
		result := MatchPartial[string, string](r,
			WithDefault[string, string](func() string { return "default" }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "default" {
			t.Errorf("expected 'default', got %s", result.OrElse(""))
		}
	})

	t.Run("OnLoading takes precedence over default", func(t *testing.T) {
		r := ResourceLoading[string]()
		result := MatchPartial[string, string](r,
			WithOnLoading[string, string](func() string { return "loading" }),
			WithDefault[string, string](func() string { return "default" }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "loading" {
			t.Errorf("expected 'loading', got %s", result.OrElse(""))
		}
	})

	t.Run("OnData takes precedence over default", func(t *testing.T) {
		r := ResourceSuccess("data")
		result := MatchPartial[string, string](r,
			WithOnData[string, string](func(d string) string { return d }),
			WithDefault[string, string](func() string { return "default" }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "data" {
			t.Errorf("expected 'data', got %s", result.OrElse(""))
		}
	})

	t.Run("OnError takes precedence over default", func(t *testing.T) {
		r := ResourceError[string](errors.New("fail"))
		result := MatchPartial[string, string](r,
			WithOnError[string, string](func(err error) string { return "err" }),
			WithDefault[string, string](func() string { return "default" }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "err" {
			t.Errorf("expected 'err', got %s", result.OrElse(""))
		}
	})

	t.Run("All handlers with Success", func(t *testing.T) {
		r := ResourceSuccess("success")
		result := MatchPartial[string, string](r,
			WithOnLoading[string, string](func() string { return "loading" }),
			WithOnError[string, string](func(err error) string { return "error" }),
			WithOnData[string, string](func(d string) string { return d }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "success" {
			t.Errorf("expected 'success', got %s", result.OrElse(""))
		}
	})

	t.Run("All handlers with Loading", func(t *testing.T) {
		r := ResourceLoading[string]()
		result := MatchPartial[string, string](r,
			WithOnLoading[string, string](func() string { return "loading" }),
			WithOnError[string, string](func(err error) string { return "error" }),
			WithOnData[string, string](func(d string) string { return d }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "loading" {
			t.Errorf("expected 'loading', got %s", result.OrElse(""))
		}
	})

	t.Run("All handlers with Error", func(t *testing.T) {
		r := ResourceError[string](errors.New("fail"))
		result := MatchPartial[string, string](r,
			WithOnLoading[string, string](func() string { return "loading" }),
			WithOnError[string, string](func(err error) string { return "error" }),
			WithOnData[string, string](func(d string) string { return d }),
		)
		if !result.IsSome() {
			t.Errorf("expected Some, got None")
		}
		if result.OrElse("") != "error" {
			t.Errorf("expected 'error', got %s", result.OrElse(""))
		}
	})
}
