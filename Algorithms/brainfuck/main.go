package main

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func BrainFuckInterpreter(code string) string {
	var result strings.Builder

	interpreter := make([]byte, 2048)
	pointer := 0
	openBracket := 0
	closeBracket := 0

	fmt.Println("interpreter:", string(interpreter), code, len(code))
	for i := 0; i < len(code); i++ {
		switch code[i] {
		case '>':
			fmt.Println(">>> :", pointer, i, code[i], string(code[i]), result.String())
			pointer++
		case '<':
			fmt.Println("<<< :", pointer, i, code[i], string(code[i]), result.String())
			pointer--
		case '+':
			fmt.Println("+++ :", pointer, i, code[i], string(code[i]), result.String())
			interpreter[pointer]++
		case '-':
			fmt.Println("--- :", pointer, i, code[i], string(code[i]), result.String())
			interpreter[pointer]--
		case '.':
			fmt.Println("... :", pointer, i, code[i], string(code[i]))
			result.WriteByte(interpreter[pointer])
		// case ',':
		// 	// TODO: implement
		case '[':
			fmt.Println("[ :", pointer, i, code[i], string(code[i]), interpreter[pointer], openBracket, result.String())
			if interpreter[pointer] == 0 {
				openBracket++
				i++
				for i < len(code) {
					if code[i] == ']' && openBracket == 1 {
						openBracket--
						break
					}
					if code[i] == '[' {
						openBracket++
					}
					if code[i] == ']' {
						openBracket--
					}
					i++
				}
			}
			if i < len(code) {
				fmt.Println("[ : after", pointer, i, code[i], string(code[i]), interpreter[pointer], openBracket)
			}
		case ']':
			fmt.Println("] :", pointer, i, code[i], string(code[i]), interpreter[pointer], closeBracket, result.String())
			if interpreter[pointer] != 0 {
				closeBracket++
				i--
				for i > 0 {
					if code[i] == '[' && closeBracket == 1 {
						closeBracket--
						break
					}
					if code[i] == ']' {
						closeBracket++
					}
					if code[i] == '[' {
						closeBracket--
					}
					i--
				}
			}
			if i < len(code) {
				fmt.Println("] : after", pointer, i, code[i], string(code[i]), interpreter[pointer], closeBracket)
			}
		default:
			// ignore
			// fmt.Println("Unknown instruction:", code[i], string(code[i]))
		}
	}
	return result.String()
}

func makeLoginTracker(maxAttempts int) func(success bool) bool {
	attempts := 0
	return func(success bool) bool {
		attempts++
		fmt.Println("Attempts:", attempts, maxAttempts)
		if success {
			attempts = 0
			return true
		}
		if attempts <= maxAttempts {
			return true
		}
		return false
	}
}

// Flatten will take an array of nested array and return
// all nested elements in an array. e.g. [[1,2,[3]],4] -> [1,2,3,4]
func Flatten(nested []any) []any {
	flattened := make([]any, 0)

	for _, i := range nested {
		switch i := i.(type) {
		case []any:
			flattenedSubArray := Flatten(i)
			flattened = append(flattened, flattenedSubArray...)
		case any:
			flattened = append(flattened, i)
		}
	}

	return flattened
}

func nestedSum(items []any) int {
	sum := 0

	for _, item := range items {
		switch item := item.(type) {
		case int:
			sum += item
		case []any:
			sum += nestedSum(item)
		}
	}

	return sum
}

type Report struct {
	EvenIndexSum int
	ScoreTotal   int
	RuneCount    int
}

func rangeReport(nums []int, scores map[string]int, word string) Report {

	var report Report
	for range word {
		report.RuneCount++
	}

	for _, v := range scores {
		report.ScoreTotal += v
	}

	for i, n := range nums {
		if i%2 == 0 {
			report.EvenIndexSum += n
		}
	}

	// var hh = "hfgdh"
	// fmt.Println(hh, hh[0:3])
	// // hh[:3] = 'a'
	// for i := range 7 {
	// 	fmt.Println(i, hh[i])
	// }

	return report
}

type Shape interface {
	Area() int
}

type Square struct {
	Side int
}

type Circle struct {
	Radius int
}

type Rectangle struct {
	Width  int
	Height int
}

func (r Rectangle) Area() int {
	return r.Width * r.Height
}

func (s Square) Area() int {
	return s.Side * s.Side
}

