package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

func getIndredientText(testContent string) string {
	lines := strings.Split(string(testContent), "\n")
	scores := make([]int, len(lines))

	for i, line := range lines {
		score := 0
		for _, word := range commonFoods {
			if strings.Contains(line, word) == true {
				// fmt.Println("commonFoods: " + word)
				score += 1
			}
		}
		for _, word := range badContextWords {
			if strings.Contains(line, word) == true {
				// fmt.Println("badContextWords: " + word)
				score -= 1
			}
		}
		for _, word := range goodContextWords {
			if strings.Contains(line, word) == true {
				// fmt.Println("gooContextWords: " + word)
				score += 1
			}
		}
		//fmt.Println(i, score, line)
		scores[i] = score
	}

	minScore := 0
	for _, score := range scores {
		if score < minScore {
			minScore = score
		}
	}

	for i, score := range scores {
		scores[i] = score - minScore
	}

	gaussianScores := make([]int, len(scores))
	for i, _ := range scores {
		gaussianScore := 0
		if i > 0 {
			gaussianScore += scores[i-1] * 18
		}
		if i > 1 {
			gaussianScore += scores[i-2] * 12
		}
		if i > 2 {
			gaussianScore += scores[i-3] * 6
		}
		gaussianScore += scores[i] * 20
		if i < len(scores)-1 {
			gaussianScore += scores[i+1] * 18
		}
		if i < len(scores)-2 {
			gaussianScore += scores[i+2] * 12
		}
		if i < len(scores)-3 {
			gaussianScore += scores[i+3] * 6
		}
		gaussianScores[i] = gaussianScore
	}

	medianScore := getMedian(gaussianScores)
	maxScore := 0
	maxScoreI := 0
	for i, score := range gaussianScores {
		gaussianScores[i] = score - medianScore
		if gaussianScores[i] > maxScore {
			maxScore = gaussianScores[i]
			maxScoreI = i
		}
	}
	minResidual := 1e8
	bestSigma := 0
	bestMu := 0
	for muAdjust := -2; muAdjust < 3; muAdjust++ {
		for testSigma := 4; testSigma < 20*4; testSigma++ {
			testGaussian := gaussianFunc(maxScoreI+muAdjust, testSigma, maxScore, gaussianScores)
			resid := residualFunc(testGaussian, gaussianScores)
			// fmt.Println(testSigma, resid)
			if resid < minResidual && resid > 0 {
				bestSigma = testSigma
				minResidual = resid
				bestMu = maxScoreI + muAdjust
			}
		}
	}
	fmt.Println(gaussianScores)
	fmt.Println(gaussianFunc(bestMu, bestSigma, maxScore, gaussianScores))
	fmt.Println(maxScoreI, bestSigma)

	text := ""
	for i, line := range lines {
		if i > (bestMu-int(float64(bestSigma)/1.5)) && i < (bestMu+int(float64(bestSigma)/1.5)) {
			text += line + " \n"
		}
	}

	return text
}

func getMedian(values []int) int {
	v := make([]int, len(values))
	for i, val := range values {
		v[i] = val
	}
	sort.Ints(v)
	return v[len(v)/2]
}

func residualFunc(one []int, two []int) float64 {
	if len(one) != len(two) {
		return -1.0
	}
	residual := float64(0)
	for i, _ := range one {
		if two[i] > 0 && one[i] > 0 {
			residual += math.Pow(float64(one[i])-float64(two[i]), 2)
		}
	}
	return residual
}

func gaussianFunc(mu2 int, sigma2 int, max2 int, values []int) []int {
	sigma := float64(sigma2)
	mu := float64(mu2)
	max := float64(max2)
	x := make([]float64, len(values))
	for i, _ := range values {
		x[i] = float64(i)
	}
	gaussian := make([]float64, len(x))
	maxVal := float64(0)
	for i, _ := range x {
		gaussian[i] = 1 / (sigma / 4 * math.Sqrt(2.0*3.1415926535)) * math.Exp(-1*math.Pow(x[i]-mu, 2)/(2*math.Pow(sigma/4, 2)))
		if gaussian[i] > maxVal {
			maxVal = gaussian[i]
		}
	}

	intGaussian := make([]int, len(gaussian))
	for i, val := range gaussian {
		intGaussian[i] = int(max * val / maxVal)
	}
	return intGaussian
}
