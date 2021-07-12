package nester

import (
	"math"

	"github.com/edanko/moses/entities"
)

func checkPart(b *entities.Nest, i int, w int, ps []*entities.Profile, matrix [][]int) {
	if i <= 0 || w <= 0 {
		return
	}

	pick := matrix[i][w]
	if pick != matrix[i-1][w] {

		b.Profiles = append(b.Profiles, ps[i-1])
		b.Bar.UsedLength += ps[i-1].FullLength
		//b.PutPart(ps[i-1])
		checkPart(b, i-1, w-int(math.Ceil(ps[i-1].Length)), ps, matrix)
	} else {
		checkPart(b, i-1, w, ps, matrix)
	}
}

func Nest(remnants []*entities.Bar, stockBar *entities.Bar, ps []*entities.Profile) ([]*entities.Nest, error) {
	nests := make([]*entities.Nest, 0, len(remnants))

	var b *entities.Bar

	for i := 0; len(ps) != 0; i++ {
		if i >= len(remnants) {
			b = stockBar
		} else {
			b = remnants[i]
		}

		numProfiles := len(ps)
		barLength := int(math.Ceil(b.Length))

		matrix := make([][]int, numProfiles+1) // rows representing parts
		for i := range matrix {
			matrix[i] = make([]int, barLength+1) // columns representing millimeters of length
		}

		// loop through table rows
		for i := 1; i <= numProfiles; i++ {
			// loop through table columns
			for w := 1; w <= barLength; w++ {
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

		n := &entities.Nest{}
		n.Bar = b
		checkPart(n, numProfiles, barLength, ps, matrix)

		if len(n.Profiles) == 0 {
			continue
		}

		n.Bar.UsedLength = float64(matrix[numProfiles][barLength])
		ps = removeNestedParts(ps, n.Profiles)

		//sort.Sort(entities.ProfileSlice(n.Profiles))
		nests = append(nests, n)
	}
	return nests, nil
}

func removeNestedParts(parts, del []*entities.Profile) []*entities.Profile {
	for _, d := range del {
		i := 0
		for _, p := range parts {
			if d.Project == p.Project && d.Section == p.Section && d.PosNo == p.PosNo {
				if p.Quantity > 1 {
					p.Quantity--
				} else {
					parts[i] = parts[len(parts)-1]
					parts[len(parts)-1] = &entities.Profile{}
					parts = parts[:len(parts)-1]
				}
				break
			}
			i++
		}
	}
	return parts
}
