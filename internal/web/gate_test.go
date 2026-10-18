package web

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/standblatt/server/internal/modules"
	"github.com/standblatt/server/internal/permissions"
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

// TestModuleEnums keeps the enum tags of module and permission lists in the API in step with modules.All and permissions.All.
func TestModuleEnums(t *testing.T) {
	s := &Server{}
	api := humago.New(http.NewServeMux(), s.apiConfig())
	s.registerAPI(api)
	spec, err := json.Marshal(api.OpenAPI())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{"event": modules.Keys(modules.Event), "club": modules.Keys(modules.Club), "permission": permissions.Keys(), "scope": permissions.Scopes()}
	found := map[string]int{}
	var walk func(v any)
	walk = func(v any) {
		switch v := v.(type) {
		case map[string]any:
			if list, ok := v["enum"].([]any); ok && len(list) > 0 {
				matched := ""
				for level, keys := range want {
					if list[0] == keys[0] && (matched == "" || fmt.Sprint(list) == fmt.Sprint(keys)) {
						matched = level
					}
				}
				if matched != "" {
					found[matched]++
					if got, keys := fmt.Sprint(list), fmt.Sprint(want[matched]); got != keys {
						t.Errorf("enum %s, want the %s list %s", got, matched, keys)
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
	if found["event"] == 0 || found["club"] == 0 || found["permission"] == 0 || found["scope"] == 0 {
		t.Fatalf("no module enums found: %v", found)
	}
}

// TestPermissionsUsed checks that every permission is required by at least one operation.
func TestPermissionsUsed(t *testing.T) {
	names := map[string]string{}
	f, err := parser.ParseFile(token.NewFileSet(), "../permissions/permissions.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok || g.Tok != token.CONST {
			continue
		}
		for _, spec := range g.Specs {
			v := spec.(*ast.ValueSpec)
			if lit, ok := v.Values[0].(*ast.BasicLit); ok && permissions.Known(strings.Trim(lit.Value, `"`)) {
				names[strings.Trim(lit.Value, `"`)] = v.Names[0].Name
			}
		}
	}
	var src strings.Builder
	files, _ := filepath.Glob("*.go")
	for _, name := range files {
		if !strings.HasSuffix(name, "_test.go") {
			b, _ := os.ReadFile(name)
			src.Write(b)
		}
	}
	for _, p := range permissions.All {
		if names[p.Key] == "" || !strings.Contains(src.String(), "permissions."+names[p.Key]) {
			t.Errorf("no operation requires %s", p.Key)
		}
	}
}
