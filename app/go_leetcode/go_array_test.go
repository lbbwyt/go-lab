package go_leetcode

import (
	"fmt"
	"testing"
)

func TestSearchRange(t *testing.T) {
	res := SearchRange([]int{1, 5, 5, 6}, 5)
	fmt.Println(fmt.Sprintf("%v", res))
}

func TestCombinationSum(t *testing.T) {
	nums := []int{2, 3, 5}
	nums = nums[:len(nums)]
	fmt.Println(fmt.Sprintf("%v", nums))
	target := 8
	res := CombinationSum(nums, target)
	fmt.Println(fmt.Sprintf("%v", res))
}
