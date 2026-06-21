package main

import (
	"fmt"
	"slices"
)

type SS interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~string
}

// type ss interface{ ss }

func minRooms[T ~[][2]E, E SS](intervals T) int {
	steps := 0
	var occupiedRooms [][2]E

	// slices.SortFunc(intervals, func(a, b [2]E) int {
	// 	if a[0] < b[0] {
	// 		return -1
	// 	}
	// 	if a[0] > b[0] {
	// 		return 1
	// 	}
	// 	return 0
	// })

	for j, interval := range intervals {
		if len(occupiedRooms) == 0 {
			steps++
			fmt.Println("step", steps, j)
			occupiedRooms = append(occupiedRooms, interval)
			continue
		}
		startInterval := interval[0]
		// check through occupied rooms for available space
		for i, room := range occupiedRooms {
			steps++
			fmt.Println("step", steps, j, i)
			fmt.Println("checking interval:", interval, occupiedRooms, i, room)
			end := room[1]
			// if available, replace occupied room with current interval
			if end <= startInterval {
				remaininRooms := occupiedRooms[i+1:]
				occupiedRooms = append(occupiedRooms[:i], interval)
				occupiedRooms = append(occupiedRooms, remaininRooms...)
				fmt.Println("replaced an occupied room:", i, interval, occupiedRooms)
				break
			}
			// if no available room found at the end, add a new room to occupied rooms.
			if i == len(occupiedRooms)-1 {
				occupiedRooms = append(occupiedRooms, interval)
				fmt.Println("added new room:", i, interval, occupiedRooms)
			}
		}
	}
	// fmt.Println("occupied rooms:", occupiedRooms, occupiedRooms[1+1:], occupiedRooms[:1])
	fmt.Println("steps:", steps)

	return len(occupiedRooms)
}

func minRooms1(intervals [][2]int) int {
	steps := 0
	if len(intervals) == 0 {
		return 0
	}

	startTimes := make([]int, len(intervals))
	endTimes := make([]int, len(intervals))

	for i, interval := range intervals {
		startTimes[i] = interval[0]
		endTimes[i] = interval[1]
	}

	slices.Sort(startTimes)
	slices.Sort(endTimes)

	rooms := 0
	maxRooms := 0
	i, j := 0, 0

	// endPointer := 0
	// for _, startTime := range startTimes {
	// 	steps++
	// 	if startTime < endTimes[endPointer] {
	// 		rooms++
	// 	} else {
	// 		endPointer++
	// 	}
	// }

	for i < len(startTimes) {
		steps++
		fmt.Println("step", steps, i)
		if startTimes[i] < endTimes[j] {
			rooms++
			maxRooms = max(maxRooms, rooms)
			i++
		} else {
			rooms--
			j++
		}
	}
	fmt.Println("steps:", steps)

	return maxRooms
}

func variadicfunc(removeLowest bool, args ...int) int {
	sum := 0
	if len(args) == 0 {
		return 0
	}

	if len(args) == 1 {
		return args[0]
	}

	slices.Sort(args)

	fmt.Println("sorted args:", args, "len of args", len(args))

	if removeLowest {
		args = args[1:]
	}

	for _, arg := range args {
		sum += arg
	}

	return sum / len(args)
}

type Account struct {
	Balance int
	Active  bool
}

func applyTransaction(acct *Account, amount int) bool {

	fmt.Println("account arg:", acct)
	// fmt.Println("account arg pp:", *acct)

	if acct != nil && acct.Active && (acct.Balance+amount) >= 0 {
		acct.Balance += amount
		return true
	}

	return false
}

type Calculator interface {
	Add(a, b int) int
}

type calculator struct {
	result int
}

func NewCalculator() *calculator {
	return &calculator{}
}

func (c *calculator) Set(num int) {
	fmt.Println("setting result:", c, "address:", &c, "result:", c.result, "to", num)
	c.result = num
}

func (c *calculator) Add(nums ...int) int {
	for _, num := range nums {
		c.result += num
	}
	return c.result
}

func (c *calculator) Subtract(nums ...int) int {
	for _, num := range nums {
		c.result -= num
	}
	return c.result
}

func (c *calculator) Multiply(nums ...int) int {
	for _, num := range nums {
		c.result *= num
	}
	return c.result
}

func (c *calculator) Divide(nums ...int) int {
	for _, num := range nums {
		if num == 0 {
			continue
		}
		c.result /= num
	}
	return c.result
}

func (c *calculator) Current() int {
	return c.result
}

func (c *calculator) Reset() int {
	c.result = 0
	return c.result
}

type SAccount struct {
	Name   string
	Active bool
}

type Seller struct {
	SAccount
	Sales  int
	Rating int
}

func topSeller(sellers []Seller) string {
	highestActiveScore := 0
	var topSeller string

	for _, seller := range sellers {
		if seller.Active {
			sellerScore := seller.Sales * seller.Rating
			if sellerScore > highestActiveScore {
				highestActiveScore = sellerScore
				topSeller = seller.Name
			}
		}
	}

	return topSeller
}

