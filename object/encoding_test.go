package object

import (
	"math"
	"reflect"
	"testing"
)

func TestIntegerEncoding(t *testing.T) {
	tests := []struct {
		obj          *Integer
		expectedRead int
	}{
		{&Integer{Value: 0}, 9},
		{&Integer{Value: 12}, 9},
		{&Integer{Value: -12}, 9},
		{&Integer{Value: math.MaxInt64}, 9},
		{&Integer{Value: -math.MaxInt64}, 9},
	}
	for _, tt := range tests {
		gotBytes, _ := tt.obj.MarshalBytes()

		var obj Integer
		gotRead, _ := obj.UnmarshalBytes(gotBytes)
		if gotRead != tt.expectedRead {
			t.Errorf("invalid bytes read: got %v - want %v", gotRead, tt.expectedRead)
			continue
		}

		if !reflect.DeepEqual(&obj, tt.obj) {
			t.Errorf("invalid unmarshal object: got %+v - want %+v", &obj, tt.obj)
			continue
		}
	}
}

func TestFloatEncoding(t *testing.T) {
	tests := []struct {
		obj          *Float
		expectedRead int
	}{
		{&Float{Value: 0}, 9},
		{&Float{Value: 1.0}, 9},
		{&Float{Value: -1.0}, 9},
		{&Float{Value: math.MaxFloat64}, 9},
		{&Float{Value: -math.MaxFloat64}, 9},
		{&Float{Value: math.MaxFloat32}, 9},
		{&Float{Value: -math.MaxFloat32}, 9},
	}
	for _, tt := range tests {
		gotBytes, _ := tt.obj.MarshalBytes()

		var obj Float
		gotRead, _ := obj.UnmarshalBytes(gotBytes)
		if gotRead != tt.expectedRead {
			t.Errorf("invalid bytes read: got %v - want %v", gotRead, tt.expectedRead)
			continue
		}

		if !reflect.DeepEqual(&obj, tt.obj) {
			t.Errorf("invalid unmarshal object: got %+v - want %+v", &obj, tt.obj)
			continue
		}
	}
}

func TestStringEncoding(t *testing.T) {
	tests := []struct {
		obj          *String
		expectedRead int
	}{
		{&String{Value: ""}, 9},
		{&String{Value: "hello, world"}, 21},
		{&String{Value: "This is a really Long string with punctuation and stuff 1234567890_-()%+{}[]"}, 85},
		{&String{Value: `This 
Is
a
Multiline
String`}, 36}, // Not sure the parser/lexer even supports this, but might as well verify the encoding
	}
	for _, tt := range tests {
		gotBytes, err := tt.obj.MarshalBytes()
		if err != nil {
			t.Error(err)
			continue
		}

		var obj String
		gotRead, err := obj.UnmarshalBytes(gotBytes)
		if err != nil {
			t.Error(err)
			continue
		}
		if gotRead != tt.expectedRead {
			t.Errorf("invalid bytes read: got %v - want %v", gotRead, tt.expectedRead)
			continue
		}

		if !reflect.DeepEqual(&obj, tt.obj) {
			t.Errorf("invalid unmarshal object: got %+v - want %+v", &obj, tt.obj)
			continue
		}
	}
}

func TestCompiledFunctionEncoding(t *testing.T) {
	tests := []struct {
		obj          *CompiledFunction
		expectedRead int
	}{
		{
			obj: &CompiledFunction{
				Instructions: []byte{0, 0, 0, 0, 0, 1, 7, 8, 0, 0, 2, 0, 0, 3, 12, 22, 0, 5, 18, 0, 0},
				NumLocals:    5,
				NumParams:    3,
			},
			expectedRead: 46,
		},
		{
			obj: &CompiledFunction{
				Instructions: []byte{},
				NumLocals:    0,
				NumParams:    0,
			},
			expectedRead: 25,
		},
	}
	for _, tt := range tests {
		gotBytes, err := tt.obj.MarshalBytes()
		if err != nil {
			t.Error(err)
			continue
		}

		var obj CompiledFunction
		gotRead, err := obj.UnmarshalBytes(gotBytes)
		if err != nil {
			t.Error(err)
			continue
		}
		if gotRead != tt.expectedRead {
			t.Errorf("invalid bytes read: got %v - want %v", gotRead, tt.expectedRead)
			continue
		}

		if !reflect.DeepEqual(&obj, tt.obj) {
			t.Errorf("invalid unmarshal object: got %+v - want %+v", &obj, tt.obj)
			continue
		}
	}
}
