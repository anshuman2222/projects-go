package main

import (
	"fmt"
)

type Stack struct {
	Top   int
	Cap   int
	Array []int
}

func (st *Stack) is_stack_overflow() bool {
	if st.Top == st.Cap {
		fmt.Println("Stack is overflow")
		return true
	} else {
		fmt.Println("Stack is not overflow")
		return false
	}
}

func (st *Stack) is_stack_underflow() bool {
	if st.Top == -1 {
		fmt.Println("Stack is underflow")
		return true
	} else {
		return false
	}
}

func (st *Stack) push(ele int) {
	if st.is_stack_overflow() == true {
		fmt.Println("Stack is overflow")
		return
	}
	st.Top++
	fmt.Println("Entering elemen : at index: ", ele, st.Top)
	st.Array[st.Top] = ele
}

func (st *Stack) top() int {
	fmt.Println("Top element: ", st.Array[st.Top])
	return st.Array[st.Top]
}

func (st *Stack) pop() {
	if st.is_stack_underflow() == true {
		fmt.Println("Stack is underflow")
		return
	}
	st.Top--
}

func main() {
	var cap, no, ele int

	fmt.Println("Enter capacity and number: ")
	fmt.Scanf("%d%d", &cap, &no)

	st := &Stack{Top: -1, Cap: cap, Array: nil}
	st.Array = make([]int, cap)

	for i := 0; i < no; i++ {
		fmt.Scanf("%d", &ele)

		st.push(ele)
	}
	fmt.Println("Top ele: ", st.top())
	st.pop()
	fmt.Println("Top ele: ", st.top())
}
