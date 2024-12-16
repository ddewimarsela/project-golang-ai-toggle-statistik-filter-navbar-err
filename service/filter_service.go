package service

import (
	"a21hc3NpZ25tZW50/model"
)

type FilterService struct{}

func (s *FilterService) FilterData(data map[string][]string, filter model.FilterRequest) (map[string][]string, error) {
	// Inisialisasi map untuk menyimpan hasil filter
	filteredData := make(map[string][]string)
	for key := range data {
		filteredData[key] = []string{}
	}

	// Loop melalui setiap baris data
	for i := range data["Room"] {
		shouldInclude := true

		// Cek filter Room
		if filter.Room != "" && data["Room"][i] != filter.Room {
			shouldInclude = false
		}
		// Cek filter Appliance
		if filter.Appliance != "" && data["Appliance"][i] != filter.Appliance {
			shouldInclude = false
		}
		// Cek filter Date
		if filter.Date != "" && data["Date"][i] != filter.Date {
			shouldInclude = false
		}

		// Jika semua filter cocok, tambahkan data ke hasil
		if shouldInclude {
			for key := range data {
				filteredData[key] = append(filteredData[key], data[key][i])
			}
		}
	}

	return filteredData, nil
}
