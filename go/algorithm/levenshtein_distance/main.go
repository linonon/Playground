package main

func main() {

}

// LevenshteinDistanceV1 version 1
func LevenshteinDistanceV1(s1, s2 string) int {
	if s1 == "" {
		return len(s2)
	}
	if s2 == "" {
		return len(s1)
	}

	len1, len2 := len(s1), len(s2)
	distance := make([][]int, len1+1)
	for i := range distance {
		distance[i] = make([]int, len2+1)
	}
	for i := 0; i < len1; i++ {
		distance[i][0] = i
	}
	for j := 0; j < len2; j++ {
		distance[0][j] = j
	}

	min := func(i, j int) int {
		if i < j {
			return i
		}
		return j
	}

	for i, ch1 := range s1 {
		for j, ch2 := range s2 {
			ins := distance[i+1][j] + 1
			del := distance[i][j+1] + 1
			sub := distance[i][j]
			if ch1 != ch2 {
				sub++
			}
			distance[i+1][j+1] = min(min(ins, del), sub)
		}
	}

	return distance[len1-1][len2-1]
}

// LevenshteinDistanceV2 version 2
func LevenshteinDistanceV2(s1, s2 string) int {
	if s1 == "" {
		return len(s2)
	}
	if s2 == "" {
		return len(s1)
	}

	var (
		len2     = len(s2)
		distance []int
	)

	distance = make([]int, len2+1)
	for i := range distance {
		distance[i] = i
	}

	for i, ch1 := range s1 {
		sub := i
		distance[0] = sub + 1
		for j, ch2 := range s2 {
			if ch1 != ch2 {
				sub++
			}
			dist := min(
				min(
					distance[j],   // insert
					distance[j+1], // delete
				)+1,
				sub,
			)
			sub = distance[j+1]
			distance[j+1] = dist
		}
	}

	return distance[len(distance)-1]
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
