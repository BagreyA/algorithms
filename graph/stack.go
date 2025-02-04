package graph

type Stack struct {
    data []int
}

func (s *Stack) Push(x int) {
    s.data = append(s.data, x)
}

func (s *Stack) Pop() (int, bool) {
    if len(s.data) == 0 {
        return 0, false
    }
    x := s.data[len(s.data)-1]
    s.data = s.data[:len(s.data)-1]
    return x, true
}

func (s *Stack) IsEmpty() bool {
    return len(s.data) == 0
}