func (c Circle) Area() int {
	return (22 / 7) * c.Radius * c.Radius
}

func largestArea(shapes []Shape) int {
	maxArea := 0

	for _, shape := range shapes {
		area := shape.Area()
		if area > maxArea {
			maxArea = area
		}
	}

	return maxArea
}

type Priority int

const (
	Low = iota
	Medium
	High
	Urgent
)

type Task struct {
	Title    string
	Priority Priority
	Done     bool
}

func orderedTasks(tasks []Task) []string {
	// sort.Slice(tasks, func(i, j int) bool {
	// 	if tasks[i].Priority == tasks[j].Priority {
	// 		return tasks[i].Title < tasks[j].Title
	// 	}
	// 	return tasks[i].Priority > tasks[j].Priority
	// })

	// TODO: check out efficiency
	var unfinisedTasks []string
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	for _, task := range tasks {
		if !task.Done {
			unfinisedTasks = append(unfinisedTasks, task.Title)
		}
	}
	return unfinisedTasks
}

func moveZeros(nums []int) []int {
	// var lastNoneZeroIndex = -1

	// for i, num := range nums {
	// 	if num != 0 {
	// 		lastNoneZeroIndex++
	// 		if lastNoneZeroIndex != i {
	// 			nums[lastNoneZeroIndex] = num
	// 			nums[i] = 0
	// 		}
	// 	}
	// }

	// return nums

	var movedSice []int
	var numOfZeros int

	for _, num := range nums {
		// fmt.Println("test:", numOfZeros)
		// numOfZeros++
		if num != 0 {
			movedSice = append(movedSice, num)
		} else {
			numOfZeros++
		}
	}

	return movedSice

	// for i := 0; i < len(nums); i++ {
	// 	if nums[i] != 0 {
	// 		nums[lastNoneZeroIndex] = nums[i]
	// 		lastNoneZeroIndex++
	// 	}
	// }

	// for i := lastNoneZeroIndex; i < len(nums); i++ {
	// 	nums[i] = 0
	// }
	// return nums
}

func isValidBrackets(text string) bool {
	// dd := text[:3]
	// fmt.Println(dd)
	seenBrackets := ""
	// var ss strings.Builder
	bracketsMap := map[rune]byte{
		'(': ')',
		')': '(',
		'[': ']',
		']': '[',
		'{': '}',
		'}': '{',
	}

	for _, brac := range text {
		// var oppBrac byte = ' '
		// switch brac {
		// case '(':
		// 	oppBrac = ')'
		// case ')':
		// 	oppBrac = '('
		// case '[':
		// 	oppBrac = ']'
		// case ']':
		// 	oppBrac = '['
		// case '{':
		// 	oppBrac = '}'
		// case '}':
		// 	oppBrac = '{'
		// }

		if oppBrac, ok := bracketsMap[brac]; ok {
			if len(seenBrackets) > 0 && seenBrackets[len(seenBrackets)-1] == oppBrac {
				seenBrackets = seenBrackets[:len(seenBrackets)-1]
				// ss.W
			} else {
				seenBrackets += string(brac)
			}
		}

	}

	fmt.Println("remaining brackets:", text, seenBrackets)
	return len(seenBrackets) == 0
}

func isAnagram(w1 string, w2 string) bool {
	if len(w1) != len(w2) {
		return false
	}

	var w1mapfreq = make(map[rune]int)

	for _, ch := range w1 {
		w1mapfreq[ch]++
	}

	for _, ch := range w2 {
		if _, ok := w1mapfreq[ch]; ok {
			w1mapfreq[ch]--
		}
	}

	for _, v := range w1mapfreq {
		if v != 0 {
			return false
		}
	}

	return true
}

func groupAnagrams(words []string) [][]string {
	// var groups [][]string
	groups := [][]string{}

	for _, word := range words {
		var anagramSeen bool
		for i, group := range groups {
			if len(group) > 0 && isAnagram(word, group[0]) {
				group = append(group, word)
				groups[i] = group
				anagramSeen = true
				break
			}
		}

		if anagramSeen {
			continue
		} else {
			groups = append(groups, []string{word})
		}
	}

	return groups
}

