package main

import (
	"fmt"
)

var commonFoods []string
var goodContextWords []string
var badContextWords []string

func init() {
	loadData()
}

func main() {

	fmt.Println(commonFoods)
	fmt.Println(badContextWords)
	fmt.Println(goodContextWords)
}
