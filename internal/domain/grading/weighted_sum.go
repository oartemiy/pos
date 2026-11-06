package grading

// WeightedSum — взвешенная сумма баллов без нормировки:
//
//	итог = Σ(вес · балл)
//
// Подходит для формул вида «0.4 · ДЗ + 0.6 · Экзамен», где все критерии
// уже в одной шкале, а веса — коэффициенты формулы.
type WeightedSum struct{}

func (WeightedSum) Name() string { return "взвешенная сумма" }

func (WeightedSum) Calculate(lines []Line) float64 {
	var sum float64
	for _, l := range lines {
		if l.Graded {
			sum += l.Criterion.Weight().Value() * l.Score
		}
	}
	return sum
}