func wrapText(words []string, maxWidth int) []string {
	// var wrapped []string
	wrapped := []string{}

	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word) <= maxWidth {
			currentLine += word
			if len(currentLine) < maxWidth {
				currentLine += " "
			}
		} else {
			remainingWidth := maxWidth - len(currentLine)
			for range remainingWidth {
				currentLine += " "
			}
			wrapped = append(wrapped, currentLine)
			currentLine = word + " "
		}
	}

	if len(currentLine) > 0 {
		remainingWidth := maxWidth - len(currentLine)
		for range remainingWidth {
			currentLine += " "
		}
		wrapped = append(wrapped, currentLine)
	}

	if wrapped == nil {
		return []string{}
	}
	return wrapped
}

func rotateGrid(grid [][]rune) [][]rune {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return [][]rune{}
	}
	// rowLen, colLen := len(grid), len(grid[0])
	rotated := make([][]rune, len(grid[0]))
	rotatedStr := make([][]string, len(grid[0]))

	for i := len(grid) - 1; i >= 0; i-- {
		for j, c := range grid[i] {
			rotated[j] = append(rotated[j], c)
			rotatedStr[j] = append(rotatedStr[j], string(c))
		}
	}

	fmt.Println(rotatedStr)

	return rotated
}

func runLengthEncode(text string) string {
	var encoded strings.Builder

	var current byte
	var count int
	for i := 0; i < len(text); i++ {
		if i == 0 {
			current = text[i]
			count++
			continue
		}

		if i > 0 {
			if text[i] == text[i-1] {
				count++
			} else {
				encoded.WriteByte(current)
				countStr := strconv.Itoa(count)
				encoded.WriteString(countStr)
				current = text[i]
				count = 1
			}
		}
	}

	if count > 0 {
		encoded.WriteByte(current)
		countStr := strconv.Itoa(count)
		encoded.WriteString(countStr)
	}

	return encoded.String()
}

func countIslands(grid [][]rune) int {
	if len(grid) == 0 {
		return 0
	}

	islandCount := 0
	seenLands := make(map[string]bool)

	var searchLandCells func(row, col int)

	searchLandCells = func(row, col int) {
		if row < 0 || row >= len(grid) || col < 0 || col >= len(grid[row]) {
			return
		}

		cellLoc := strconv.Itoa(row) + "," + strconv.Itoa(col)

		// if seenLands[cellLoc] || grid[row][col] != '#' {
		// 	return
		// }

		if grid[row][col] == '#' && !seenLands[cellLoc] {
			seenLands[cellLoc] = true
			searchLandCells(row+1, col)
			searchLandCells(row, col+1)
			searchLandCells(row-1, col)
			searchLandCells(row, col-1)
		}

	}

	for row := range grid {
		for col := range grid[row] {
			cell := grid[row][col]
			if cell == '#' {
				cellLoc := strconv.Itoa(row) + "," + strconv.Itoa(col)
				if !seenLands[cellLoc] {
					islandCount++
					searchLandCells(row, col)
				}
			}
		}
	}

	return islandCount
}

func countIslands1(grid [][]rune) int {
	if len(grid) == 0 {
		return 0
	}

	rows := len(grid)
	cols := len(grid[0])

	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	var dfs func(int, int)

	dfs = func(r, c int) {
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return
		}

		if visited[r][c] || grid[r][c] != '#' {
			return
		}

		visited[r][c] = true

		dfs(r+1, c)
		dfs(r-1, c)
		dfs(r, c+1)
		dfs(r, c-1)
	}

	count := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == '#' && !visited[r][c] {
				count++
				dfs(r, c)
			}
		}
	}

	return count
}

