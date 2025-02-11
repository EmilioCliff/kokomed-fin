package reports

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type excelGenerator struct {
    file     *excelize.File
    currentSheet string
}

func newExcelGenerator() *excelGenerator {
    return &excelGenerator{
        file: excelize.NewFile(),
    }
}

func (e *excelGenerator) createSheet(name string) {
    e.file.NewSheet(name)
    e.currentSheet = name
}

func (e *excelGenerator) closeExcel() error {
	return e.file.Close()
}

func (e *excelGenerator) writeHeader(columns []string, styleID int) {
	end, _ := excelize.CoordinatesToCellName(len(columns), 1) // lastHeader1
	e.file.SetCellStyle(e.currentSheet, "A1", end, styleID)

    for i, col := range columns {
		cell := fmt.Sprintf("%s1", string(rune(65+i)))
        e.file.SetCellValue(e.currentSheet, cell, col)
    }
}

func (e *excelGenerator) writeRow(row int, data []interface{}) {
	for i, val := range data {
		cell := fmt.Sprintf("%s%d", string(rune(65+i)), row)
		e.file.SetCellValue(e.currentSheet, cell, val)
	}
}

func (e *excelGenerator) createHeaderStyle() int {
    // style, _ := e.file.NewStyle(&excelize.Style{
    //     Font:      &excelize.Font{Bold: true},
    //     Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DFF2FF"}, Pattern: 1},
    //     Alignment: &excelize.Alignment{Horizontal: "center"},
    // })
	headerStyle, _ := e.file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "center",
			Indent:          1,
			JustifyLastLine: true,
			ReadingOrder:    0,
			RelativeIndent:  1,
			ShrinkToFit:     true,
			Vertical:        "top",
			WrapText:        true,
		},
		Font: &excelize.Font{
			Bold:   true,
			Italic: true,
			Family: "Times New Roman",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#0094B7"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
    return headerStyle
}

func (e *excelGenerator) createDateStyle() int {
	dateStyle, _ := e.file.NewStyle(&excelize.Style{
		NumFmt: 14,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})
    return dateStyle
}

func (e *excelGenerator) createPercentageStyle() int {
	percentageStyle, _ := e.file.NewStyle(&excelize.Style{
		NumFmt: 9,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})
	return percentageStyle
}

func (e *excelGenerator) createMoneyStyle() int {
	moneyStyle, _ := e.file.NewStyle(&excelize.Style{
		NumFmt: 4,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})
    return moneyStyle
}

func (e *excelGenerator) createQuantityStyle() int {
	quantityStyle, _ := e.file.NewStyle(&excelize.Style{
		NumFmt: 3,
	})
    return quantityStyle
}