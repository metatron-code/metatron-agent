package app

import "errors"

var errorRunSnapd = errors.New("application must run only in ubuntu core or snapd")
