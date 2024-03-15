package rbac

import (
	"github.com/AvyChanna/nginx-token-authz/lib/set"
)

func validateUid(uid string) error {
	if len(uid) == 0 {
		return ErrUidInvalid
	}

	return nil
}

func validateGid(gid string) error {
	if len(gid) == 0 {
		return ErrGidInvalid
	}

	return nil
}

func validatePmap(perm map[string]bool) error {
	for k := range perm {
		if !strRegex.MatchString(k) {
			return ErrPermInvalid
		}
	}

	return nil
}

func validatePset(pset set.StrSet) error {
	for p := range pset {
		if !strRegex.MatchString(p) {
			return ErrPermInvalid
		}
	}

	return nil
}

func (d *Config) validateData() error {
	for gid, group := range d.Groups {
		if err := validateGid(gid); err != nil {
			return err
		}
		if err := validatePmap(group.Pmap); err != nil {
			return err
		}
	}

	for uid, user := range d.Users {
		if err := validateUid(uid); err != nil {
			return err
		}
		if user.Groups != nil {
			for gid := range *user.Groups {
				if _, ok := d.Groups[gid]; !ok {
					return ErrGidInvalid
				}
			}
		}

		if err := validatePmap(user.Pmap); err != nil {
			return err
		}
	}
	if d.AllowAllSet != nil {
		if err := validatePset(*d.AllowAllSet); err != nil {
			return err
		}
	}

	return nil
}
