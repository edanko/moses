package nest

import (
	"fmt"
	"math"
	"sort"

	"github.com/edanko/moses/internal/config"
	"github.com/edanko/moses/internal/models"
	"github.com/edanko/moses/internal/remnant"

	"github.com/spf13/viper"
)

func checkPart(b *models.Bar, i int, w int, ps []*models.Part, matrix [][]int) {
	if i <= 0 || w <= 0 {
		return
	}

	pick := matrix[i][w]
	if pick != matrix[i-1][w] {
		b.PutPart(ps[i-1])
		checkPart(b, i-1, w-int(math.Ceil(ps[i-1].FullLength)), ps, matrix)
	} else {
		checkPart(b, i-1, w, ps, matrix)
	}
}

// Nest - dynamic programming version, columns method
// load remnants first, then add needed amount of stock bars
func (n *Nest) Nest() {
	fmt.Println("[i]", n.Parts[0].Dim, "-", len(n.Parts), "parts")
	p := n.Parts[0]
	rems := remnant.Remnants(p.Project, p.Dim, p.Quality)
	for _, rem := range rems {
		n.AddBar(rem)
	}
	sort.Sort(models.BarSlice(n.Bars))

	for i := 0; len(n.Parts) != 0; i++ {
		// skip parts bigger than bar length from n.Parts
		// TODO: move this somwhere else, preferable to load drawings
		if n.Parts[0].FullLength > config.BarSize(n.Parts[0].Dim) {
			fmt.Println(n.Parts[0].Project + "-" + n.Parts[0].Section + "-" + n.Parts[0].PosNo + "(len: " + n.Parts[0].LengthString() + ") > bar len - " + n.Parts[0].Dim)

			// unfit current part
			//n.remove([]*models.Part{n.Parts[0]})
			//i--
			//continue

			fmt.Println("[i] adding overlenged bar")
			overLenBar := models.NewBar(n.Parts[0].FullLength)
			overLenBar.PutPart(n.Parts[0])
			n.remove([]*models.Part{n.Parts[0]})
			n.AddBar(overLenBar)
			i--
			continue

		}

		// add new stock bar
		if i >= len(n.Bars) {
			n.AddBar(models.NewBar(config.BarSize(n.Parts[0].Dim)))
		}

		b := n.Bars[i]
		ps := n.Parts

		numItems := len(ps)                    // number of parts
		capacity := int(math.Ceil(b.Capacity)) // capacity of current bar

		if capacity == int(math.Ceil(b.Length)) {
			continue
		}

		// create the empty matrix
		matrix := make([][]int, numItems+1) // rows representing parts
		for i := range matrix {
			matrix[i] = make([]int, capacity+1) // columns representing millimeters of length
		}

		// loop through table rows
		for i := 1; i <= numItems; i++ {
			// loop through table columns
			for w := 1; w <= capacity; w++ {
				// if weight of part matching this index can fit at the current capacity column...
				if int(math.Ceil(ps[i-1].FullLength)) <= w {
					// length of this subset without this part
					valueOne := float64(matrix[i-1][w])
					// length of this subset without the previous part, and this part instead
					valueTwo := float64(int(math.Ceil(ps[i-1].FullLength)) + matrix[i-1][w-int(math.Ceil(ps[i-1].FullLength))])
					// take maximum of either valueOne or valueTwo
					matrix[i][w] = int(math.Max(valueOne, valueTwo))
					// if the new length is not more, carry over the previous length
				} else {
					matrix[i][w] = matrix[i-1][w]
				}
			}
		}

		checkPart(b, numItems, capacity, ps, matrix)
		// add other statistics to the bag
		b.Length = float64(matrix[numItems][capacity])
		// remove nested parts from n.Parts
		n.remove(b.Parts)
		// sort from small to big
		sort.Sort(models.PartSlice(n.Bars[i].Parts))
	}

	for i := 0; i < len(n.Bars); i++ {
		n.Bars[i].ID = i + 1
	}

	// skip remnant bars, without nested parts
	bars := make([]*models.Bar, 0)
	for _, b := range n.Bars {
		if len(b.Parts) > 0 {
			bars = append(bars, b)
		}
	}
	n.Bars = bars
	sort.Sort(models.BarSlice(n.Bars))

	for i := 0; i < len(n.Bars); i++ {
		n.Bars[i].ID = i + 1

		remnant.Use(n.Bars[i].UsedRemnantID, n.Bars[i].NestName())

		if n.Bars[i].RemnantLength() > viper.GetFloat64("mincut") {
			n.Bars[i].MarkRemnant = n.Bars[i].NestName() + "-L" + n.Bars[i].RemnantLengthString()

			b := n.Bars[i]

			remnant.Add(b.Project(), b.Section(), b.Dim(), b.Quality(), b.RemnantLengthString(), b.MarkRemnant)
		}
	}
}

func (n *Nest) remove(parts []*models.Part) {
	for _, partToDel := range parts {
		i := 0
		for _, part := range n.Parts {
			if partToDel.PosNo == part.PosNo {
				n.Parts[i] = n.Parts[len(n.Parts)-1]
				n.Parts[len(n.Parts)-1] = &models.Part{}
				n.Parts = n.Parts[:len(n.Parts)-1]
				break
			}
			i++
		}
	}
}
