# Array
---
# Topics
- insertion : Given a static array of n + 1 size, insert at index idx.
- deletion : Given a static array of n size, delete at index idx.
- reverse : Given a static array of n size, reverse using loop and two pointer. 
- two pointer
- merge : Given two static array, merge them. Given two sorted array, merge them and output array must be sorted.
- sort : Given a static array, sort them using bubble sort, quick sort and merge sort.
- rotation : Given an static array and k, rotate the array by k left or right. Use mod if k >= n, and also use reverse technique.
---

# Quiz
This quizs are from ChatGPT. I will try to add more real interview quizs.
---
```text
Q1. Given two sorted array.
    A = [1, 4, 7, 8, 10]
    B = [2, 3, 9]

    Constraints:
        - Don't use extra arrays
        - Don't use built-in sort
        - Final result must be sorted
        - Time complexity better than O((n + m) log(n + m))
    
    Task: Explain step-by-step how to merge them in place and give
        - Core idea
        - Pointer movement logic
        - Time Complexity
        - Edge cases that break naive solutions.

    Bonus Task:
        - If you use Gap Method, then explain how is the initial gap computed and how does it shrink.
```

```text
Q2. Given an array. 
    A = [1, 2, 3, 4, 5, 6, 7]
    You must reverse only the EVEN numbers in place, keeping ODD numbers fixed in position.

    Constraints:
        - O(n) time
        - O(1) space
        - No extra array

    Explain:
        - Pointer initialization
        - Pointer movement rules
        - Termination condition
```

```text
Q3. Given an array, detect if it is sorted in:
    - strictly ascending
    - strictly descending
    - not sorted

    Constraints:
        - One pass
        - No sorting
        - Handle duplicates correctly

    Explain the flags / invariants
```

```text
Q4. You are given an array of distinct integers.
    A = [4, 5, 6, 7, 1, 2, 3]
    This is a rotation of a strictly ascending sorted array.

    Task: In one pass, determine whether the array is:
        - rotated sorted ascending
        - strictly sorted ascending (no rotation)
        - not rotated sorted

    Constraints:
        - O(n) time
        - O(1) space
        - No sorting
        - No extra array

    Explain
        - What invariant you track
        - What condition invalidates the array
        - How you distinguish pure sorted vs rotated sorted
```

```text
Q5. You are given an array of size n+1 containing integers in the range [1, n].
    A = [3, 1, 3, 4, 2]

    Constraints:
        - Exactly one number is duplicated
        - It may appear multiple times
        - Array is read-only
        - No extra space
        - Do not modify the array

    Task:
        - Explain the algorithm to find the duplicate.

    You must state:
        - The core idea (mapping / invariant)
        - Why it works mathematically
        - Time complexity
        - Space complexity

    Bonus Question:
        - Can this problem be solved using binary search on value range without modifying the array?
```