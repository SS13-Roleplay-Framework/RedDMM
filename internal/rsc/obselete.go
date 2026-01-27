package rsc

import (
	_ "embed"
)

var (
	//go:embed obselete/obselete.dm
	ObsoleteDM string

	//go:embed obselete/obselete.dmi
	ObsoleteDMI []byte
)
