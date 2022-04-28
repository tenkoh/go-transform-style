package internal

import (
	"reflect"
	"testing"
)

func TestInRange(t *testing.T) {
	tests := []struct {
		src  uint8
		min  uint8
		max  uint8
		want bool
	}{
		{15, 10, 20, true},
		{10, 10, 20, true},
		{20, 10, 20, true},
		{1, 10, 20, false},
	}
	for _, tt := range tests {
		if got := inRange(tt.src, tt.min, tt.max); got != tt.want {
			t.Errorf("want %v, got %v\n", tt.want, got)
		}
	}
}

func TestIsRegularLower(t *testing.T) {
	tests := []struct {
		src  []byte
		want bool
	}{
		{[]byte("a"), true},
		{[]byte("z"), true},
		{[]byte("A"), false},
		{[]byte("Z"), false},
		{[]byte("0"), false},
		{[]byte("9"), false},
		{[]byte("あ"), false},
		{[]byte(""), false},
	}

	for _, tt := range tests {
		if got := isRegularLower(tt.src); got != tt.want {
			t.Errorf("input: %s, got %v, want %v\n", string(tt.src), got, tt.want)
		}
	}
}

func TestIsRegularUpper(t *testing.T) {
	tests := []struct {
		src  []byte
		want bool
	}{
		{[]byte("A"), true},
		{[]byte("Z"), true},
		{[]byte("a"), false},
		{[]byte("z"), false},
		{[]byte("0"), false},
		{[]byte("9"), false},
		{[]byte("あ"), false},
		{[]byte(""), false},
	}

	for _, tt := range tests {
		if got := isRegularUpper(tt.src); got != tt.want {
			t.Errorf("input: %s, got %v, want %v\n", string(tt.src), got, tt.want)
		}
	}
}

func TestIsRegularDigit(t *testing.T) {
	tests := []struct {
		src  []byte
		want bool
	}{
		{[]byte("0"), true},
		{[]byte("9"), true},
		{[]byte("a"), false},
		{[]byte("z"), false},
		{[]byte("A"), false},
		{[]byte("Z"), false},
		{[]byte("あ"), false},
		{[]byte(""), false},
	}

	for _, tt := range tests {
		if got := isRegularDigit(tt.src); got != tt.want {
			t.Errorf("input: %s, got %v, want %v\n", string(tt.src), got, tt.want)
		}
	}
}

func TestReplacer_Replace(t *testing.T) {
	type args struct {
		r   *replacer
		src []byte
	}
	r1 := &replacer{
		lowerFunc: func(b uint8) []byte { return []byte{b + 1} },
		upperFunc: func(b uint8) []byte { return []byte{b + 2} },
		digitFunc: func(b uint8) []byte { return []byte{b + 3} },
	}
	tests := []struct {
		arg  args
		want []byte
	}{
		{args{r1, []byte("aA1あ")}, []byte("bC4あ")},
		{args{r1, []byte{250, 100, 70, 50, 250}}, []byte{250, 101, 72, 53, 250}},
	}
	for _, tt := range tests {
		if got := tt.arg.r.replace(tt.arg.src); !reflect.DeepEqual(tt.want, got) {
			t.Errorf("want %v, got %v\n", tt.want, got)
		}
	}
}

func TestTransformer_Transform(t *testing.T) {
	type args struct {
		src   []byte
		dst   []byte
		atEOF bool
	}
	type wants struct {
		wrote            []byte
		nSrc             int
		nDst             int
		stockToTransform []byte
		stockToWrite     []byte
		err              error
	}
	r1 := &replacer{
		lowerFunc: func(b uint8) []byte { return []byte{b + 1} },
		upperFunc: func(b uint8) []byte { return []byte{b + 2} },
		digitFunc: func(b uint8) []byte { return []byte{b + 3} },
	}
	tr1 := &Transformer{
		rep:              r1,
		stockToWrite:     nil,
		stockToTransform: nil,
	}
	tests := []struct {
		tr    *Transformer
		args  []args
		wants []wants
	}{
		{
			tr1,
			[]args{{[]byte("aA1あ"), make([]byte, 10), false}},
			[]wants{{fillBytes([]byte("bC4あ"), 10), 6, 6, nil, nil, nil}},
		},
	}

	for _, tt := range tests {
		if len(tt.args) != len(tt.wants) {
			t.Fatal("invalid test condition: len(args) must be same as len(wants)")
		}
		for i, arg := range tt.args {
			nDst, nSrc, err := tt.tr.Transform(arg.dst, arg.src, arg.atEOF)
			want := tt.wants[i]
			if !reflect.DeepEqual(arg.dst, want.wrote) {
				t.Errorf("wrote bytes mismatched; want %v, got %v\n", want.wrote, arg.dst)
			}
			if nDst != want.nDst {
				t.Errorf("nDst mismatched; want %d, got %d\n", want.nDst, nDst)
			}
			if nSrc != want.nSrc {
				t.Errorf("nSrc mismatched; want %d, got %d\n", want.nSrc, nSrc)
			}
			if err != want.err {
				t.Errorf("error mismatched; want %s, got %s\n", want.err, err)
			}
			if got := tt.tr.stockToTransform; !reflect.DeepEqual(got, want.stockToTransform) {
				t.Errorf("stockToTransform mismatched; want %v, got %v\n", got, want.stockToTransform)
			}
			if got := tt.tr.stockToWrite; !reflect.DeepEqual(got, want.stockToWrite) {
				t.Errorf("stockToWrite mismatched; want %v, got %v\n", got, want.stockToWrite)
			}
		}
		tt.tr.Reset()
	}
}

func fillBytes(src []byte, nDst int) []byte {
	var n int
	dst := make([]byte, nDst)
	if nDst < len(src) {
		n = nDst
	} else {
		n = len(src)
	}
	for i := 0; i < n; i++ {
		dst[i] = src[i]
	}
	return dst
}
