package auther

import (
	"encoding/json"

	"github.com/AvyChanna/nginx-token-authz/internal/set"
)

type Auther struct {
	data   map[string]set.StrSet // uid:{permName}
	admins set.StrSet            // {uid}
}

func New(uidPermMap map[string]set.StrSet, admins set.StrSet) *Auther {
	return &Auther{
		data:   uidPermMap,
		admins: admins,
	}
}

func (a Auther) IsValidClaim(uid, claim string) bool {
	if a.admins.Contains(uid) {
		return true
	}

	user, found := a.data[uid]
	if !found {
		return false
	}

	return user.Contains(claim)
}

func (a Auther) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Data   map[string]set.StrSet `json:"data"`
		Admins set.StrSet            `json:"admins"`
	}{
		Data:   a.data,
		Admins: a.admins,
	})
}
