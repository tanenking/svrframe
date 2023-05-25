package helper

func GetRandomBetweenI64(min int64, max int64) int64 {
	if min == max {
		return min
	}
	_min := min
	_max := max
	if _min > _max {
		_min = max
		_max = min
	}
	v := _max - _min
	if v <= 0 {
		return _min
	}
	v++

	return min + RD.Int63n(v)
}
func GetRandomBetween(min int32, max int32) int32 {
	v := GetRandomBetweenI64(int64(min), int64(max))
	return int32(v)
}

func GetProbResult(prob int) bool {
	//万分比
	g := RD.Intn(10000) + 1

	return g <= prob
}

// 二维数组 [[id,weight]]
func GetWeightFromJsonArray(weights [][]int) int {
	if len(weights) <= 0 {
		return 0
	}
	type wd_t struct {
		idx    int
		weight int
	}

	wdsum := 0
	wds := []wd_t{}
	for _, v := range weights {
		if len(v) < 2 {
			continue
		}
		wdsum += v[1]
		wds = append(wds, wd_t{
			idx:    v[0],
			weight: wdsum,
		})
	}
	r := GetRandomBetween(0, int32(wdsum)-1)

	for i := 0; i < len(wds); i++ {
		if r < int32(wds[i].weight) {
			return wds[i].idx
		}
	}
	return 0
}
func GetWeightFromMaps(weights map[int]int) int {
	if len(weights) <= 0 {
		return 0
	}
	type wd_t struct {
		idx    int
		weight int
	}

	wdsum := 0
	wds := []wd_t{}
	for k, v := range weights {
		if v > 0 {
			wdsum += v
			wds = append(wds, wd_t{
				idx:    k,
				weight: wdsum,
			})
		}
	}
	r := GetRandomBetween(0, int32(wdsum)-1)

	for i := 0; i < len(wds); i++ {
		if r < int32(wds[i].weight) {
			return wds[i].idx
		}
	}
	return 0
}
func GetWeightFromProbs(probs []int32) int {
	if len(probs) <= 0 {
		return 0
	}
	type wd_t struct {
		idx    int
		weight int32
	}

	wdsum := int32(0)
	wds := []wd_t{}
	for k, v := range probs {
		if v > 0 {
			wdsum += v
			wds = append(wds, wd_t{
				idx:    k,
				weight: wdsum,
			})
		}
	}
	r := GetRandomBetween(0, wdsum-1)

	for i := 0; i < len(wds); i++ {
		if r < wds[i].weight {
			return wds[i].idx
		}
	}
	return 0
}

func GetWeightFromProbsInt(probs []int) int {
	if len(probs) <= 0 {
		return 0
	}
	probs_32 := []int32{}
	for i := 0; i < len(probs); i++ {
		probs_32 = append(probs_32, int32(probs[i]))
	}
	return GetWeightFromProbs(probs_32)
}