func main() {
	if len(os.Args) != 2 {
		// fmt.Println("Usage: go run main.go <brainfuck_code>")
		fmt.Println()
		// return
	}
	// codeText := os.Args[1]
	// output := BrainFuckInterpreter(codeText)
	// fmt.Print(output)

	tracker := makeLoginTracker(2)
	test := []string{"abc", "", "abcd", "test", "test1"}
	for _, t := range test {
		fmt.Println(t, " : ", tracker(t == "test"))
	}

	fmt.Println(nestedSum([]any{10, []any{-2, []any{5}}, 1}))    // 14
	fmt.Println(nestedSum([]any{1, 2, 3, 4, 5}))                 // 15
	fmt.Println(nestedSum([]any{[]any{}, []any{1, []any{}}, 2})) // 3
	fmt.Println(nestedSum([]any{}))                              // 0

	fmt.Println(rangeReport([]int{2, 3, 4, 5}, map[string]int{"a": 10, "b": 20}, "go"))

	fmt.Println(largestArea([]Shape{Rectangle{Width: 10, Height: 2}, Square{Side: 5}, Circle{Radius: 2}}))

	fmt.Println(orderedTasks([]Task{{Title: "Task 1", Priority: Medium, Done: false}, {Title: "Task 2", Priority: High, Done: false}, {Title: "Task 3", Priority: Low, Done: false}, {Title: "Task 4", Priority: Low, Done: false}, {Title: "Task 5", Priority: High, Done: false}}))

	fmt.Println(moveZeros([]int{}))
	fmt.Println(moveZeros([]int{0}))
	fmt.Println(moveZeros([]int{0, 1, 0, 3, 12}))
	fmt.Println(moveZeros([]int{4, 0, -2, 0, 5}))
	fmt.Println(moveZeros([]int{1, 2, 0, 3}))

	ss := []string{"hg", "kj", "gg"}

	fmt.Println("testing slices delete before:", ss)

	sd := slices.Replace(ss, 1, 2)

	fmt.Println("testing slices delete:", ss, sd)

	fmt.Println(isValidBrackets("()[]{}"))
	fmt.Println(isValidBrackets(""))
	fmt.Println(isValidBrackets("([)]"))
	fmt.Println(isValidBrackets("{[]}"))
	fmt.Println(isValidBrackets("("))
	fmt.Println(isValidBrackets("]"))

	fmt.Println("is anagram test", isAnagram("ate", "eat"))

	// group := []string{}

	// fmt.Println("Initial group:", group)

	// // // group[len(group)+1] = "hello"
	// // group1 := append(*&group, "hello")
	// // group = group1

	// fmt.Println("After group:", group)
	fmt.Println(groupAnagrams([]string{"eat", "tea", "tan", "ate", "nat", "bat"}))
	cc := groupAnagrams([]string{})
	fmt.Println(cc)
	tt := groupAnagrams([]string{""})
	fmt.Println(tt, len(tt), len(tt[0]), tt[0][0], len(tt[0][0]))
	fmt.Printf("%#v\n%#v\n", tt, cc)

	fmt.Println(wrapText([]string{"This", "is", "an", "example", "of", "text"}, 10))
	fmt.Println(wrapText([]string{"a", "b", "c", "d"}, 3))
	fmt.Println(wrapText([]string{"a"}, 1))
	fmt.Println(wrapText([]string{"ab", "cd", "ef"}, 5))
	fmt.Println(wrapText([]string{""}, 5))
	fmt.Println(wrapText([]string{}, 5))

	fmt.Println(rotateGrid([][]rune{{'a', 'b', 'c'}, {'d', 'e', 'f'}}))
	fmt.Println(rotateGrid([][]rune{}))
	fmt.Println(rotateGrid([][]rune{{'x'}}))
	fmt.Println(rotateGrid([][]rune{{'a'}, {'b'}, {'c'}}))

	fmt.Println(runLengthEncode(""))
	fmt.Println(runLengthEncode("aaabbc"))
	fmt.Println(runLengthEncode("abcd"))
	fmt.Println(runLengthEncode("aaaaaaaaaaaa"))

	// fmt.Println(countIslands([][]rune{{'#', '#', '.', '.'}, {'.', '#', '#', '.'}, {'.', '.', '#', '#'}, {'#', '.', '.', '#'}}))
	// fmt.Println(countIslands([][]rune{
	// 	{'#', '.', '#'},
	// 	{'#', '.', '#'},
	// 	{'#', '#', '#'},
	// }))
	// fmt.Println(countIslands([][]rune{{'.', '.', '.'}, {'.', '.', '.'}, {'.', '.', '.'}}))
	fmt.Println(countIslands([][]rune{
		// {},
		{'.', '.', '.', '.', '.', '.'},
		{'#', '#', '.', '.', '.', '.', '#'},
		{'#', '.', '.', '.', '.', '.', '#'},
		{'.', '.', '.', '.', '.', '.', '.'},
		{'.', '.', '.', '.', '.', '#'},
		{'#', '#', '.', '.', '.', '#', '.'},
		// {'#', '.', '.', '.', '.'},
	}))
	fmt.Println(countIslands([][]rune{}))

	// var artist struct {
	// 	ID int `json:"id"`
	// }
	// json.Unmarshal([]byte(`{"id": 1}`), &artist)

	// fmt.Println(artist)
}
