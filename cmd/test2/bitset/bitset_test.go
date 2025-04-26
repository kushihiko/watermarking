package bitset

import "testing"

func TestBitSet(t *testing.T) {
	type action struct {
		method string
		index  int
		value  bool // for Set/Clear (only relevant for Set if true)
	}

	tests := []struct {
		name        string
		size        int
		actions     []action
		expected    map[int]bool
		expectedLen int
	}{
		{
			name: "Simple Set and Get",
			size: 10,
			actions: []action{
				{method: "set", index: 2},
			},
			expected:    map[int]bool{2: true, 1: false},
			expectedLen: 10,
		},
		{
			name: "Set and Clear",
			size: 5,
			actions: []action{
				{method: "set", index: 3},
				{method: "clear", index: 3},
			},
			expected:    map[int]bool{3: false},
			expectedLen: 5,
		},
		{
			name: "Auto Expand Set",
			size: 5,
			actions: []action{
				{method: "set", index: 20},
			},
			expected:    map[int]bool{20: true},
			expectedLen: 21,
		},
		{
			name:        "Get Out Of Bounds",
			size:        5,
			actions:     []action{},
			expected:    map[int]bool{100: false},
			expectedLen: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitSet(tt.size)
			for _, act := range tt.actions {
				switch act.method {
				case "set":
					bs.Set(act.index)
				case "clear":
					bs.Clear(act.index)
				}
			}
			for index, want := range tt.expected {
				got := bs.Get(index)
				if got != want {
					t.Errorf("at index %d: expected %v, got %v", index, want, got)
				}
			}
			if bs.Len() != tt.expectedLen {
				t.Errorf("expected length %d, got %d", tt.expectedLen, bs.Len())
			}
		})
	}
}
