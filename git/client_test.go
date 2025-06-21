package git

import "testing"

func TestCommit_Change(t *testing.T) {
	tests := []struct {
		name   string
		commit Commit
		want   Change
	}{
		{
			"Breaking",
			"chore!: drop support for Node 6",
			Breaking,
		},
		{
			"Breaking (scope)",
			"feat(api)!: send an email to the customer when a product is shipped",
			Breaking,
		},

		{
			"Feat",
			"feat: allow provided config object to extend other configs",
			Feat,
		},
		{
			"Feat (scope)",
			"feat(lang): add Polish language",
			Feat,
		},
		{
			"Fix",
			"fix: prevent racing of requests",
			Fix,
		},
		{
			"Fix (scope)",
			"fix(lang): prevent racing of requests",
			Fix,
		},
		{
			"None",
			"docs: correct spelling of CHANGELOG",
			None,
		},
		{
			"None (scope)",
			"docs(lang): update ref docs",
			None,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.commit.Change(); got != tt.want {
				t.Errorf("Change() = %v, want %v", got, tt.want)
			}
		})
	}
}
