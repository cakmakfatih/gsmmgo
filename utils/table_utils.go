package utils

import "echochat/models"

type TableData struct {
	SizeClasses []string
	Header      []string
	Rows        [][]string
}

func ProvidersToTableData(providers []models.ProviderModel) TableData {
	var result TableData

	result.SizeClasses = []string{"w-56", "w-52", "w-36", "w-56", "w-56"}
	result.Header = []string{"ID", "URL", "Method", "Method Data", "Alias"}

	for _, p := range providers {
		result.Rows = append(result.Rows, []string{
			p.ID,
			p.URL,
			p.Method,
			p.MethodDataReadable,
			p.Alias,
		})
	}

	return result
}
