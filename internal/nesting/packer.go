package nesting

import (
	"context"
	"fmt"
	"sort"
	"uuid"

	"github.com/edanko/moses/internal/domain"
	"github.com/edanko/moses/pkg/binpack"
)

// PackOptions configures the production packer.
type PackOptions struct {
	Spacing     func(profileType, endcut string) (float64, error)
	StockLength func(profileType string) (float64, error)
	Algorithm   string
}

// PackResult is the output of the production packer.
type PackResult struct {
	Nests       []domain.Nest
	Utilization float64
}

type packingProfile struct {
	profile   domain.Part
	effective float64
}

type groupKey struct {
	project     string
	profileType string
	quality     string
}

// Pack nests profiles into remnants and stock bars. Profiles are expanded by
// their quantity, grouped by project/profile/quality, and packed with the
// configured algorithm (knapsack by default).
func Pack(
	ctx context.Context,
	profiles []domain.Part,
	remnants []domain.Remnant,
	options PackOptions,
) (PackResult, error) {
	var result PackResult
	if options.Spacing == nil || options.StockLength == nil {
		return result, fmt.Errorf("pack: spacing and stock length resolvers are required")
	}
	if len(profiles) == 0 {
		return result, fmt.Errorf("pack: no profiles")
	}

	expand := func(profile domain.Part) int {
		if profile.Quantity() <= 0 {
			return 1
		}
		return profile.Quantity()
	}

	groups := make(map[groupKey][]domain.Part)
	order := make([]groupKey, 0)
	for _, profile := range profiles {
		key := groupKey{
			project:     profile.Project(),
			profileType: profile.DesignationKey(),
			quality:     profile.Grade(),
		}
		if _, exists := groups[key]; !exists {
			order = append(order, key)
		}
		for range expand(profile) {
			partCopy := profile
			partCopy.SetQuantity(1)
			groups[key] = append(groups[key], partCopy)
		}
	}
	sort.Slice(order, func(i, j int) bool {
		if order[i].project != order[j].project {
			return order[i].project < order[j].project
		}
		if order[i].profileType != order[j].profileType {
			return order[i].profileType < order[j].profileType
		}
		return order[i].quality < order[j].quality
	})

	spacing := newSpacingCache(options.Spacing)

	var totalEffective, totalCapacity float64
	for _, key := range order {
		groupProfiles := groups[key]
		nests, effective, capacity, err := packGroup(ctx, key, groupProfiles, remnants, options, spacing)
		if err != nil {
			return result, err
		}
		result.Nests = append(result.Nests, nests...)
		totalEffective += effective
		totalCapacity += capacity
	}

	if totalCapacity > 0 {
		result.Utilization = totalEffective / totalCapacity * 100
	}
	return result, nil
}

