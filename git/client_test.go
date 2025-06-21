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
			"b2504f5 chore!: drop support for Node 6",
			Breaking,
		},
		{
			"Breaking (scope)",
			"4d0442b (origin/fixes, fixes) feat(api)!: send an email to the customer when a product is shipped",
			Breaking,
		},

		{
			"Feat",
			"db5a126 (HEAD -> feat/ws) feat: allow provided config object to extend other configs",
			Feat,
		},
		{
			"Feat (scope)",
			"b2504f5 feat(lang): add Polish language",
			Feat,
		},
		{
			"Fix",
			"6f99d20 fix: prevent racing of requests",
			Fix,
		},
		{
			"Fix (scope)",
			"da17a5d (origin/feat/client_switch, feat/client_switch) fix(lang): prevent racing of requests",
			Fix,
		},
		{
			"None",
			"6d52e21 (origin/env/prod) docs: correct spelling of CHANGELOG",
			None,
		},
		{
			"None (scope)",
			"e4ae48f docs(lang): update ref docs",
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
