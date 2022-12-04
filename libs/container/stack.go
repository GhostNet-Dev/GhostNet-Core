package container

type Stack []interface{}

func (s *Stack) Count() int {
	return len(*s)
}

// IsEmpty - 스택이 비어있는지 확인하는 함수
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push - 스택에 값을 추가하는 함수.
func (s *Stack) Push(data interface{}) {
	*s = append(*s, data)
}

// Pop - 스택에 값을 제거하고 top위치에 값을 반환하는 함수.
func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		return nil
	} else {
		top := len(*s) - 1
		data := (*s)[top]
		*s = (*s)[:top]
		return data
	}
}

func (s *Stack) Peek() interface{} {
	if s.IsEmpty() {
		return nil
	} else {
		top := len(*s) - 1
		data := (*s)[top]
		return data
	}
}

func (s *Stack) Clear() {
	*s = (*s)[:0]
}
