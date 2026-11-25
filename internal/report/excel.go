package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/xuri/excelize/v2"
)

const (
	barSheet  = "Ведомость расхода материала"
	partSheet = "Партлист"
)

// Excel writes the two-sheet bar and part workbook.
func (r *Reporter) Excel(path string, bars []Bar) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create report directory: %w", err)
	}

	book := excelize.NewFile()
	defer func() { _ = book.Close() }()

	if err := book.SetSheetName("Sheet1", barSheet); err != nil {
		return err
	}
	if err := book.SetSheetRow(barSheet, "A1", &[]any{
		"Карта раскроя",
		"Тип оборудования",
		"Исп. ДМО",
		"Типоразмер",
		"Марка материала",
		"Длина, мм",
		"Масса заготовки, кг",
		"Коэф. раскроя, %",
		"Маркировка ДМО",
		"Длина ДМО, мм",
		"Масса ДМО, кг",
	}); err != nil {
		return err
	}

	row := 2
	for _, bar := range bars {
		stock := bar.Nest.StockBar()
		profileType := stock.DesignationKey()
		used, err := r.usedLength(bar)
		if err != nil {
			return fmt.Errorf("mass for %q: %w", bar.Name, err)
		}
		stockMass, err := mass(profileType, stock.Length(), r.Spec)
		if err != nil {
			return fmt.Errorf("mass for %q: %w", bar.Name, err)
		}

		consumed := ""
		if bar.Consumed != nil {
			consumed = bar.Consumed.Marking
		}

		barRow := []any{
			bar.Name,
			"Камера резки профиля",
			consumed,
			stock.DesignationName(),
			stock.Grade(),
			stock.Length(),
			stockMass,
			percent(used, stock.Length()),
			"",
			"",
			"",
		}
		if bar.Created != nil {
			createdMass, err := mass(profileType, bar.Created.Length, r.Spec)
			if err != nil {
				return fmt.Errorf("mass for %q: %w", bar.Name, err)
			}
			barRow[8] = bar.Created.Marking
			barRow[9] = bar.Created.Length
			barRow[10] = createdMass
		}
		cell := "A" + strconv.Itoa(row)
		if err := book.SetSheetRow(barSheet, cell, &barRow); err != nil {
			return err
		}
		row++
	}

	if _, err := book.NewSheet(partSheet); err != nil {
		return err
	}
	if err := book.SetSheetRow(partSheet, "A1", &[]any{
		"Чертеж",
		"Заказ",
		"Секция",
		"Позиция",
		"Карта раскроя",
		"Типоразмер детали",
		"Марка материала",
		"Кол-во, шт",
		"Длина, мм",
		"Масса 1 детали, кг",
		"Общая масса, кг",
		"Маршрут обработки",
		"Примечание",
	}); err != nil {
		return err
	}

	row = 2
	for _, bar := range bars {
		for _, profile := range bar.Nest.Parts() {
			profileType := profile.DesignationKey()
			partMass, err := mass(profileType, profile.Length(), r.Spec)
			if err != nil {
				return fmt.Errorf("mass for %q: %w", bar.Name, err)
			}
			position := 0
			if parsed, err := strconv.Atoi(profile.PositionNo()); err == nil {
				position = parsed
			}
			partRow := []any{
				"",
				profile.Project(),
				profile.Section(),
				position,
				bar.Name,
				profile.DesignationName(),
				profile.Grade(),
				profile.Quantity(),
				profile.Length(),
				partMass,
				partMass * float64(profile.Quantity()),
				"",
				"",
			}
			cell := "A" + strconv.Itoa(row)
			if err := book.SetSheetRow(partSheet, cell, &partRow); err != nil {
				return err
			}
			row++
		}
	}

	index, err := book.GetSheetIndex(partSheet)
	if err != nil {
		return err
	}
	book.SetActiveSheet(index)
	if err := book.SaveAs(path); err != nil {
		return fmt.Errorf("save report %q: %w", path, err)
	}
	return nil
}
