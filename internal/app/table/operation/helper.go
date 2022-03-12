package operation

// mergeIndexes объединяет два слайса, сохраняя их порядок и исключая дубликаты
func mergeIndexes(first []int, second []int) []int {
	var ret []int

	if len(first) > len(second) {
		ret = make([]int, 0, len(first))
	} else {
		ret = make([]int, 0, len(second))
	}

	i, j := 0, 0
	lastValue := -1

	for {
		// Проверяем, не выходим ли мы за пределы слайсов
		if i == len(first) || j == len(second) {
			break
		}
		// Игнорируем дубликаты
		if first[i] == lastValue {
			i++
			continue
		}
		if second[j] == lastValue {
			j++
			continue
		}
		if first[i] < second[j] {
			ret = append(ret, first[i])
			lastValue = first[i]
			i++
			continue
		}
		ret = append(ret, second[j])
		lastValue = second[j]
		j++
	}

	// копируем оставшиеся элементы из первого слайса при их наличии
	if i < len(first) {
		for ; i < len(first); i++ {
			if lastValue != first[i] {
				ret = append(ret, first[i])
			}
		}
	}

	// копируем оставшиеся элементы из второго слайса при их наличии
	if j < len(second) {
		for ; j < len(second); j++ {
			if lastValue != second[j] {
				ret = append(ret, second[j])
			}
		}
	}

	return ret
}
