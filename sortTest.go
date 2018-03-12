package main

import (
	"fmt"
	"sort"
	"strings"
)
func main() {
	fmt.Println("\n######## SearchInts not works in descending order  ######## ")
	intSlice := []int{55, 54, 53, 52, 51, 50, 48, 36, 15, 5}    // sorted slice in descending
	x := 36
	pos := sort.SearchInts(intSlice,x)
	fmt.Printf("Found %d at index %d in %v\n", x, pos, intSlice)
	fmt.Println("\n######## Search works in descending order  ########")
	i := sort.Search(len(intSlice), func(i int) bool { return intSlice[i] <= x })
	fmt.Printf("Found %d at index %d in %v\n", x, i, intSlice)
	fmt.Println("\n\n######## SearchStrings not works in descending order  ######## ")
	// sorted slice in descending
	strSlice := []string{"Washington","Texas","Ohio","Nevada","Montana","Indiana","Alaska"}
	y := "Montana"
	posstr := sort.SearchStrings(strSlice,y)
	fmt.Printf("Found %s at index %d in %v\n", y, posstr, strSlice)
	fmt.Println("\n######## Search works in descending order  ########")
	j := sort.Search(len(strSlice), func(j int) bool {return strSlice[j] <= y})
	fmt.Printf("Found %s at index %d in %v\n", y, j, strSlice)
	fmt.Println("\n######## Search works in ascending order  ########")
	fltSlice := []float64{10.10, 20.10, 30.15, 40.15, 58.95} // string slice in float64
	z := 40.15
	k := sort.Search(len(fltSlice), func(k int) bool {return fltSlice[k] >= z})
	fmt.Printf("Found %f at index %d in %v\n", z, k, fltSlice)


var n = 5
var m = 10
 a  := make([][]int, 90)
fmt.Println(a)
OuterLoop:
	for i = 0; i < n; i++ {
		fmt.Println("I ES :")
		fmt.Println(i)
		for j = 0; j < m; j++ {
			if j == 6{
				fmt.Println("J es: ")
				fmt.Println(j)
				break OuterLoop
			}

		}
		fmt.Println("SALIO DEL CICLO INTERNO")
	}
	fmt.Println("SALIO DEL CICLO EXTERNO")
	w1 := "bcad"
	w2 := SortString(w1)

	fmt.Println(w1)
	fmt.Println(w2)
	var abc = []string{"F1","F10","F11","F12","F13","F14","F15","F16","F17","F18","F19","F2","F20","F21","F22","F23","F24","F25","F26","F27",
		"F28","F29","F3","F30","F31","F32","F33","F34","F35","F36","F37","F38","F39","F4","F40","F41"}



	fmt.Println("10 antes de los dos puntos")
	fmt.Println(abc[10:])
	fmt.Println("10 despues de los dos puntos")

	fmt.Println(abc[:10])


	var uno = 2
	var dos = 5
	hits := make(map[int]map[int]string)
	fmt.Println("vacio",len(hits))
	hits[0] =make(map [int]string)
	hits[1] =make(map [int]string)
	hits[2] =make(map [int]string)
	fmt.Println("lllll")
	hits[uno][dos] = "GRACE"
	fmt.Println("Con un elemento",len(hits))
	fmt.Println(hits[uno][dos])



	x1 := "chars@arefun"

	ii := strings.Index(x1, "@")
	fmt.Println("Index: ", ii)
	if i > -1 {
		chars := x1[:ii]
		arefun := x1[ii+1:]
		fmt.Println(chars)
		fmt.Println(arefun)
	} else {
		fmt.Println("Index not found")
		fmt.Println(x)
	}
	str := "a space-separated string"
	str = strings.Replace(str, "t", ",", 1)
	fmt.Println(str)
}


type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func SortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}
