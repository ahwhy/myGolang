package collector

import (
	"encoding/json"
	"strconv"
)


func ParseLine(line string) *RocketMQMetric {
	words := []string{}
	chars := []rune{}
	for _, c := range line {
		if c == ' ' {
			if len(chars) > 0 {
				words = append(words, string(chars))
			}
			chars = []rune{}
		} else {
			chars = append(chars, c)
		}
	}

	if len(chars) > 0 {
		words = append(words, string(chars))
	}

	m := &RocketMQMetric{
		Group:   words[0],
		Count:   words[1],
		Version: words[2],
		Type:    words[3],
	}

	if len(words) <= 6 {
		m.TPS = words[4]
		m.DiffTotal = words[5]
	} else {
		m.Model = words[4]
		m.TPS = words[5]
		m.DiffTotal = words[6]
	}

	return m
}

type RocketMQMetric struct {
	Group     string
	Count     string
	Version   string
	Type      string
	Model     string
	TPS       string
	DiffTotal string
}

func (m *RocketMQMetric) IntCount() float64 {
	i, _ := strconv.ParseFloat(m.Count, 64)

	return i
}

func (m *RocketMQMetric) IntTPS() float64 {
	i, _ := strconv.ParseFloat(m.TPS, 64)

	return i
}

func (m *RocketMQMetric) IntDiffTotal() float64 {
	i, _ := strconv.ParseFloat(m.DiffTotal, 64)

	return i
}

func (m *RocketMQMetric) String() string {
	v, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(v)
}

func NewRocketMQMetricSet() *RocketMQMetricSet {
	return &RocketMQMetricSet{
		Items: []*RocketMQMetric{},
	}
}

type RocketMQMetricSet struct {
	Items []*RocketMQMetric
}

func (s *RocketMQMetricSet) Add(item *RocketMQMetric) {
	s.Items = append(s.Items, item)
}

func (s *RocketMQMetricSet) String() string {
	v, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(v)
}
