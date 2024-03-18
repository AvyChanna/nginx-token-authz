package auther

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AvyChanna/nginx-token-authz/lib/set"
	"github.com/AvyChanna/nginx-token-authz/lib/tabletest"
)

func TestAuther(t *testing.T) {
	type fields struct {
		data   map[string]set.StrSet
		admins set.StrSet
	}
	type arg struct {
		uid   string
		claim string
		want  bool
	}

	tests := tabletest.T[fields, arg]{
		{
			"Admin user",
			fields{
				map[string]set.StrSet{"admin": set.New("perm1")},
				set.New("admin"),
			},
			[]arg{
				{"admin", "perm1", true},
			},
		},
		{
			"Non-admin user with correct claim",
			fields{map[string]set.StrSet{"user": set.New("perm2")}, set.New[string]()},
			[]arg{
				{"user", "perm2", true},
			},
		},
		{
			"Non-admin user with incorrect claim",
			fields{map[string]set.StrSet{"user": set.New("perm3")}, set.New[string]()},
			[]arg{
				{"user", "perm4", false},
			},
		},
		{
			"User not found in data",
			fields{map[string]set.StrSet{"admin": set.New("perm5")}, set.New("admin")},
			[]arg{
				{"non-existing user", "perm6", false},
			},
		},
	}

	tests.Run(t, func(t *testing.T, f fields, args []arg) {
		testAuther := New(f.data, f.admins)

		for _, arg := range args {
			assert.Equal(t, arg.want, testAuther.IsValidClaim(arg.uid, arg.claim))
		}
	})
}
