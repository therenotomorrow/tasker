package config_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/therenotomorrow/tasker/internal/config"
)

func TestUnitPath(t *testing.T) {
	t.Parallel()

	type args struct {
		parts []string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "without parts", args: args{parts: make([]string, 0)}, want: "/"},
		{name: "with parts", args: args{parts: []string{"test", "data", "empty"}}, want: "/test/data/empty"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := config.Path(test.args.parts...)
			want := filepath.Clean(test.want)

			if !filepath.IsAbs(got) || !strings.Contains(got, want) {
				t.Errorf("Path() got = %v, want = %v", got, want)
			}
		})
	}
}
