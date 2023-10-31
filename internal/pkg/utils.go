package utils

// Ptr returns a point of e (type E) using generics.
func Ptr[E any](e E) *E {
	return &e
}

// GetBytesBetween gets the bytes between the first matched 'open' and its corresponding 'close'.
func GetBytesBetween(dat []byte, open byte, close byte, sx, se int) (b []byte) {
	loc := FindOpenAndCloseLocations(dat, open, close, sx, se)
	if loc != nil {
		s := dat[loc[0]:loc[1]]
		b = make([]byte, len(s))
		copy(b, s)
	}
	return
}

// FindOpenAndCloseLocations find the first 'open' and its matching 'close' index within 'dat'
func FindOpenAndCloseLocations(dat []byte, open byte, close byte, sx, se int) []int {
	if len(dat) < 1 {
		return nil
	}
	c := 0
	for c != 0 || dat[se] != close {
		se++
		if se > len(dat)-1 {
			return nil
		}
		if dat[se] == open {
			c++
			if dat[sx] != open {
				sx = se
			}
		} else if dat[se] == close {
			c--
		}
	}
	return []int{sx, se + 1}
}
