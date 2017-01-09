package main

func LCP(l []string) string {
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0]
	}

	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}

	return min
}
