package drops

type Flag uint

func SetFlag(set uint, flag Flag) uint {
	return set ^ uint(flag)
}

func UnsetFlag(set uint, flag Flag) uint {
	return set &^ uint(flag)
}

func IsFlag(set uint, flag Flag) bool {
	if set&uint(flag) > 0 {
		return true
	}
	return false
}
