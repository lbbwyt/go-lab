package go_leetcode

import "sort"

func SearchRange(nums []int, target int) []int {
	leftmost := sort.SearchInts(nums, target)
	if leftmost == len(nums) || nums[leftmost] != target {
		return []int{-1, -1}
	}
	rightmost := sort.SearchInts(nums, target+1) - 1
	return []int{leftmost, rightmost}
}

func CombinationSum(candidates []int, target int) [][]int {
	if len(candidates) == 0 {
		return [][]int{}
	}
	sort.Ints(candidates)
	c, res := []int{}, [][]int{}
	findComninationSun(candidates, target, 0, c, &res)
	return res
}

func findComninationSun(nums []int, target int, index int, c []int, res *[][]int) {
	if target <= 0 {
		if target == 0 { //到达目的输出解
			b := make([]int, len(c))
			copy(b, c)
			*res = append(*res, b)
		}
		return
	}

	for i := index; i < len(nums); i++ {
		if nums[i] > target { //剪支
			break
		}
		c = append(c, nums[i])
		findComninationSun(nums, target-nums[i], i, c, res)
		c = c[:len(c)-1] // 回溯
	}

}
