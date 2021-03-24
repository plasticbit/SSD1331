package OLED

func maxint(nums ...int) (res int) {
	for _, val := range nums {
		if val > res {
			res = val
		}
	}

	return
}

func minint(nums ...int) (res int) {
	for _, val := range nums {
		if val < res {
			res = val
		}
	}

	return
}

func maxbyte(nums ...byte) (res byte) {
	for _, val := range nums {
		if val > res {
			res = val
		}
	}

	return
}

func minbyte(nums ...byte) (res byte) {
	for _, val := range nums {
		if val < res {
			res = val
		}
	}

	return
}
