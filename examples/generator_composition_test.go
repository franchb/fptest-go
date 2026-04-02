package examples_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/option"
	"github.com/franchb/fptest/gen"
	"github.com/franchb/fptest/prop"
	"pgregory.net/rapid"
)

// --- Domain types ---

type GenUser struct {
	Name  string
	Age   int
	Email string
}

type GenAddress struct {
	Street  string
	City    string
	ZipCode string
}

type LineItem struct {
	Product  string
	Quantity int
	Price    int // in cents
}

type Order struct {
	ID     int
	UserID int
	Items  []LineItem
	Total  int // in cents, must equal sum of qty * price
}

// --- Generators ---

// genNameG generates a human name using the Gen monad.
func genNameG() gen.Gen[string] {
	return gen.FromRapid(rapid.StringMatching(`[A-Z][a-z]{2,10}`))
}

// genAgeG generates an age in [0, 150].
func genAgeG() gen.Gen[int] {
	return gen.FromRapid(rapid.IntRange(0, 150))
}

// genEmailG generates an email address.
func genEmailG() gen.Gen[string] {
	return gen.Map(genNameG(), func(name string) string {
		return name + "@example.com"
	})
}

// genLineItem generates a single line item with positive quantity and price.
func genLineItem() gen.Gen[LineItem] {
	return gen.Map3(
		gen.FromRapid(rapid.StringMatching(`[A-Z][a-z]{2,8}`)),
		gen.FromRapid(rapid.IntRange(1, 20)),
		gen.FromRapid(rapid.IntRange(100, 10000)),
		func(product string, qty, price int) LineItem {
			return LineItem{Product: product, Quantity: qty, Price: price}
		},
	)
}

// --- Tests ---

// TestUserGenWithMap3 demonstrates composing three independent field generators
// into a struct generator using gen.Map3. Each field is generated independently,
// then combined — the Applicative pattern for test data.
func TestUserGenWithMap3(t *testing.T) {
	userGen := gen.Map3(
		genNameG(),
		genAgeG(),
		genEmailG(),
		func(name string, age int, email string) GenUser {
			return GenUser{Name: name, Age: age, Email: email}
		},
	)

	rapid.Check(t, func(t *rapid.T) {
		u := userGen(t)
		if u.Name == "" {
			t.Fatal("empty name")
		}
		if u.Age < 0 || u.Age > 150 {
			t.Fatalf("age out of range: %d", u.Age)
		}
	})
}

// TestOrderWithDependentTotal demonstrates gen.Chain for dependent generation:
// first generate a list of items, then compute the total from those items. This
// ensures the Total field is always consistent with Items — a property that
// independent generation cannot guarantee.
func TestOrderWithDependentTotal(t *testing.T) {
	orderGen := gen.Chain(
		gen.Slice(genLineItem(), 1, 5),
		func(items []LineItem) gen.Gen[Order] {
			total := 0
			for _, item := range items {
				total += item.Quantity * item.Price
			}
			return gen.Map(
				gen.FromRapid(rapid.IntRange(1, 10000)),
				func(id int) Order {
					return Order{ID: id, UserID: 1, Items: items, Total: total}
				},
			)
		},
	)

	rapid.Check(t, func(t *rapid.T) {
		order := orderGen(t)

		// Verify the structural invariant: Total == sum of item subtotals
		computed := 0
		for _, item := range order.Items {
			computed += item.Quantity * item.Price
		}
		if order.Total != computed {
			t.Fatalf("order total %d != computed %d", order.Total, computed)
		}
		if len(order.Items) == 0 {
			t.Fatal("order has no items")
		}
	})
}

// TestOptionalAddress demonstrates gen.MonadicOption to generate Option[Address]
// values — sometimes Some(address), sometimes None. This is useful for testing
// code that handles optional fields gracefully.
func TestOptionalAddress(t *testing.T) {
	addressGen := gen.Map3(
		gen.FromRapid(rapid.String()),
		gen.FromRapid(rapid.String()),
		gen.FromRapid(rapid.StringMatching(`[0-9]{5}`)),
		func(street, city, zip string) GenAddress {
			return GenAddress{Street: street, City: city, ZipCode: zip}
		},
	)

	optAddressGen := gen.MonadicOption(addressGen)

	someCount := 0
	noneCount := 0
	rapid.Check(t, func(t *rapid.T) {
		opt := optAddressGen(t)
		if option.IsSome(opt) {
			someCount++
		} else {
			noneCount++
		}
	})

	// With random bool, we expect both Some and None to appear
	if someCount == 0 {
		t.Fatal("no Some values generated")
	}
	if noneCount == 0 {
		t.Fatal("no None values generated")
	}
}

// TestFilterAdults demonstrates gen.Filter to produce only users with age >= 18.
// Combined with prop.Invariant, this verifies the generator actually respects
// the constraint across all generated values.
func TestFilterAdults(t *testing.T) {
	adultGen := gen.Filter(
		gen.Map3(
			genNameG(),
			genAgeG(),
			genEmailG(),
			func(name string, age int, email string) GenUser {
				return GenUser{Name: name, Age: age, Email: email}
			},
		),
		func(u GenUser) bool { return u.Age >= 18 },
	)

	prop.Invariant(t, "AdultAge",
		gen.ToRapid(adultGen),
		func(u GenUser) bool { return u.Age >= 18 },
	)
}

// TestUserSlice demonstrates gen.Slice to generate lists of users of varying
// length. This is useful for testing collection-processing functions like
// filters, reducers, and batch operations.
func TestUserSlice(t *testing.T) {
	userGen := gen.Map3(
		genNameG(),
		genAgeG(),
		genEmailG(),
		func(name string, age int, email string) GenUser {
			return GenUser{Name: name, Age: age, Email: email}
		},
	)

	usersGen := gen.Slice(userGen, 1, 10)

	rapid.Check(t, func(t *rapid.T) {
		users := usersGen(t)
		if len(users) < 1 || len(users) > 10 {
			t.Fatalf("slice length %d outside [1, 10]", len(users))
		}
	})
}

// TestGenToRapidBridge demonstrates gen.ToRapid, which converts a monadic
// Gen[A] to a *rapid.Generator[A]. This is the bridge between monadic generator
// composition and the prop/laws packages that accept rapid generators.
func TestGenToRapidBridge(t *testing.T) {
	// Build a monadic generator
	userGen := gen.Map3(
		genNameG(),
		genAgeG(),
		genEmailG(),
		func(name string, age int, email string) GenUser {
			return GenUser{Name: name, Age: age, Email: email}
		},
	)

	// Convert to rapid.Generator for use with prop.Invariant
	rapidGen := gen.ToRapid(userGen)

	prop.Invariant(t, "UserStructValid",
		rapidGen,
		func(u GenUser) bool {
			return u.Name != "" && u.Age >= 0 && u.Age <= 150
		},
	)
}
