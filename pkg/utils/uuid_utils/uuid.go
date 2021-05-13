package uuid_utils

import "github.com/rs/xid"

func New() string {
	return xid.New().String()
}
