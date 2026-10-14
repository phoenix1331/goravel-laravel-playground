// Package factories holds model factories. A model links to its factory
// with a Factory() method, the equivalent of the HasFactory trait.
package factories

import (
	"github.com/brianvoe/gofakeit/v7"
)

// towns gives the demo league a grassroots English feel; gofakeit's own
// place names are American.
var towns = []string{
	"Ashford", "Barnsley", "Crewe", "Dorking", "Eastleigh", "Frome", "Gosport",
	"Hitchin", "Ilkley", "Kendal", "Ludlow", "Malton", "Newbury", "Oakham",
	"Penrith", "Ripon", "Stroud", "Thirsk", "Ulverston", "Whitby",
}

type ClubFactory struct{}

// Definition returns default values keyed by struct field name (not column
// name). Any key can be overridden when the factory is called.
//
// Laravel:
//
//	public function definition(): array
//	{
//	    return ['name' => fake()->city().' '.fake()->randomElement([...]), ...];
//	}
//
//	Club::factory()->count(3)->create();
//
// Goravel: facades.Orm().Factory().Count(3).Create(&clubs). gofakeit stands in
// for Faker. Keys are field names ("FoundedYear", not "founded_year"), and a
// misspelt key is silently ignored.
func (f *ClubFactory) Definition() map[string]any {
	town := gofakeit.RandomString(towns)
	suffix := gofakeit.RandomString([]string{"Athletic", "Rovers", "United", "Town", "Wanderers"})
	year := gofakeit.IntRange(1870, 2020)

	return map[string]any{
		"Name":        town + " " + suffix,
		"Ground":      town + " " + gofakeit.RandomString([]string{"Park", "Rec", "Lane", "Road"}),
		"FoundedYear": &year,
	}
}
