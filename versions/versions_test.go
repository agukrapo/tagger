package versions

import (
	"reflect"
	"testing"
)

func TestCommit_Change(t *testing.T) {
	tests := []struct {
		name  string
		msg   string
		want1 Change
		want2 string
	}{
		{
			"Breaking",
			"chore!: drop support for Node 6",
			Breaking,
			"drop support for Node 6",
		},
		{
			"Breaking (scope)",
			"feat(api)!: ASD123: send an email to the customer when a product is shipped",
			Breaking,
			"ASD123: send an email to the customer when a product is shipped",
		},
		{
			"Feat",
			"feat: allow provided config object to extend other configs",
			Feat,
			"allow provided config object to extend other configs",
		},
		{
			"Feat (scope)",
			"feat(lang): add Polish language",
			Feat,
			"add Polish language",
		},
		{
			"Fix",
			"fix: qwerty:prevent racing of requests",
			Fix,
			"qwerty:prevent racing of requests",
		},
		{
			"Fix (scope)",
			"fix(lang): prevent racing of requests",
			Fix,
			"prevent racing of requests",
		},
		{
			"None",
			"docs: correct spelling of CHANGELOG",
			None,
			"correct spelling of CHANGELOG",
		},
		{
			"None (scope)",
			"docs(lang): update ref docs",
			None,
			"update ref docs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commit := &Commit{
				message: tt.msg,
			}
			if got1, got2 := commit.Change(); got1 != tt.want1 || got2 != tt.want2 {
				t.Errorf("Change() = %v %v, want %v %v", got1, got2, tt.want1, tt.want2)
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

func TestVersion_String(t *testing.T) {
	tests := []struct {
		major int
		minor int
		patch int
		want  string
	}{
		{
			want: "v0",
		},
		{
			major: 1,
			want:  "v1",
		},
		{
			major: 2,
			minor: 3,
			want:  "v2.3",
		},
		{
			major: 4,
			minor: 5,
			patch: 6,
			want:  "v4.5.6",
		},
		{
			major: 7,
			patch: 8,
			want:  "v7.0.8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			v := Version{
				major: tt.major,
				minor: tt.minor,
				patch: tt.patch,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
