package bitboard

import (
	"testing"
)

func TestEqual(t *testing.T) {
	type args struct {
		b1 Bitboard
		b2 Bitboard
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Equal",
			args: args{
				0b1001,
				0b1001,
			},
			want: true,
		},
		{
			name: "Not Equal",
			args: args{
				0b1001,
				0b1000,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.b1, tt.args.b2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	type args struct {
		b1 Bitboard
		b2 Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "Intersect",
			args: args{
				0b1001,
				0b1001,
			},
			want: 0b1001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersection(tt.args.b1, tt.args.b2); got != tt.want {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	type args struct {
		b1 Bitboard
		b2 Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "Union",
			args: args{
				0b1010,
				0b1001,
			},
			want: 0b1011,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Union(tt.args.b1, tt.args.b2); got != tt.want {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComplement(t *testing.T) {
	type args struct {
		b Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "Emtpy",
			args: args{
				Empty,
			},
			want: Universal,
		},
		{
			name: "Universal",
			args: args{
				Universal,
			},
			want: Empty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Complement(tt.args.b); got != tt.want {
				t.Errorf("Complement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	type args struct {
		b1 Bitboard
		b2 Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "Difference",
			args: args{
				0b1100,
				0b1000,
			},
			want: 0b0100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Difference(tt.args.b1, tt.args.b2); got != tt.want {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImplication(t *testing.T) {
	type args struct {
		b1 Bitboard
		b2 Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "Implication 0 0",
			args: args{
				Empty,
				Empty,
			},
			want: Universal,
		},
		{
			name: "Implication 0 1",
			args: args{
				Empty,
				Universal,
			},
			want: Universal,
		},
		{
			name: "Implication 1 0",
			args: args{
				Universal,
				Empty,
			},
			want: Empty,
		},
		{
			name: "Implication 1 1",
			args: args{
				Universal,
				Universal,
			},
			want: Universal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Implication(tt.args.b1, tt.args.b2); got != tt.want {
				t.Errorf("Implication() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSymmetricDifference(t *testing.T) {
	type args struct {
		b1 Bitboard
		b2 Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "Symmetric Difference",
			args: args{
				0b1101,
				0b0110,
			},
			want: 0b1011,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SymmetricDifference(tt.args.b1, tt.args.b2); got != tt.want {
				t.Errorf("SymmetricDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquivalence(t *testing.T) {
	type args struct {
		b1 Bitboard
		b2 Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "Universal",
			args: args{
				0b1101,
				Universal,
			},
			want: 0b1101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equivalence(tt.args.b1, tt.args.b2); got != tt.want {
				t.Errorf("Equivalence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMajority(t *testing.T) {
	type args struct {
		b1 Bitboard
		b2 Bitboard
		b3 Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "Equivalence",
			args: args{
				0b1101,
				0b0110,
				0b0011,
			},
			want: 0b0111,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Majority(tt.args.b1, tt.args.b2, tt.args.b3); got != tt.want {
				t.Errorf("Majority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGreaterOne(t *testing.T) {
	type args struct {
		bs []Bitboard
	}
	tests := []struct {
		name string
		args args
		want Bitboard
	}{
		{
			name: "nil",
			args: args{nil},
			want: Empty,
		},
		{
			name: "1 Bitboard",
			args: args{[]Bitboard{Universal}},
			want: Empty,
		},
		{
			name: "2 Bitboard",
			args: args{[]Bitboard{0b1001, 0b1100}},
			want: Intersection(0b1001, 0b1100),
		},
		{
			name: "3 Bitboard",
			args: args{[]Bitboard{0b1001, 0b1100, 0b0110}},
			want: Majority(0b1001, 0b1100, 0b0110),
		},
		{
			name: "4 Bitboard",
			args: args{
				[]Bitboard{
					0b100100,
					0b110010,
					0b011001,
					0b111001,
				},
			},
			want: 0b111001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GreaterOne(tt.args.bs...); got != tt.want {
				t.Errorf("GreaterOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
