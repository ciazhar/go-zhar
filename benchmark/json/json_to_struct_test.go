package json

import (
	"bytes"
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/a8m/djson"
	"github.com/antonholmquist/jason"
	"github.com/bitly/go-simplejson"
	"github.com/buger/jsonparser"
	"github.com/francoispqt/gojay"
	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
	"github.com/mreiferson/go-ujson"
	"github.com/pquerna/ffjson/ffjson"
	segmentio "github.com/segmentio/encoding/json"
	"github.com/ugorji/go/codec"
	"testing"
)

func BenchmarkDJSON(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		object, _ := djson.DecodeObject([]byte(Json100Byte))
		_ = Person{
			ID:      int(object["id"].(float64)),
			Name:    object["name"].(string),
			Age:     int(object["age"].(float64)),
			City:    object["city"].(string),
			Country: object["country"].(string),
		}
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkDJSONAllocString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dec := djson.NewDecoder([]byte(Json100Byte))
		dec.AllocString()
		object, _ := dec.DecodeObject()
		_ = Person{
			ID:      int(object["id"].(float64)),
			Name:    object["name"].(string),
			Age:     int(object["age"].(float64)),
			City:    object["city"].(string),
			Country: object["country"].(string),
		}
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkEasyJson(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		easyjson.Unmarshal([]byte(Json100Byte), &person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkFFjson(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		ffjson.Unmarshal([]byte(Json100Byte), &person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkGabs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseJSON, _ := gabs.ParseJSON([]byte(Json100Byte))
		_ = Person{
			ID:      int(parseJSON.S("id").Data().(float64)),
			Name:    parseJSON.S("name").Data().(string),
			Age:     int(parseJSON.S("age").Data().(float64)),
			Country: parseJSON.S("country").Data().(string),
			City:    parseJSON.S("city").Data().(string),
		}
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkGojayUnmarshal(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		gojay.UnmarshalJSONObject([]byte(Json100Byte), &person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkGojayDecode(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		dec := gojay.NewDecoder(bytes.NewReader([]byte(Json100Byte)))
		dec.DecodeObject(&person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkGoCodec(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		decoder := codec.NewDecoderBytes([]byte(Json100Byte), new(codec.JsonHandle))
		decoder.Decode(&person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkGoJson(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		gojson.Unmarshal([]byte(Json100Byte), &person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkGoSimpleJson(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newJson, _ := simplejson.NewJson([]byte(Json100Byte))
		_ = Person{
			ID:      newJson.Get("id").MustInt(),
			Name:    newJson.Get("name").MustString(),
			Age:     newJson.Get("age").MustInt(),
			City:    newJson.Get("city").MustString(),
			Country: newJson.Get("country").MustString(),
		}
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkGoUJson(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newJson, _ := ujson.NewFromBytes([]byte(Json100Byte))
		_ = Person{
			ID:      int(newJson.Get("id").Int64()),
			Name:    newJson.Get("name").String(),
			Age:     int(newJson.Get("age").Int64()),
			City:    newJson.Get("city").String(),
			Country: newJson.Get("country").String(),
		}
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkJason(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, _ := jason.NewObjectFromBytes([]byte(Json100Byte))

		id, _ := v.GetInt64("id")
		name, _ := v.GetString("name")
		age, _ := v.GetInt64("age")
		city, _ := v.GetString("city")
		country, _ := v.GetString("country")

		_ = Person{
			ID:      int(id),
			Name:    name,
			Age:     int(age),
			City:    city,
			Country: country,
		}
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

// NOT SUPPORTED
//func BenchmarkJettison(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//
//	}
//}

// NOT SUPPORTED
//func BenchmarkJingo(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//
//	}
//}

func BenchmarkJsonParser(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, _ := jsonparser.GetInt([]byte(Json100Byte), "id")
		name, _ := jsonparser.GetString([]byte(Json100Byte), "name")
		age, _ := jsonparser.GetInt([]byte(Json100Byte), "age")
		city, _ := jsonparser.GetString([]byte(Json100Byte), "city")
		country, _ := jsonparser.GetString([]byte(Json100Byte), "country")

		_ = Person{
			ID:      int(id),
			Name:    name,
			Age:     int(age),
			City:    city,
			Country: country,
		}

		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkJsonIterator(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		jsoniter.Unmarshal([]byte(Json100Byte), &person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

func BenchmarkSegmentio(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		segmentio.Unmarshal([]byte(Json100Byte), &person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}

//func BenchmarkSimdJSON(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//
//	}
//}

func BenchmarkStdlib(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person Person
		json.Unmarshal([]byte(Json100Byte), &person)
		//fmt.Println(fmt.Sprintf("%+v", person))
	}
}
