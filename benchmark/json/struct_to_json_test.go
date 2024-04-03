package json

import (
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/bet365/jingo"
	"github.com/francoispqt/gojay"
	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
	"github.com/pquerna/ffjson/ffjson"
	segmentio "github.com/segmentio/encoding/json"
	"github.com/ugorji/go/codec"
	"github.com/wI2L/jettison"
	"strings"
	"testing"
)

//NOT SUPPORTED
//func BenchmarkDjsonSTJ(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//
//	}
//}

func BenchmarkEasyJsonSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshal, _ := easyjson.Marshal(Struct100Byte)
		_ = string(marshal)
		//println(str)
	}
}

func BenchmarkFFjsonSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshal, _ := ffjson.Marshal(Struct100Byte)
		_ = string(marshal)
		//println(str)
	}
}

func BenchmarkGabsSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		consume, _ := gabs.Consume(Struct100Byte)
		_ = consume.String()
		//println(str)
	}
}

func BenchmarkGojayMarshalSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshal, _ := gojay.MarshalJSONObject(&Struct100Byte)
		_ = string(marshal)
		//println(str)
	}
}

func BenchmarkGojayEncodeSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := strings.Builder{}
		enc := gojay.NewEncoder(&b)
		enc.Encode(&Struct100Byte)
		_ = b.String()
		//println(str)
	}
}

func BenchmarkGoCodecSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var b []byte
		decoder := codec.NewEncoderBytes(&b, new(codec.JsonHandle))
		decoder.Encode(&Struct100Byte)
		_ = string(b)
		//println(str)
	}
}

func BenchmarkGoJsonSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshal, _ := gojson.Marshal(Struct100Byte)
		_ = string(marshal)
		//println(str)
	}
}

//NOT SUPPORTED
//func BenchmarkGoSimpleJsonSTJ(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//
//	}
//}

//NOT SUPPORTED
//func BenchmarkGoUJsonSTJ(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//
//	}
//}

//NOT SUPPORTED
//func BenchmarkJasonSTJ(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {

//	}
//}

func BenchmarkJettisonSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b, _ := jettison.Marshal(Struct100Byte)
		_ = string(b)
		//println(str)
	}
}

func BenchmarkJingoSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var enc = jingo.NewStructEncoder(Person{})
		buf := jingo.NewBufferFromPool()
		enc.Marshal(&Struct100Byte, buf)
		_ = buf.String()
		//println(str)
	}
}

//NOT SUPPORTED
//func BenchmarkJsonParserSTJ(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//
//	}
//}

func BenchmarkJsonIteratorSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b, _ := jsoniter.Marshal(Struct100Byte)
		_ = string(b)
		//println(str)
	}
}

func BenchmarkSegmentioSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshal, _ := segmentio.Marshal(Struct100Byte)
		_ = string(marshal)
		//println(str)
	}
}

//func BenchmarkSimdJSONSTJ(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//
//	}
//}

func BenchmarkStdlibSTJ(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshal, _ := json.Marshal(Struct100Byte)
		_ = string(marshal)
		//println(str)
	}
}
