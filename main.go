package main

import (
	"fmt"
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
	testContent = parseURL("http://allrecipes.com/recipe/31848/jambalaya/")

	text := getIndredientText(testContent)
	fmt.Println(text)

	text = strings.Replace(text, ",", "", -1)
	text = strings.Replace(text, ".", "", -1)
	text = strings.Replace(text, "!", "", -1)
	text = strings.Replace(text, "?", "", -1)

	for _, food := range commonFoods {
		if strings.Contains(text, " "+food+" ") {
			fmt.Println(food)
		}
	}

}
