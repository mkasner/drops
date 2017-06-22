package drops

type Flag uint

func SetFlag(set Flag, flag Flag) Flag {
	return set ^ flag
}

func UnsetFlag(set Flag, flag Flag) Flag {
	return set &^ flag
}

func IsFlag(set Flag, flag Flag) bool {
	if set&flag > 0 {
		return true
	}
	return false
}
