package web_test

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/standblatt/server/internal/web"
)

// viaTokens makes logged-in test clients call every path an API token reaches with a token of the user instead of
// the session; TestMain runs all tests once with sessions and once with tokens.
var viaTokens bool

// called collects the operations that succeeded, per mode.
var called = struct {
	sync.Mutex
	session, token map[string]bool
}{session: map[string]bool{}, token: map[string]bool{}}

// TestMain runs the tests with sessions and again with API tokens, then checks on a full run that every operation
// succeeded at least once and every operation a token reaches succeeded with a token.
func TestMain(m *testing.M) {
	web.ObserveOperationsForTest(func(op string, status int, token bool) {
		if status >= 400 {
			return
		}
		called.Lock()
		defer called.Unlock()
		if token {
			called.token[op] = true
		} else {
			called.session[op] = true
		}
	})
	code := m.Run()
	viaTokens = true
	if c := m.Run(); code == 0 {
		code = c
	}
	if code == 0 && flag.Lookup("test.run").Value.String() == "" && flag.Lookup("test.skip").Value.String() == "" {
		code = checkCalled()
	}
	os.Exit(code)
}

// checkCalled reports operations no test called successfully, and token-reachable ones no token called.
func checkCalled() int {
	code := 0
	for _, op := range web.OperationsForTest() {
		switch {
		case !called.session[op.ID] && !called.token[op.ID]:
			fmt.Printf("never called successfully: %s %s (%s)\n", op.Method, op.Path, op.ID)
			code = 1
		case web.TokenMayForTest(op.Method, op.Path) && !called.token[op.ID]:
			fmt.Printf("never called successfully with an API token: %s %s (%s)\n", op.Method, op.Path, op.ID)
			code = 1
		}
	}
	return code
}
