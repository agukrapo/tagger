package versions

import (
	"reflect"
	"testing"
)

func TestCommit_Change(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want Change
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
			commit := &Commit{
				message: tt.msg,
			}
			if got := commit.change(); got != tt.want {
				t.Errorf("Change() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTag_asVersion(t *testing.T) {
	tests := []struct {
		tag     Tag
		version Version
		error   string
	}{
		{
			tag:   "latest",
			error: `invalid tag "latest"`,
		},
		{
			tag:   "valpha",
			error: `invalid tag "valpha"`,
		},
		{
			tag:   "v",
			error: `invalid tag "v"`,
		},
		{
			tag:     "v1",
			version: Version{1, 0, 0},
		},
		{
			tag:     "v1.2",
			version: Version{1, 2, 0},
		},
		{
			tag:     "v1.2.3",
			version: Version{1, 2, 3},
		},
		{
			tag:   "v1.2.3.4",
			error: `invalid tag "v1.2.3.4"`,
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.tag), func(t *testing.T) {
			got, err := tt.tag.asVersion()
			if errNotEqual(tt.error, err) {
				t.Errorf("asVersion() err = %v, error %v", err, tt.error)
				return
			}
			if !reflect.DeepEqual(got, tt.version) {
				t.Errorf("asVersion() got = %v, want %v", got, tt.version)
			}
		})
	}
}

func errNotEqual(errStr string, err error) bool {
	if errStr == "" {
		return err != nil
	} else {
		return err == nil || err.Error() != errStr
	}
}
