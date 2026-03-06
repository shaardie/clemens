package nnue

import (
	"embed"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"unsafe"
)

const (
	InputSize  = 768
	HiddenSize = 128
	L1Size     = 16
	L2Size     = 16
	EvalScale  = 400.0
)

// model contains the neural network
type model struct {
	// Feature-Transformer: (768, 256) – transponiert gespeichert,
	// damit ein Feature-Vektor zusammenhängend im Speicher liegt.
	FTWeight [InputSize][HiddenSize]float32
	FTBias   [HiddenSize]float32

	L1Weight [L1Size][HiddenSize * 2]float32 // (32, 512)
	L1Bias   [L1Size]float32

	L2Weight [L2Size][L1Size]float32 // (32, 32)
	L2Bias   [L2Size]float32

	OutWeight [L2Size]float32 // (32,)
	OutBias   float32
}

var m model

//go:embed nnue.bin
var fs embed.FS

func init() {
	if err := load(); err != nil {
		panic(err)
	}
}

// load the neural network from the embed file
func load() error {
	data, err := fs.ReadFile("nnue.bin")
	if err != nil {
		return fmt.Errorf("failed to read embed nnue.bin, %w", err)
	}

	off := 0

	read32 := func() uint32 {
		v := binary.LittleEndian.Uint32(data[off:])
		off += 4
		return v
	}
	readF32 := func() float32 {
		return math.Float32frombits(read32())
	}

	// Check Header
	if read32() != InputSize || read32() != HiddenSize ||
		read32() != L1Size || read32() != L2Size {
		return errors.New("model size does not match")
	}

	// FT weights (768 × 256)
	for i := range m.FTWeight {
		for j := range m.FTWeight[i] {
			m.FTWeight[i][j] = readF32()
		}
	}
	for j := range m.FTBias {
		m.FTBias[j] = readF32()
	}

	// L1 (32 × 512)
	for i := range m.L1Weight {
		for j := range m.L1Weight[i] {
			m.L1Weight[i][j] = readF32()
		}
	}
	for i := range m.L1Bias {
		m.L1Bias[i] = readF32()
	}

	// L2 (32 × 32)
	for i := range m.L2Weight {
		for j := range m.L2Weight[i] {
			m.L2Weight[i][j] = readF32()
		}
	}
	for i := range m.L2Bias {
		m.L2Bias[i] = readF32()
	}

	// Output (32 + 1)
	for i := range m.OutWeight {
		m.OutWeight[i] = readF32()
	}
	m.OutBias = readF32()

	expected := 16 + (InputSize*HiddenSize+HiddenSize+
		L1Size*HiddenSize*2+L1Size+
		L2Size*L1Size+L2Size+
		L2Size+1)*int(unsafe.Sizeof(float32(0)))
	if off != expected {
		return fmt.Errorf("unexpected file size (%d vs %d)", off, expected)
	}
	return nil
}
