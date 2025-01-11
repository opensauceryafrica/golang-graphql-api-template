// a database interface for moving parts while ensuring not to break things
package migration

import (
	auth "cendit.io/auth/schema"
	garage "cendit.io/garage/schema"
)

var Tables = []interface{}{
	&auth.Activity{},
	&auth.Factory{},
	&auth.Transaction{},
	&auth.Card{},
	&auth.KYC{},
	&auth.User{},
	&auth.Address{},
	&auth.SecuritySetting{},
	&auth.Wallet{},
	&garage.PlatformConfig{},
	&garage.Rate{},
	&garage.Template{},
}
