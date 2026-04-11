package hegelgen_test

import (
	"strings"
	"testing"
	"time"

	"github.com/franchb/fptest-go/engine"
	fpthegel "github.com/franchb/fptest-go/hegel"
	"github.com/franchb/fptest-go/hegel/hegelgen"
)

func TestEmailsGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/emails", func(et engine.T) {
		email := hegelgen.Emails().Draw(et, "email")
		if !strings.Contains(email, "@") {
			t.Fatalf("expected email to contain @, got %q", email)
		}
	})
}

func TestURLsGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/urls", func(et engine.T) {
		url := hegelgen.URLs().Draw(et, "url")
		if !strings.HasPrefix(url, "http") {
			t.Fatalf("expected URL to start with http, got %q", url)
		}
	})
}

func TestDatesGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/dates", func(et engine.T) {
		d := hegelgen.Dates().Draw(et, "date")
		var zero time.Time
		if d == zero {
			t.Fatal("expected non-zero date")
		}
	})
}

func TestIntegersGenerator(t *testing.T) {
	runner := fpthegel.HegelRunner{}
	runner.MakeCheck(t, "test/integers", func(et engine.T) {
		n := hegelgen.Integers(1, 100).Draw(et, "n")
		if n < 1 || n > 100 {
			t.Fatalf("expected value in [1,100], got %d", n)
		}
	})
}
