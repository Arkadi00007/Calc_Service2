package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func MyLimit(num []int, limit int, n int) ([]int, error) {
	if num == nil {
		return []int{}, errors.New("false")
	}
	if n < 0 {
		return nil, errors.New("[]")
	}
	if len(num) == 0 {
		return num, nil
	}
	b := make([]int, 0, n)
	if n == 0 {
		return b, nil
	}
	for _, i := range num {
		if i < limit && len(b) < n {
			b = append(b, i)
		}
		if len(b) == n {
			break
		}
	}
	return b, nil
}

//	func Clean(nums []int, x int) []int {
//		var a int = 1
//		for i, n := range nums {
//			if n == x {
//				nums[i] = nums[len(nums)-a]
//				a++
//			}
//		}
//		return nums[a-1:]
//	}

func IsLatin(input string) bool {
	a := strings.ToLower(input)
	b := true
	for _, i := range a {
		if i < 97 || i > 122 {
			b = false
		}

	}
	return b
}

func containsDuplicate(nums []int) bool {
	set := make(map[int]bool)
	for _, num := range nums {
		if _, ok := set[num]; ok { // если такой ключ существует, переходим к return
			return true
		}
		set[num] = true
	}
	return false
}

func CountingSort(contacts []string) map[string]int {
	mapp := make(map[string]int)
	for _, i := range contacts {
		if _, ok := mapp[i]; ok {
			mapp[i]++
		} else {
			mapp[i] = 1
		}
	}
	return mapp
}

func DeleteLongKeys(m map[string]int) map[string]int {
	for k := range m {
		if len(k) < 6 {
			delete(m, k)
		}
	}
	return m
}

type Task struct {
	summary     string
	description string
	deadline    time.Time
	priority    int
}

type Tasks interface {
	IsOverdue() bool
}

func (s Task) IsOverdue() bool {
	return time.Now().After(s.deadline)
}

func (s Task) IsTopPriority() bool {
	if s.priority == 4 || s.priority == 5 {
		return true
	}
	return false
}

type Note struct {
	title string
	text  string
}

type ToDoList struct {
	name  string
	tasks []Task
	notes []Note
}

func (s ToDoList) TasksCount() int {
	return len(s.tasks)
}

func (s ToDoList) NotesCount() int {
	return len(s.notes)
}

func (s ToDoList) CountTopPrioritiesTasks() int {
	count := 0
	for _, task := range s.tasks {
		if task.IsTopPriority() {
			count++
		}
	}
	return count
}

func (s ToDoList) CountOverdueTasks() int {
	count := 0
	for _, task := range s.tasks {
		if task.IsOverdue() {
			count++
		}
	}
	return count
}

func PrintType(x interface{}) {
	switch x.(type) {
	case string:
		fmt.Println("string")
	case int:
		fmt.Println("int")
	case bool:
		fmt.Println("bool")
	default:
		fmt.Println("unknown")
	}
}

func GetCharacterAtPosition(str string, position int) (rune, error) {
	if position >= len(str) {
		return 0, errors.New("position out of range")
	}
	return []rune(str)[position], nil
}

func Factorial(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("factorial is not defined for negative numbers")
	}
	var a int
	for i := 1; i <= n; i++ {
		a *= i
	}
	return a, nil
}

func SumTwoIntegers(a, b string) (int, error) {
	num, err := strconv.Atoi(a)
	num1, err1 := strconv.Atoi(b)
	if err != nil && err1 != nil {
		return 0, fmt.Errorf("invalid input, please provide two integers")
	}
	return num + num1, nil
}

func AreAnagrams(str1, str2 string) bool {
	dict := make(map[string]int)
	dict2 := make(map[string]int)
	var set1 []string
	var set2 []string
	sort.Strings(set1)
	sort.Strings(set2)

	for _, i := range str1 {
		set1 = append(set1, string(i))
	}
	for _, i := range str2 {
		set2 = append(set2, string(i))
	}
	for _, i := range set1 {
		dict[i]++
	}
	for _, i := range set2 {
		dict2[i]++
	}
	b := true
	for k, _ := range dict {
		if dict[k] != dict2[k] {
			b = false
		} else {
			b = true
		}
	}
	return b
}

func PrintHello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func Multiply(a, b int) int {
	return a * b
}

//func TestPrintHello(t *testing.T) {
//	got := PrintHello("Igor")
//	expected := "Hello, Igor!"
//
//	if got != expected {
//		t.Fatalf(`PrintHello("Igor") = %q, want %q`, got, expected)
//	}
//}

func main() {
	a := []struct {
		name string
		age  int
	}{
		{"arkadi", 7},
		{"artyom чёрт", 11},
		{"tigran", 111},
	}
	type ara struct {
		name string
		age  int
	}

	arr := ara{"arkadi", 111}

	fmt.Println(arr.name)
	for _, j := range a {
		fmt.Printf("%s - %d\n", j.name, j.age)

	}

}

//_, err := fmt.Scanln(&a)
//if err != nil {
//	fmt.Println("Некорректный ввод")
//	return
//}

//if n <= 0 {
//fmt.Println("NO")
//}
//if (n & (n - 1)) == 0 {
//fmt.Println("YES")
//} else {
//fmt.Println("NO")
//}
//
//}
