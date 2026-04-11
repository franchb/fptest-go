package examples_test

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/franchb/fptest-go/prop"
	"pgregory.net/rapid"
)

// --- Domain types ---

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip"`
}

type UserProfile struct {
	User    User    `json:"user"`
	Address Address `json:"address"`
	Bio     string  `json:"bio"`
}

// --- Generators ---

func genUser() *rapid.Generator[User] {
	return rapid.Custom(func(t *rapid.T) User {
		return User{
			Name:  rapid.String().Draw(t, "name"),
			Age:   rapid.IntRange(0, 150).Draw(t, "age"),
			Email: rapid.String().Draw(t, "email"),
		}
	})
}

func genAddress() *rapid.Generator[Address] {
	return rapid.Custom(func(t *rapid.T) Address {
		return Address{
			Street: rapid.String().Draw(t, "street"),
			City:   rapid.String().Draw(t, "city"),
			Zip:    rapid.StringMatching(`[0-9]{5}`).Draw(t, "zip"),
		}
	})
}

func genUserProfile() *rapid.Generator[UserProfile] {
	return rapid.Custom(func(t *rapid.T) UserProfile {
		return UserProfile{
			User:    genUser().Draw(t, "user"),
			Address: genAddress().Draw(t, "address"),
			Bio:     rapid.String().Draw(t, "bio"),
		}
	})
}

// --- Tests ---

// TestUserJSONRoundTrip verifies that json.Marshal/json.Unmarshal form a valid
// codec pair for User structs. Property-based testing catches edge cases that
// hand-picked examples miss: empty strings, zero values, Unicode, etc.
func TestUserJSONRoundTrip(t *testing.T) {
	prop.RoundTripError(t, "UserJSON",
		genUser(),
		func(a, b User) bool { return a == b },
		func(u User) []byte {
			bs, err := json.Marshal(u)
			if err != nil {
				panic(err) // Marshal should never fail for this struct
			}
			return bs
		},
		func(bs []byte) (User, error) {
			var u User
			err := json.Unmarshal(bs, &u)
			return u, err
		},
	)
}

// TestNestedStructJSONRoundTrip verifies JSON round-trip for a struct containing
// nested structs. Ensures no data is lost at any nesting level.
func TestNestedStructJSONRoundTrip(t *testing.T) {
	prop.RoundTripError(t, "UserProfileJSON",
		genUserProfile(),
		func(a, b UserProfile) bool { return a == b },
		func(p UserProfile) []byte {
			bs, err := json.Marshal(p)
			if err != nil {
				panic(err)
			}
			return bs
		},
		func(bs []byte) (UserProfile, error) {
			var p UserProfile
			err := json.Unmarshal(bs, &p)
			return p, err
		},
	)
}

// TestBase64RoundTrip verifies that base64 encode/decode is a valid codec pair
// for arbitrary byte sequences. Uses the non-error RoundTrip variant.
func TestBase64RoundTrip(t *testing.T) {
	prop.RoundTrip(t, "Base64",
		rapid.SliceOf(rapid.Byte()),
		func(a, b []byte) bool {
			if len(a) == 0 && len(b) == 0 {
				return true
			}
			if len(a) != len(b) {
				return false
			}
			for i := range a {
				if a[i] != b[i] {
					return false
				}
			}
			return true
		},
		base64.StdEncoding.EncodeToString,
		func(s string) []byte {
			bs, _ := base64.StdEncoding.DecodeString(s)
			return bs
		},
	)
}

// TestIntParseRoundTripPartial demonstrates RoundTripPartial with strconv.Itoa
// and a decode function that returns (int, bool) instead of (int, error).
func TestIntParseRoundTripPartial(t *testing.T) {
	prop.RoundTripPartial(t, "IntParse",
		rapid.IntRange(-1_000_000, 1_000_000),
		func(a, b int) bool { return a == b },
		strconv.Itoa,
		func(s string) (int, bool) {
			n, err := strconv.Atoi(s)
			return n, err == nil
		},
	)
}
