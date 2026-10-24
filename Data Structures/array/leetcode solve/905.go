// https://leetcode.com/problems/sort-array-by-parity/description/

// use two pointer to solve.
// one for odd, one for even
// O(n) time complexity
// O(1) memory complexity
func sortArrayByParity(nums []int) []int {
	odd, even := 0, len(nums)-1
	for odd < even {
		for nums[odd]%2 == 0 && odd < even { // skip evens
			odd++
		}
		for nums[even]%2 == 1 && odd < even { // skip odds
			even--
		}
		nums[odd], nums[even] = nums[even], nums[odd]
		odd, even = odd+1, even-1
	}
	return nums
}