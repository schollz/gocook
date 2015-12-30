package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

var commonFoods []string
var goodContextWords []string
var badContextWords []string
var pairing map[[2]string]int
var uniqueFoods []string

// func gaussianFilter(data []int) (filtered []int) {

func init() {
	loadData()

	pairing = make(map[[2]string]int, 221777)

	f, err := ioutil.ReadFile("resources/pairing.csv")
	if err != nil {
		panic(err)
	}
	for _, l := range strings.Split(string(f), "\n") {
		arr := strings.Split(string(l), ",")
		if len(arr) == 3 {
			k := arr[:2]
			sort.Sort(sort.StringSlice(k))
			key := [2]string{k[0], k[1]}

			val, err := strconv.Atoi(arr[2])
			if err != nil {
				panic(err)
			}
			pairing[key] = val

		}
	}
	fmt.Println(pairing[[2]string{"mint oil", "pinto bean"}])
	fmt.Println(len(pairing), len(uniqueFoods))

}

func main() {

	// fmt.Println(commonFoods)
	// fmt.Println(badContextWords)
	// fmt.Println(goodContextWords)

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
	testContent = parseURL("http://www.foodnetwork.com/recipes/food-network-kitchens/slow-cooker-pork-tacos-recipe.html")

	text := getIndredientText(testContent)

	text = strings.Replace(text, ",", "", -1)
	text = strings.Replace(text, ".", "", -1)
	text = strings.Replace(text, "!", "", -1)
	text = strings.Replace(text, "?", "", -1)
	text = strings.Replace(text, "olive oil", "olivoil", -1)
	text = strings.Replace(text, "cider vinegar", "vinegar", -1)
	fmt.Println(text)

	fmt.Println("\nFrom food pairing:")
	var ingredients []string
	for _, food := range commonFoods {
		if strings.Contains(text, " "+food) && len(food) > 1 {
			ingredients = append(ingredients, food)
		}
	}

	score := float64(0)
	foods := float64(0)
	for i, food1 := range ingredients {
		for j, food2 := range ingredients {
			if j > i {
				score += float64(pairing[[2]string{food1, food2}])
				score += float64(pairing[[2]string{food2, food1}])
				if pairing[[2]string{food1, food2}] > 0 {
					fmt.Println(food1, food2, pairing[[2]string{food1, food2}])
				}
				if pairing[[2]string{food2, food1}] > 0 {
					fmt.Println(food2, food1, pairing[[2]string{food2, food1}])
				}
				if score > 0 {
					foods += float64(1)
				}
			}
		}
	}
	fmt.Println(score, foods, score/foods)

}
