package banner

import (
	"os"
	"strings"
	"testing"
)

func TestRender_EnvironmentSuffixGating(t *testing.T) {
	tests := []struct {
		name          string
		opts          Options
		envVarValue   string
		wantSuffix    string
		notWantSuffix string
	}{
		{
			name: "explicit Environment => no suffix",
			opts: Options{
				ServiceName: "test-svc",
				Environment: "production",
			},
			wantSuffix:    "",
			notWantSuffix: "(production)",
		},
		{
			name: "empty Environment + env var set => suffix shown",
			opts: Options{
				ServiceName: "test-svc",
			},
			envVarValue: "staging",
			wantSuffix:  "(staging)",
		},
		{
			name: "empty Environment + env var not set => no suffix",
			opts: Options{
				ServiceName: "test-svc",
			},
			wantSuffix:    "",
			notWantSuffix: "(dev)", // current default that should be gone
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVarValue != "" {
				os.Setenv(EnvVarEnvironment, tt.envVarValue)
				defer os.Unsetenv(EnvVarEnvironment)
			} else {
				os.Unsetenv(EnvVarEnvironment)
			}

			out := Render(tt.opts, BuildInfo{})

			// The suffix appears in the header line: :: goStartyUpy :: (suffix)
			headerPart := ":: goStartyUpy ::"
			if tt.wantSuffix != "" {
				want := headerPart + " " + tt.wantSuffix
				if !strings.Contains(out, want) {
					t.Errorf("output missing expected suffix %q\nOutput:\n%s", want, out)
				}
			} else if tt.notWantSuffix != "" {
				notWant := headerPart + " " + tt.notWantSuffix
				if strings.Contains(out, notWant) {
					t.Errorf("output contains unexpected suffix %q\nOutput:\n%s", notWant, out)
				}
			}
		})
	}
}
