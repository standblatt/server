package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/standblatt/server/internal/modules"
)

// TestModuleOperations checks that every operation of a module names a known module and has a path parameter the
// gate finds the event by.
func TestModuleOperations(t *testing.T) {
	s := &Server{}
	api := humago.New(http.NewServeMux(), s.apiConfig())
	s.registerAPI(api)
	n := 0
	for path, item := range api.OpenAPI().Paths {
		for _, op := range []*huma.Operation{item.Get, item.Post, item.Put, item.Patch, item.Delete} {
			if op == nil {
				continue
			}
			key := modules.Of(op.Metadata)
			if key == "" {
				continue
			}
			n++
			if _, ok := modules.Get(key); !ok {
				t.Errorf("%s: unknown module %q", op.OperationID, key)
			}
			if m, _ := modules.Get(key); m.Level == modules.Event && !slices.ContainsFunc(moduleParams, func(p string) bool { return strings.Contains(path, "{"+p+"}") }) {
				t.Errorf("%s: the gate cannot find the event of %s", op.OperationID, path)
			}
		}
	}
	if n == 0 {
		t.Fatal("no module operations")
	}
}

// TestModuleEnums keeps the enum tags of module lists in the API in step with modules.All.
func TestModuleEnums(t *testing.T) {
	s := &Server{}
	api := humago.New(http.NewServeMux(), s.apiConfig())
	s.registerAPI(api)
	spec, err := json.Marshal(api.OpenAPI())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{"event": modules.Keys(modules.Event), "club": modules.Keys(modules.Club)}
	found := map[string]int{}
	var walk func(v any)
	walk = func(v any) {
		switch v := v.(type) {
		case map[string]any:
			if list, ok := v["enum"].([]any); ok && len(list) > 0 {
				for level, keys := range want {
					if list[0] == keys[0] {
						found[level]++
						if got := fmt.Sprint(list); got != fmt.Sprint(keys) {
							t.Errorf("enum %s, want the %s modules %v", got, level, keys)
						}
					}
				}
			}
			for _, x := range v {
				walk(x)
			}
		case []any:
			for _, x := range v {
				walk(x)
			}
		}
	}
	var doc any
	json.Unmarshal(spec, &doc)
	walk(doc)
	if found["event"] == 0 || found["club"] == 0 {
		t.Fatalf("no module enums found: %v", found)
	}
}
