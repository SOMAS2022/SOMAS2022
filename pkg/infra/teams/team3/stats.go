package team3

type Stats struct {
	Health  uint
	Stamina uint
	Attack  uint
	Shield  uint
}

type AveragedStats struct {
	Health  float64
	Stamina float64
	Attack  float64
	Shield  float64
}

type StatsQueue struct {
	MaxQueueLen float64
	Queue       []Stats
}

func makeStatsQueue(len float64) *StatsQueue {
	return &StatsQueue{MaxQueueLen: len, Queue: make([]Stats, int(len))}
}

func (sq *StatsQueue) addStat(stat Stats) {
	sq.Queue = append(sq.Queue, stat)
	if len(sq.Queue) > int(sq.MaxQueueLen) {
		sq.Queue = sq.Queue[1:]
	}
}

func (sq *StatsQueue) averageStats() AveragedStats {
	avgStats := AveragedStats{Health: 0, Stamina: 0, Attack: 0, Shield: 0}
	for _, stat := range sq.Queue {
		avgStats.Health += float64(stat.Health)
		avgStats.Stamina += float64(stat.Stamina)
		avgStats.Attack += float64(stat.Attack)
		avgStats.Shield += float64(stat.Shield)
	}
	avgStats.Health /= sq.MaxQueueLen
	avgStats.Stamina /= sq.MaxQueueLen
	avgStats.Attack /= sq.MaxQueueLen
	avgStats.Shield /= sq.MaxQueueLen

	return avgStats
}
