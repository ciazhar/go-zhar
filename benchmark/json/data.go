package json

import "github.com/francoispqt/gojay"

type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	City    string `json:"city"`
}

func (u *Person) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "id":
		return dec.Int(&u.ID)
	case "name":
		return dec.String(&u.Name)
	case "age":
		return dec.Int(&u.Age)
	case "country":
		return dec.String(&u.Country)
	case "city":
		return dec.String(&u.City)
	}
	return nil
}
func (u *Person) NKeys() int {
	return 5
}

func (u *Person) MarshalJSONObject(enc *gojay.Encoder) {
	enc.IntKey("id", u.ID)
	enc.StringKey("name", u.Name)
	enc.IntKey("age", u.Age)
	enc.StringKey("country", u.Country)
	enc.StringKey("city", u.City)
}
func (u *Person) IsNil() bool {
	return u == nil
}

var Struct100Byte = Person{
	ID:      123,
	Name:    "John Doe",
	Age:     30,
	Country: "USA",
	City:    "New York",
}

var Json100Byte = `{"id":123,"name":"John Doe","age":30,"country":"USA","city":"New York"}`
