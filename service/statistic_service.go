package service

import (
	"fmt"
	"sort"
	"strconv"
)

type StatisticService struct{}

func (s *StatisticService) CalculateStatistics(data map[string][]string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Mengambil data konsumsi energi
	consumptions, ok := data["Energy_Consumption"]
	if !ok {
		return nil, fmt.Errorf("kolom Energy_Consumption tidak ditemukan")
	}

	// Konversi string ke float64
	values := make([]float64, 0)
	var total float64
	for _, v := range consumptions {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("error konversi data: %v", err)
		}
		values = append(values, val)
		total += val
	}

	// Hitung statistik dasar
	n := float64(len(values))
	mean := total / n

	// Sort untuk mendapatkan min, max, dan median
	sort.Float64s(values)
	min := values[0]
	max := values[len(values)-1]

	// Hitung median
	var median float64
	if len(values)%2 == 0 {
		median = (values[len(values)/2-1] + values[len(values)/2]) / 2
	} else {
		median = values[len(values)/2]
	}

	// Masukkan hasil ke map
	stats["total_konsumsi"] = total
	stats["rata_rata"] = mean
	stats["minimum"] = min
	stats["maximum"] = max
	stats["median"] = median

	return stats, nil
}
