package delivery

import (
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
)

const (
	scoreStreet   = 40
	scoreLandmark = 30
	scoreCity     = 20
	scoreRegion   = 10
	scoreCountry  = 5
)

func containsFold(haystack, needle string) bool {
	if needle == "" || haystack == "" {
		return false
	}
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

func fieldsMatch(target, field string) bool {
	if field == "" {
		return false
	}
	if target == "" {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(target), strings.TrimSpace(field)) ||
		containsFold(target, field) ||
		containsFold(field, target)
}

func matchLocation(street, city, region, country string, loc db.DeliveryRouteLocation, storeLocation string) (bool, int, string) {
	if loc.UseStoreLocation {
		if storeLocation != "" {
			matched, score, _ := matchLocation(storeLocation, "", "", "", loc, "")
			if matched {
				return true, score, "store location"
			}
		}
		return true, scoreCity, "store location"
	}

	matched := false
	score := 0
	var parts []string

	if loc.Street != "" && fieldsMatch(street, loc.Street) {
		score += scoreStreet
		matched = true
		parts = append(parts, loc.Street)
	}
	if loc.Landmark != "" && (fieldsMatch(street, loc.Landmark) || fieldsMatch(city, loc.Landmark)) {
		score += scoreLandmark
		matched = true
		parts = append(parts, loc.Landmark)
	}
	if loc.City != "" && fieldsMatch(city, loc.City) {
		score += scoreCity
		matched = true
		parts = append(parts, loc.City)
	}
	if loc.Region != "" && fieldsMatch(region, loc.Region) {
		score += scoreRegion
		matched = true
		parts = append(parts, loc.Region)
	}
	if !matched && loc.Country != "" && loc.Street == "" && loc.City == "" && loc.Region == "" {
		if country == "" || fieldsMatch(country, loc.Country) {
			score += scoreCountry
			matched = true
			parts = append(parts, loc.Country)
		}
	}

	summary := strings.Join(parts, " / ")
	if loc.Label != "" && summary == "" {
		summary = loc.Label
	}
	return matched, score, summary
}

type routeMatch struct {
	RouteID        int64
	PickupScore    int
	DeliveryScore  int
	PickupSummary  string
	DeliverySummary string
}

func scoreRoute(route db.DeliveryRoute, storeLocation string, pickupCity, pickupRegion string,
	deliveryStreet, deliveryCity, deliveryRegion, deliveryCountry string) (bool, routeMatch) {
	var pickupMatched bool
	var pickupScore int
	var pickupSummary string

	for _, loc := range route.Locations {
		if loc.LocationType != constant.DeliveryLocationTypePickup {
			continue
		}
		street := ""
		city := pickupCity
		region := pickupRegion
		if loc.UseStoreLocation {
			m, s, sum := matchLocation(street, city, region, "", loc, storeLocation)
			if m {
				pickupMatched = true
				if s > pickupScore {
					pickupScore = s
					pickupSummary = sum
				}
			}
			continue
		}
		m, s, sum := matchLocation(loc.Street, loc.City, loc.Region, loc.Country, loc, storeLocation)
		if m {
			pickupMatched = true
			if s > pickupScore {
				pickupScore = s
				pickupSummary = sum
			}
		}
		// Also try matching store location against route pickup fields
		if storeLocation != "" {
			m2, s2, sum2 := matchLocation(storeLocation, city, region, "", loc, storeLocation)
			if m2 && s2 > pickupScore {
				pickupMatched = true
				pickupScore = s2
				pickupSummary = sum2
			}
		}
	}

	if !pickupMatched {
		return false, routeMatch{}
	}

	var deliveryMatched bool
	var deliveryScore int
	var deliverySummary string

	for _, loc := range route.Locations {
		if loc.LocationType != constant.DeliveryLocationTypeDelivery {
			continue
		}
		m, s, sum := matchLocation(deliveryStreet, deliveryCity, deliveryRegion, deliveryCountry, loc, storeLocation)
		if m {
			deliveryMatched = true
			if s > deliveryScore {
				deliveryScore = s
				deliverySummary = sum
			}
		}
	}

	if !deliveryMatched {
		return false, routeMatch{}
	}

	return true, routeMatch{
		RouteID:         route.ID,
		PickupScore:     pickupScore,
		DeliveryScore:   deliveryScore,
		PickupSummary:   pickupSummary,
		DeliverySummary: deliverySummary,
	}
}
