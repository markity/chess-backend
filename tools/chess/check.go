package chess

func CheckChessPostsionVaild(x rune, y int) bool {
	if x != 'a' && x != 'b' && x != 'c' && x != 'd' && x != 'e' && x != 'f' && x != 'g' && x != 'h' {
		return false
	}

	if y < 1 || x > 8 {
		return false
	}

	return true
}
