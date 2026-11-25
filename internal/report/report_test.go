package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"uuid"

	"github.com/xuri/excelize/v2"

	"github.com/edanko/moses/internal/domain"
)

func reporter() *Reporter {
	return &Reporter{
		Spacing: func(_, _ string) (float64, error) { return 6, nil },
		Spec: func(string) (domain.SectionGeometry, error) {
			return domain.SectionGeometry{
				Shape:        domain.ShapeHP,
				B:            160,
				C:            28,
				S:            8,
				R:            7,
				MassPerMeter: 14.08,
			}, nil
		},
		MinRemnantLength: 1100,
	}
}

func bar(t *testing.T) Bar {
	t.Helper()
	profileType, err := domain.ParseProfileDesignation("RP160X8")
	if err != nil {
		t.Fatal(err)
	}
	end, err := domain.ParseEndCut("21 v1=0 v2=0")
	if err != nil {
		t.Fatal(err)
	}
	profile := domain.NewPart(domain.PartParams{
		ID:          uuid.NewV4(),
		Project:     "056001",
		Section:     "05102",
		PositionNo:  "4835",
		Designation: profileType,
		Grade:       "A40",
		Length:      6000,
		Quantity:    1,
		LeftEndCut:  end,
		RightEndCut: end,
	})
	stock := domain.NewStockBar(domain.StockBarParams{
		ID:          uuid.NewV4(),
		Project:     "056001",
		Designation: profileType,
		Grade:       "A40",
		Length:      12000,
		Type:        domain.StockTypeBar,
	})
	nest := domain.NewNest(domain.NestParams{
		ID:       uuid.NewV4(),
		StockBar: stock,
		Parts:    []domain.Part{profile},
	})
	return Bar{
		Nest:  *nest,
		Index: 1,
		Name:  "05102-RP160X8-A40-1",
		Created: &domain.Remnant{
			Length:  5988,
			Marking: "05102-RP160X8-A40-1-L5988",
		},
	}
}

func TestBarList(t *testing.T) {
	got := reporter().BarList([]Bar{bar(t)})
	want := "   RP160X8 |   A40 |  1 bars (12000.000000)\r\n"
	if got != want {
		t.Fatalf("bar list = %q, want %q", got, want)
	}
}

func TestNestingList(t *testing.T) {
	got, err := reporter().NestingList([]Bar{bar(t)})
	if err != nil {
		t.Fatalf("nesting list: %v", err)
	}
	for _, expected := range []string{
		"Имя программы      05102-RP160X8-A40-1\r\n",
		"Длина: 12000 / Исп.: 6000",
		"Отход: 5988",
		"Маркировать отход  05102-RP160X8-A40-1-L5988",
		"|  05102 |    4835 |  6000 |",
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("nesting list missing %q:\n%s", expected, got)
		}
	}
}

func TestTxt(t *testing.T) {
	got := Txt([]Bar{bar(t)}, false)
	for _, expected := range []string{
		"12000\r\n\r\n4835\r\n6000\r\n21 v1=0 v2=0\r\n21 v1=0 v2=0\r\n",
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("txt missing %q:\n%s", expected, got)
		}
	}
	name, err := TxtFileName(bar(t))
	if err != nil {
		t.Fatal(err)
	}
	if name != "056001-05102_RP160X8_A40" {
		t.Fatalf("name = %q", name)
	}
}

func TestExcel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.xlsx")
	if err := reporter().Excel(path, []Bar{bar(t)}); err != nil {
		t.Fatalf("excel: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	book, err := excelize.OpenFile(path)
	if err != nil {
		t.Fatalf("open excel: %v", err)
	}
	defer func() { _ = book.Close() }()

	barRows, err := book.GetRows(barSheet)
	if err != nil {
		t.Fatal(err)
	}
	if len(barRows) != 2 || barRows[1][0] != "05102-RP160X8-A40-1" || barRows[1][3] != "RP160X8" {
		t.Fatalf("unexpected bar rows: %+v", barRows)
	}

	partRows, err := book.GetRows(partSheet)
	if err != nil {
		t.Fatal(err)
	}
	if len(partRows) != 2 || partRows[1][4] != "05102-RP160X8-A40-1" || partRows[1][8] != "6000" {
		t.Fatalf("unexpected part rows: %+v", partRows)
	}
}
