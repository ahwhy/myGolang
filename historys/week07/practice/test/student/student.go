package main

import "fmt"

type Student struct {
	Name      string
	ChiScore  int
	EngScore  int
	MathScore int
}

func NewStudent(name string) (*Student, error) {
	if name == "" {
		return nil, fmt.Errorf("name为空")
	}

	return &Student{
		Name: name,
	}, nil
}

func (s *Student) GetAvgScore() (int, error) {
	score := s.ChiScore + s.EngScore + s.MathScore
	if score == 0 {
		return 0, fmt.Errorf("全都是0分")
	}

	return score / 3, nil
}
