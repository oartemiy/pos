package grading

// WeightedAverage — взвешенное среднее по 10-балльной шкале:
//
//	итог = 10 · Σ(вес · балл/максимум) / Σ(вес)
//
// Балл нормируется к максимуму критерия, чтобы критерий «из 100»
// не перевешивал критерий «из 5» сильнее, чем задано весом.
// Неоценённый критерий считается как 0: итог не завышается,
// если выставлены оценки только по части критериев.
type WeightedAverage struct{}

func (WeightedAverage) Name() string { return "взвешенное среднее" }

func (WeightedAverage) Calculate(lines []Line) float64 {
	var weighted, totalWeight float64
	for _, l := range lines {
		w := l.Criterion.Weight().Value()
		totalWeight += w
		if l.Graded {
			weighted += w * l.Score / l.Criterion.MaxScore().Value()
		}
	}
	if totalWeight == 0 {
		return 0
	}
	return 10 * weighted / totalWeight
}
