package web

import (
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
			if !slices.ContainsFunc(moduleParams, func(p string) bool { return strings.Contains(path, "{"+p+"}") }) {
				t.Errorf("%s: the gate cannot find the event of %s", op.OperationID, path)
			}
		}
	}
	if n == 0 {
		t.Fatal("no module operations")
	}
}
