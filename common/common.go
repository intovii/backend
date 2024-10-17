package common

import "errors"

var StatusInvalidParams = 0
var ErrInvalidParams = errors.New("invalid params").Error()
var StatusGetInfo = 1
var ErrGetInfo = errors.New("can not get info").Error()