func packGroup(
	ctx context.Context,
	key groupKey,
	profiles []domain.Part,
	remnants []domain.Remnant,
	options PackOptions,
	spacing *spacingCache,
) ([]domain.Nest, float64, float64, error) {
	items := make([]packingProfile, 0, len(profiles))
	var totalEffective float64
	for _, profile := range profiles {
		select {
		case <-ctx.Done():
			return nil, 0, 0, ctx.Err()
		default:
		}

		effective, err := effectiveLength(profile, spacing)
		if err != nil {
			return nil, 0, 0, fmt.Errorf(
				"%s %s position %s: %w",
				profile.DesignationName(),
				profile.Grade(),
				profile.PositionNo(),
				err,
			)
		}
		items = append(items, packingProfile{profile: profile, effective: effective})
		totalEffective += effective
	}

	stockLength, err := options.StockLength(key.profileType)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("stock length for %s: %w", key.profileType, err)
	}
	if stockLength <= 0 {
		return nil, 0, 0, fmt.Errorf("stock length for %s must be positive", key.profileType)
	}

	var (
		nests        []domain.Nest
		overItems    []packingProfile
		fitItems     []packingProfile
		capacityUsed float64
	)
	for _, item := range items {
		if item.effective > stockLength {
			overItems = append(overItems, item)
			continue
		}
		fitItems = append(fitItems, item)
	}

	for _, item := range overItems {
		stock := domain.NewStockBar(domain.StockBarParams{
			ID:          uuid.NewV4(),
			Project:     key.project,
			Designation: item.profile.Designation(),
			Grade:       item.profile.Grade(),
			Length:      item.effective,
			Type:        domain.StockTypeOverlength,
		})
		nest := domain.NewNest(domain.NestParams{
			ID:       uuid.NewV4(),
			StockBar: stock,
			Parts:    []domain.Part{item.profile},
		})
		nests = append(nests, *nest)
		capacityUsed += item.effective
	}

	if len(fitItems) > 0 {
		bins := make([]*binpack.Bin[domain.StockBar, domain.Part], 0)
		for _, remnant := range remnants {
			if remnant.Used ||
				remnant.Project != key.project ||
				remnant.Designation.Key() != key.profileType ||
				remnant.Grade != key.quality {
				continue
			}
			stock, err := remnant.ToStockBar()
			if err != nil {
				return nil, 0, 0, err
			}
			bins = append(bins, binpack.NewBin[domain.StockBar, domain.Part](*stock, stock.Length()))
		}

		// One stock bar per remaining part guarantees progress even when the
		// chosen heuristic wastes capacity. Unused bins are dropped below.
		for range fitItems {
			stock := domain.NewStockBar(domain.StockBarParams{
				ID:          uuid.NewV4(),
				Project:     key.project,
				Designation: fitItems[0].profile.Designation(),
				Grade:       key.quality,
				Length:      stockLength,
				Type:        domain.StockTypeBar,
			})
			bins = append(bins, binpack.NewBin[domain.StockBar, domain.Part](stock, stockLength))
		}

		binItems := make([]binpack.Item[domain.Part], len(fitItems))
		for i, item := range fitItems {
			binItems[i] = binpack.Item[domain.Part]{
				Original: item.profile,
				Length:   item.effective,
			}
		}

		algorithm, err := algorithmFor(options.Algorithm)
		if err != nil {
			return nil, 0, 0, err
		}
		packed, err := algorithm.Pack(ctx, binItems, bins)
		if err != nil {
			return nil, 0, 0, err
		}
		if len(packed.UnpackedItems) > 0 {
			return nil, 0, 0, fmt.Errorf(
				"%s %s: %d profiles could not be packed",
				key.profileType,
				key.quality,
				len(packed.UnpackedItems),
			)
		}

		effectiveByID := make(map[uuid.UUID]float64, len(fitItems))
		for _, item := range fitItems {
			effectiveByID[item.profile.ID()] = item.effective
		}

		for _, bin := range packed.PackedBins {
			if len(bin.Items) == 0 {
				continue
			}
			grouped := make([]packingProfile, 0, len(bin.Items))
			for _, item := range bin.Items {
				grouped = append(grouped, packingProfile{
					profile:   item.Original,
					effective: effectiveByID[item.Original.ID()],
				})
			}
			sort.SliceStable(grouped, func(i, j int) bool {
				if grouped[i].effective != grouped[j].effective {
					return grouped[i].effective < grouped[j].effective
				}
				return grouped[i].profile.PositionNo() < grouped[j].profile.PositionNo()
			})

			groupProfiles := make([]domain.Part, len(grouped))
			for i, item := range grouped {
				groupProfiles[i] = item.profile
			}

			nest := domain.NewNest(domain.NestParams{
				ID:       uuid.NewV4(),
				StockBar: bin.Original,
				Parts:    groupProfiles,
				Consumed: bin.Original.ConsumedRemnant(),
			})
			nests = append(nests, *nest)
			capacityUsed += bin.Original.Length()
		}
	}

	sort.SliceStable(nests, func(i, j int) bool {
		leftLength, rightLength := nests[i].StockLength(), nests[j].StockLength()
		if leftLength != rightLength {
			return leftLength < rightLength
		}
		left, right := nests[i].Parts(), nests[j].Parts()
		if len(left) == 0 || len(right) == 0 {
			return len(left) < len(right)
		}
		if left[0].PositionNo() != right[0].PositionNo() {
			return left[0].PositionNo() < right[0].PositionNo()
		}
		return false
	})

	return nests, totalEffective, capacityUsed, nil
}

type spacingResult struct {
	value float64
	err   error
}

type spacingKey struct {
	profileType string
	endcut      string
}

type spacingCache struct {
	resolve func(string, string) (float64, error)
	values  map[spacingKey]spacingResult
}

func newSpacingCache(resolve func(string, string) (float64, error)) *spacingCache {
	return &spacingCache{
		resolve: resolve,
		values:  make(map[spacingKey]spacingResult),
	}
}

func (cache *spacingCache) get(profileType, endcut string) (float64, error) {
	key := spacingKey{profileType: profileType, endcut: endcut}
	if cached, ok := cache.values[key]; ok {
		return cached.value, cached.err
	}
	value, err := cache.resolve(profileType, endcut)
	cache.values[key] = spacingResult{value: value, err: err}
	return value, err
}

func effectiveLength(profile domain.Part, spacing *spacingCache) (float64, error) {
	left, err := spacing.get(profile.DesignationKey(), profile.LeftEndKind())
	if err != nil {
		return 0, fmt.Errorf("left end %q: %w", profile.LeftEndKind(), err)
	}
	right, err := spacing.get(profile.DesignationKey(), profile.RightEndKind())
	if err != nil {
		return 0, fmt.Errorf("right end %q: %w", profile.RightEndKind(), err)
	}
	return left + profile.Length() + right, nil
}

func algorithmFor(name string) (binpack.Algorithm[domain.StockBar, domain.Part], error) {
	switch name {
	case "", "best-fit", "best-fit-decreasing":
		return binpack.NewBestFitDecreasingAlgorithm[domain.StockBar, domain.Part](), nil
	case "knapsack":
		return binpack.NewKnapsackAlgorithm[domain.StockBar, domain.Part](), nil
	case "first-fit":
		return binpack.NewFirstFitAlgorithm[domain.StockBar, domain.Part](), nil
	case "first-fit-decreasing":
		return binpack.NewFirstFitDecreasingAlgorithm[domain.StockBar, domain.Part](), nil
	default:
		return nil, fmt.Errorf("unknown nesting algorithm %q", name)
	}
}
