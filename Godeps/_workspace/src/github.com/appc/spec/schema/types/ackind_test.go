package types

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestACKindMarshalBad(t *testing.T) {
	tests := map[string]error{
		"Foo": ACKindError("bad ACKind: Foo"),
		"ApplicationManifest": ACKindError("bad ACKind: ApplicationManifest"),
		"": ErrNoACKind,
	}
	for in, werr := range tests {
		a := ACKind(in)
		b, gerr := json.Marshal(a)
		if b != nil {
			t.Errorf("ACKind(%q): want b=nil, got %v", in, b)
		}
		if jerr, ok := gerr.(*json.MarshalerError); !ok {
			t.Errorf("expected JSONMarshalerError")
		} else {
			if e := jerr.Err; e != werr {
				t.Errorf("err=%#v, want %#v", e, werr)
			}
		}
	}
}

func TestACKindMarshalGood(t *testing.T) {
	for i, in := range []string{
		"ImageManifest",
		"PodManifest",
	} {
		a := ACKind(in)
		b, err := json.Marshal(a)
		if !reflect.DeepEqual(b, []byte(`"`+in+`"`)) {
			t.Errorf("#%d: marshalled=%v, want %v", i, b, []byte(in))
		}
		if err != nil {
			t.Errorf("#%d: err=%v, want nil", i, err)
		}
	}
}

func TestACKindUnmarshalBad(t *testing.T) {
	tests := []string{
		"ImageManifest", // Not a valid JSON-encoded string
		`"garbage"`,
		`"AppManifest"`,
		`""`,
	}
	for i, in := range tests {
		var a, b ACKind
		err := a.UnmarshalJSON([]byte(in))
		if err == nil {
			t.Errorf("#%d: err=nil, want non-nil", i)
		} else if !reflect.DeepEqual(a, b) {
			t.Errorf("#%d: a=%v, want empty", i, a)
		}
	}
}

func TestACKindUnmarshalGood(t *testing.T) {
	tests := map[string]ACKind{
		`"PodManifest"`:   ACKind("PodManifest"),
		`"ImageManifest"`: ACKind("ImageManifest"),
	}
	for in, w := range tests {
		var a ACKind
		err := json.Unmarshal([]byte(in), &a)
		if err != nil {
			t.Errorf("%v: err=%v, want nil", in, err)
		} else if !reflect.DeepEqual(a, w) {
			t.Errorf("%v: a=%v, want %v", in, a, w)
		}
	}
}
