package main

import (
	"fmt"
)

var subjects []string = []string{
	"数 学",
	"语 文",
	"英 语",
}

var scores [][]int = [][]int{
	{88, 88, 90},
	{66, 99, 94},
	{75, 84, 98},
	{93, 77, 66},
}

func main() {
	scores = append(scores, []int{98, 97, 99})

	mathAvg := avgScores(0)
	chineseAvg := avgScores(1)
	englishAvg := avgScores(2)

	fmt.Printf("%10s %7s %7s\n", subjects[0], subjects[1], subjects[2])
	for i, v := range scores {
		// fmt.Printf("%5d %7d %7d\n", scores[i][0], scores[i][1], scores[i][2])
		for j := range v {
			fmt.Printf("%10d", scores[i][j])
		}
		fmt.Println()
	}
	fmt.Printf("%10d %9d %9d\n", mathAvg, chineseAvg, englishAvg)
}

func avgScores(num int) int {
	All := 0
	for i := 0; i < len(scores); i++ {
		All += scores[i][num]
	}
	Avg := All / len(scores)
	return Avg
}
