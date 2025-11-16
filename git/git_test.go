package git

import (
	"reflect"
	"testing"

	"github.com/agukrapo/tagger/versions"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		line string
		want *versions.Commit
	}{
		{
			line: "7b4ca00 (origin/main, origin/HEAD) ci: tagger version fix",
			want: versions.NewCommit("7b4ca00", "ci: tagger version fix"),
		},
		{
			line: "b0e9838 fix: debug info on error",
			want: versions.NewCommit("b0e9838", "fix: debug info on error"),
		},
		{
			line: "7b746b5 (tag: v0.4.2, tag: v0.4.1) fix: show debug info on push",
			want: versions.NewCommit("7b746b5", "fix: show debug info on push"),
		},
		{
			line: `185e42e (HEAD -> main) Revert "fix: tree tag type"`,
			want: versions.NewCommit("185e42e", `Revert "fix: tree tag type"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got, ok := parse(tt.line)
			if !ok {
				t.Error("parse() ok = false")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
