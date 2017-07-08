package drops

type Flag uint64

// SetFlag sets flag on a new flagset
// Use it when you need only one flag set and others discarded
func SetFlag(flag Flag) Flag {
	var newSet Flag
	return newSet ^ flag
}

// AddFlag adds flag to existing flagset.
// Useful when multiple flags want to be set
func AddFlag(set Flag, flag ...Flag) Flag {
	for _, f := range flag {
		set = set ^ f
	}
	return set
}

// UnsetFlag unsets specified flag from existing flag set
func UnsetFlag(set Flag, flag Flag) Flag {
	return set &^ flag
}

// IsFlag checks if specific flag is set in provided flagset
func IsFlag(set Flag, flag Flag) bool {
	if set&flag > 0 {
		return true
	}
	return false
}
