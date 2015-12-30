package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

var commonFoods []string
var goodContextWords []string
var badContextWords []string

// func gaussianFilter(data []int) (filtered []int) {

// }

func init() {
	loadData()
}

func main() {

	fmt.Println(commonFoods)
	fmt.Println(badContextWords)
	fmt.Println(goodContextWords)

	testContent := `Hi
Ingredients
2 cups all-purpose flour
1/4 cup sugar
1 tablespoon baking powder
1/2 teaspoon salt
1 1/2 cups milk
1 tablespoon plus 1 teaspoon vanilla extract
2 large eggs, separated, plus 2 additional egg whites
1 stick (8 tablespoons) salted butter, melted, plus softened butter, for serving
Warm syrup, for serving
Directions
Watch how to make this recipe.
Special equipment: Waffle iron
Preheat the waffle iron to the regular setting.
Sift together the flour, sugar, baking powder and salt in a bowl. In a separate bowl, whisk together the milk, vanilla and 2 egg yolks. Pour over the dry ingredients and very gently stir         until halfway combined. Pour in the melted butter and continue mixing very gently until combined.
In a separate bowl using a
whisk
(or a
mixer
), beat the 4 egg whites until stiff. Slowly fold them into the batter, stopping short of mixing them all the way         through.
Scoop the
batter
into your
waffle iron
in batches and cook according to its directions (lean toward the
waffles
being a little deep golden and
crisp
!). Serve immediately with softened butter and warm syrup.
Recipe courtesy of Ree Drummond`

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
		fmt.Println(i, score, line)
		scores[i] = score
	}

	fmt.Println(scores)

	minScore := 0
	for _, score := range scores {
		if score < minScore {
			minScore = score
		}
	}

	for i, score := range scores {
		scores[i] = score - minScore
	}
	fmt.Println(scores)

	gaussianScores := make([]int, len(scores))
	maxScore := 0
	maxScoreI := 0
	for i, _ := range scores {
		gaussianScore := 0
		if i > 0 {
			gaussianScore += scores[i-1] * 64
		}
		if i > 1 {
			gaussianScore += scores[i-2] * 32
		}
		gaussianScore += scores[i] * 100
		if i < len(scores)-1 {
			gaussianScore += scores[i+1] * 64
		}
		if i < len(scores)-2 {
			gaussianScore += scores[i+1] * 32
		}
		gaussianScores[i] = gaussianScore
		if gaussianScores[i] > maxScore {
			maxScore = gaussianScores[i]
			maxScoreI = i
		}
	}
	fmt.Println(gaussianScores)

	for testSigma := 1; testSigma < 10; testSigma++ {
		testGaussian := gaussianFunc(maxScoreI, testSigma, getMedian(gaussianScores), maxScore, gaussianScores)
		fmt.Println(residualFunc(testGaussian, gaussianScores))
	}

	//  / (sigma * sqrt(2 * pi)) * exp( -1 * (x - mu)^2 / (2 * sigma^2))

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
		residual += math.Pow(float64(one[i])-float64(two[i]), 2)
	}
	return residual
}

func gaussianFunc(mu2 int, sigma2 int, baseline2 int, max2 int, values []int) []int {
	sigma := float64(sigma2)
	mu := float64(mu2)
	baseline := float64(baseline2)
	max := float64(max2 - baseline2)
	x := make([]float64, len(values))
	for i, _ := range values {
		x[i] = float64(i)
	}
	fmt.Println(mu, sigma, x)
	gaussian := make([]int, len(x))
	for i, _ := range x {
		gaussian[i] = int(baseline + max/(sigma*math.Sqrt(2.0*3.1415926535))*math.Exp(-1*math.Pow(x[i]-mu, 2)/(2*math.Pow(sigma, 2))))
	}
	return gaussian
}
