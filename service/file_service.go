package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	// Validasi input
	if fileContent == "" {
		return nil, errors.New("file content kosong")
	}

	// Buat CSV reader dari string content
	reader := csv.NewReader(strings.NewReader(fileContent))

	// Baca header
	headers, err := reader.Read()
	if err != nil {
		return nil, errors.New("gagal membaca header CSV")
	}

	// Validasi header yang dibutuhkan
	requiredHeaders := []string{"Date", "Time", "Appliance", "Energy_Consumption", "Room", "Status"}
	headerMap := make(map[string]bool)
	for _, h := range headers {
		headerMap[h] = true
	}
	for _, required := range requiredHeaders {
		if !headerMap[required] {
			return nil, errors.New("header yang dibutuhkan tidak ditemukan: " + required)
		}
	}

	//proses data csv
	// Inisialisasi map untuk menyimpan data
	result := make(map[string][]string)
	for _, header := range headers {
		result[header] = []string{}
	}

	// Baca data baris per baris
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("gagal membaca baris data: " + err.Error())
		}

		// Validasi jumlah kolom
		if len(record) != len(headers) {
			return nil, errors.New("jumlah kolom tidak sesuai dengan header")
		}

		// Masukkan data ke map
		for i, value := range record {
			result[headers[i]] = append(result[headers[i]], value)
		}
	}

	// Validasi hasil akhir
	if len(result["Appliance"]) == 0 {
		return nil, errors.New("tidak ada data yang valid dalam file")
	}

	return result, nil
}
