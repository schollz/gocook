package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var commonFoods []string
var goodContextWords []string
var badContextWords []string
var pairing map[[2]string]int
var uniqueFoods []string
var cdfTable map[int]int

// func gaussianFilter(data []int) (filtered []int) {

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func loadCdfData() {
	cdfTable = make(map[int]int)
	f, err := ioutil.ReadFile("resources/cdf.tab")
	if err != nil {
		panic(err)
	}
	for _, l := range strings.Split(string(f), "\n") {
		arr := strings.Split(string(l), " ")
		if len(arr) == 2 {
			val1, err := strconv.Atoi(arr[0])
			if err != nil {
				panic(err)
			}
			val2, err := strconv.Atoi(arr[1])
			if err != nil {
				panic(err)
			}
			cdfTable[val1] = val2
		}
	}

	fmt.Println(cdfTable)

}

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
	fmt.Println(len(pairing), len(commonFoods))

	loadCdfData()
	getCdf(49)
}

type FoodJson struct {
	Score int
	Good  []string
	Bad   []string
}

func main() {
	fmt.Println(getScore("http://www.foodnetwork.com/recipes/alton-brown/eggs-benedict-recipe.html"))
	runtime.GOMAXPROCS(runtime.NumCPU() - 1) // one core for wrk
	http.HandleFunc("/score", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			url := r.URL.Query()["url"]
			fmt.Println(url)
			if len(url) > 0 {
				defer timeTrack(time.Now(), "/url="+url[0])
				b, err := json.Marshal(getScore(url[0]))
				if err != nil {
					panic(err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(b))
			}
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	fmt.Println("Running on Port 4000")
	http.ListenAndServe(":4000", nil)
}

func getScore(url string) FoodJson {
	testContent := parseURL(url)
	text := getIndredientText(testContent)

	text = strings.ToLower(text)
	text = strings.Replace(text, ",", "", -1)
	text = strings.Replace(text, ".", "", -1)
	text = strings.Replace(text, "!", "", -1)
	text = strings.Replace(text, "?", "", -1)
	text = strings.Replace(text, "olive oil", "olivoil", -1)
	text = strings.Replace(text, "cider vinegar", "vinegar", -1)
	text = strings.Replace(text, "nutrition", "", -1)
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
	var badCombos []string
	var goodCombos []string
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
					if score > 70 {
						goodCombos = append(goodCombos, food1+","+food2)
					}
					if score < 30 {
						badCombos = append(badCombos, food1+","+food2)
					}
				}
			}
		}
	}
	newScore := -1
	if score > 0 {
		newScore = getCdf(int(score / foods))
	}
	return FoodJson{
		Score: newScore,
		Bad:   badCombos,
		Good:  goodCombos,
	}

}

func getCdf(val int) int {
	bestDiff := 10000
	bestVal := 0
	for k, _ := range cdfTable {
		diff := k - val
		if diff < 0 {
			diff = diff * -1
		}
		if diff < bestDiff {
			bestDiff = diff
			bestVal = k
		}
	}
	return cdfTable[bestVal]
}
