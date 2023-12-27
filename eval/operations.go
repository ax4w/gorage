package eval

func less(l, r *token) *token {
	checkIfCompatible(l, r)
	if l.tokenType == tokenTypeDate && compareDates(string(l.value), string(r.value)) == 1 {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	} else if convertBytesToFloat(l.value) < convertBytesToFloat(r.value) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func greater(l, r *token) *token {
	checkIfCompatible(l, r)
	if l.tokenType == tokenTypeDate && compareDates(string(l.value), string(r.value)) == -1 {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	} else if convertBytesToFloat(l.value) > convertBytesToFloat(r.value) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func greaterThan(l, r *token) *token {
	checkIfCompatible(l, r)
	if l.tokenType == tokenTypeDate && compareDates(string(l.value), string(r.value)) <= 0 {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	} else if convertBytesToFloat(l.value) >= convertBytesToFloat(r.value) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func lessThan(l, r *token) *token {
	checkIfCompatible(l, r)
	if l.tokenType == tokenTypeDate && compareDates(string(l.value), string(r.value)) >= 0 {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	} else if convertBytesToFloat(l.value) <= convertBytesToFloat(r.value) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func equal(l, r *token) *token {
	checkIfCompatible(l, r)
	if l.tokenType == tokenTypeDate && compareDates(string(l.value), string(r.value)) == 0 {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	} else if compareByteArray(l.value, r.value) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func notEqual(l, r *token) *token {
	checkIfCompatible(l, r)
	if l.tokenType == tokenTypeDate && compareDates(string(l.value), string(r.value)) != 0 {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	} else if !compareByteArray(l.value, r.value) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func and(l, r *token) *token {
	if !(l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
		panic("&& expects both sides to be a boolean")
	}
	if compareByteArray(l.value, []byte("t")) && compareByteArray(r.value, []byte("t")) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func or(l, r *token) *token {
	if !(l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
		panic("&& expects both sides to be a boolean")
	}
	if compareByteArray(l.value, []byte("t")) || compareByteArray(r.value, []byte("t")) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func nand(l, r *token) *token {
	if !(l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
		panic("&& expects both sides to be a boolean")
	}
	if !(compareByteArray(l.value, []byte("t")) && compareByteArray(r.value, []byte("t"))) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}

func nor(l, r *token) *token {
	if !(l.tokenType == tokenTypeBoolean && r.tokenType == tokenTypeBoolean) {
		panic("&& expects both sides to be a boolean")
	}
	if !(compareByteArray(l.value, []byte("t")) || compareByteArray(r.value, []byte("t"))) {
		return &token{value: []byte("t"), tokenType: tokenTypeBoolean}
	}
	return &token{value: []byte("f"), tokenType: tokenTypeBoolean}
}