func firstDuplicate[T comparable](items []T) (T, bool) {
	var seenItems []T
	for _, item := range items {
		if slices.Contains(seenItems, item) {
			return item, true
		}
		seenItems = append(seenItems, item)
	}

	var zero T
	return zero, false
}

type IntIterator func(yield func(int) bool)

func (f IntIterator) Map(g func(int) int) IntIterator {
	return func(yield func(int) bool) {
		f(func(i int) bool {
			return yield(g(i))
		})
	}
}

// func (f IntIterator) Map(g func(int) int) IntIterator {
// 	return func(yield func(int) bool) bool {
// 		return f(func(i int) bool {
// 			return yield(g(i))
// 		})
// 	}
// }

// func numbers1() IntIterator {
// 	return func(yield func(int) bool) bool {
// 		for i := 0; i < 10; i++ {
// 			if !yield(i) {
// 				break
// 			}
// 		}
// 		return true
// 	}
// }

var numbers = func(kk func(int) bool) {
	kk(10)
	kk(20)
	kk(30)

	// return true
}

func sumIterators(iter IntIterator) int {
	sum := 0

	fmt.Println("iterators:", iter)

	iter(func(i int) bool {
		sum += i
		return true
	})

	iter = iter.Map(func(i int) int {
		return i * 2
	})

	sum2 := 0

	// iter(func(i int) bool {
	// 	sum2 += i
	// 	return true
	// })

	for i := range iter {
		sum2 += i
	}

	fmt.Println("iterators after map:", sum2)

	return sum
}

func main() {
	// m := minRooms1([][2]int{{30, 75}, {0, 50}, {60, 150}})
	// m := minRooms1([][2]int{{5, 10}, {0, 30}, {15, 20}})
	// m := minRooms1([][2]int{{1, 5}, {2, 4}, {4, 6}})
	// m := minRooms([][2]int{{1, 2}, {2, 3}, {3, 4}})
	// fmt.Println("minimum rooms:", m)

	// fmt.Println("average:", variadicfunc(true, 80, 90, 100, 70, 120))

	// var acct Account
	// acct = Account{Balance: 500, Active: true}
	// fmt.Println("account:", applyTransaction(&acct, -600))
	// fmt.Println("account:", acct.Balance)

	calc := NewCalculator()
	fmt.Println("pointer calculator:", calc, calc.Add(1, 2))
	fmt.Println("pointer calculator:", calc.result)

	calc2 := calculator{}
	// fmt.Println("calculator:", calc2)
	// fmt.Println("calculator:", calc2.Add(1, 2))
	// fmt.Println("calculator:", calc2.result)
	// calc2.result = 10
	// fmt.Println("calculator set to 10:", calc2)
	// fmt.Println("calculator set to 10:", calc2.Add(1, 2))
	// fmt.Println("calculator set to 10:", calc2.result)
	// calc2.Set(30)
	// fmt.Println("calculator:", calc2)
	// fmt.Println("calculator:", calc2.Add(1, 2))
	// fmt.Println("calculator:", calc2.result)

	fmt.Println("calc current:", calc2.Current())
	fmt.Println("calc add:", calc2.Add(1, 2))
	fmt.Println("calc add:", calc2.Add(3))
	fmt.Println("calc multiply:", calc2.Multiply(2))
	fmt.Println("calc subtract:", calc2.Subtract(5))
	fmt.Println("calc divide:", calc2.Divide(2))
	fmt.Println("calc current:", calc2.Current())

	calc3 := calculator{}
	fmt.Println("calc3 sub:", calc3.Subtract(10, 3, 2))
	fmt.Println(calc3.Divide(0))
	fmt.Println(calc3.Reset())
	fmt.Println(calc3.Current())

	sellers := []Seller{
		{
			SAccount: SAccount{
				Name:   "Ada",
				Active: true,
			},
			Sales:  100,
			Rating: 5,
		},
		{
			SAccount: SAccount{
				Name:   "Ben",
				Active: true,
			},
			Sales:  200,
			Rating: 3,
		},
		{
			SAccount: SAccount{
				Name:   "Tola",
				Active: false,
			},
			Sales:  200,
			Rating: 3,
		},
	}
	fmt.Println("top seller:", topSeller(sellers))
	item, dup := firstDuplicate([]string{"a", "b", "c", "d", "a"})
	fmt.Println("first duplicate:", item, dup)
	item2, dup2 := firstDuplicate([]string{"a", "b", "c", "d", "e"})
	fmt.Println("first duplicate:", item2, dup2)

	item3, dup3 := firstDuplicate([]int{1, 2, 3, 4, 5, 6, 2})
	fmt.Println("first duplicate:", item3, dup3)
	item4, dup4 := firstDuplicate([]int{1, 2, 3, 4, 5, 6, 7})
	fmt.Println("first duplicate:", item4, dup4)

	fmt.Println(firstDuplicate([]float64{1.1, 2.2, 3.3, 2.2}))

	fmt.Println("sum iterators:", sumIterators(numbers))

	EightQueens()
}
