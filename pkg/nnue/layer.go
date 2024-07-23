package nnue

type linearLayer struct {
	Weights [][]float64 `json:"weight"`
	Bias    []float64   `json:"bias"`
}

func linear(layer *linearLayer, input, output []float64) {
	copy(output, layer.Bias)
	for i := range input {
		for j := range output {
			output[j] += input[i] * layer.Weights[i][j]
		}
	}
}

func crelu(input, output []float64) {
	for i := range output {
		output[i] = min(max(input[i], 0), 1)
	}
}
