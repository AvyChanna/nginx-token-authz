package parser

import (
	"github.com/AvyChanna/nginx-token-authz/internal/rbac/auther"
	"github.com/AvyChanna/nginx-token-authz/internal/set"
)

func New() *Config {
	return &Config{
		Users:       map[string]User{},
		Groups:      map[string]Group{},
		AllowAllSet: &set.StrSet{},
	}
}

type User struct {
	Admin  bool            `yaml:"admin"`
	Groups *set.StrSet     `yaml:"groups"`
	Pmap   map[string]bool `yaml:"perms"`
}

type Group struct {
	Admin bool            `yaml:"admin"`
	Pmap  map[string]bool `yaml:"perms"`
}

type Config struct {
	Users       map[string]User  `yaml:"users"`
	Groups      map[string]Group `yaml:"groups"`
	AllowAllSet *set.StrSet      `yaml:"globalAllows"`
}

func (d *Config) getAllowedPerms(user User) set.StrSet {
	// this is the final whitelist for user. missing perms are blocked
	whiteList := set.StrSet{}

	userBlackList := set.StrSet{}

	for perm, val := range user.Pmap {
		if val {
			whiteList.Add(perm)
		} else {
			userBlackList.Add(perm)
		}
	}

	groupWhiteList := set.StrSet{}
	groupBlackList := set.StrSet{}

	for gid := range *user.Groups {
		group := d.Groups[gid]

		for perm, val := range group.Pmap {
			if val {
				groupWhiteList.Add(perm)
			} else {
				groupBlackList.Add(perm)
			}
		}
	}

	// group blocking takes precedence over group allowing
	// remove all blocked perms from the group whitelist
	for bPerm := range groupBlackList {
		groupWhiteList.Remove(bPerm)
	}

	// only add to user whitelist if it is not blocked on user level
	for wPerm := range groupWhiteList {
		if !userBlackList.Contains(wPerm) {
			whiteList.Add(wPerm)
		}
	}

	// add global whitelist rules
	// even user/group level blocking won't block these
	for dPerm := range *d.AllowAllSet {
		whiteList.Add(dPerm)
	}

	return whiteList
}

func (d *Config) isAdmin(user User) bool {
	if user.Admin {
		return true
	}

	for gid := range *user.Groups {
		group := d.Groups[gid]

		if group.Admin {
			return true
		}
	}

	return false
}

func (d *Config) Done() (*auther.Auther, error) {
	err := validateData(d)
	if err != nil {
		return nil, err
	}

	uidPerms := map[string]set.StrSet{}
	admins := set.StrSet{}

	for uid, user := range d.Users {
		if d.isAdmin(user) {
			admins.Add(uid)
			continue
		}

		uidPerms[uid] = d.getAllowedPerms(user)
	}

	return auther.New(uidPerms, admins), nil
}
