package chip8

import (
	"testing"
)

type treg struct {
	reg   uint8
	value uint8
}

func setup() *CPU {
	m := Memory{}
	g := NewGraphics(&m)
	return NewCPU(&m, g)
}

// TODO: Split unit tests into seperate files.

func TestPush(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		expected  uint16
		expectErr bool
	}{
		{name: "Valid data FFF", inst: 0x1FFF, expected: 0x1FFF, expectErr: false},
		{name: "Valid data AEB", inst: 0x1AEB, expected: 0x1AEB, expectErr: false},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			err := c8.Push(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute push: %v", err)
			}
			if c8.Stack[c8.SP] != tc.expected {
				t.Errorf("expected: %x; got: %x", tc.expected, c8.Stack[c8.SP])
			}
		})
	}
	t.Run("Exceed stack limit.", func(t *testing.T) {
		c8 := setup()
		i := 0
		//Fill stack with data.
		for i < len(c8.Stack) {
			err := c8.Push(0xFFFF)
			if err != nil {
				t.Fatalf("failed execute push: %v", err)
			}
			i++
		}
		//Try and push data when stack is full.
		err := c8.Push(0xFFFF)
		if err == nil {
			t.Errorf("expected stack limit error; sp: %d; stack: %v", c8.SP, c8.Stack)
		}

	})
}

func TestJump(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		expected  uint16
		expectErr bool
	}{
		{name: "Valid Address FFF", inst: 0x1FFF, expected: 0x0FFF, expectErr: false},
		{name: "Valid Address AEB", inst: 0x1AEB, expected: 0x0AEB, expectErr: false},
		{name: "Invalid Jump Instruction", inst: 0x2AEB, expected: 0x0AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			err := c8.Jump(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute jump: %v", err)
			}
			if c8.PC != tc.expected {
				t.Errorf("expected: %x; got: %x", tc.expected, c8.PC)
			}
		})
	}
}

func TestCall(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		pc        uint16
		s         uint16
		expectErr bool
	}{
		//Stack doesn't reset after each test.
		{name: "Valid Address FFF", inst: 0x2FFF, pc: 0x0FFF, s: 0x200, expectErr: false},
		{name: "Valid Address AEB", inst: 0x2AEB, pc: 0x0AEB, s: 0x200, expectErr: false},
		{name: "Invalid Call Instruction", inst: 0x1AEB, pc: 0x0AEB, s: 0xFFF, expectErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			err := c8.Call(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute call: %v", err)
			}
			if c8.PC != tc.pc {
				t.Errorf("PC - expected: %x; got: %x", tc.pc, c8.PC)
			}
			if c8.Stack[c8.SP] != tc.s {
				t.Errorf("SP - expected: %x; got: %x, stack: %v", tc.s, c8.Stack[c8.SP], c8.Stack)
			}
		})
	}
	t.Run("Call Stack limit test.", func(t *testing.T) {
		c8 := setup()
		i := 0
		//Fill stack with data.
		for i < len(c8.Stack) {
			err := c8.Call(0x2FFF)
			if err != nil {
				t.Fatalf("failed execute Call: %v", err)
			}
			i++
		}
		//Try and push data when stack is full.
		err := c8.Push(0x2FFF)
		if err == nil {
			t.Errorf("expected stack limit error; sp: %d; stack: %v", c8.SP, c8.Stack)
		}

	})
}

func TestSkipEqualVal(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       treg
		expected  uint16
		expectErr bool
	}{
		{name: "Equal reg/value.", inst: 0x330A, reg: treg{reg: 0x03, value: 10}, expected: PCInit + 1, expectErr: false},
		{name: "Unequal reg/value", inst: 0x33FF, reg: treg{reg: 0x03, value: 10}, expected: PCInit, expectErr: false},
		{name: "Invalid Instruction", inst: 0x23FF, reg: treg{reg: 0x03, value: 10}, expected: PCInit, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			c8.V[tc.reg.reg] = tc.reg.value
			err := c8.SkipEqualVal(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute SkipEqualVal: %v", err)
			}
			if c8.PC != tc.expected {
				t.Errorf("expected: %x; got: %x", tc.expected, c8.PC)
			}
		})
	}
}

func TestSkipNotEqualVal(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       treg
		expected  uint16
		expectErr bool
	}{
		{name: "Equal reg/value.", inst: 0x4A0A, reg: treg{reg: 0x0A, value: 10}, expected: PCInit, expectErr: false},
		{name: "Unequal reg/value", inst: 0x4AFF, reg: treg{reg: 0x0A, value: 10}, expected: PCInit + 1, expectErr: false},
		{name: "Invalid Instruction", inst: 0x2AFF, reg: treg{reg: 0x0A, value: 10}, expected: PCInit, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			c8.V[tc.reg.reg] = tc.reg.value
			err := c8.SkipNotEqualVal(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute SkipEqualVal: %v", err)
			}
			if c8.PC != tc.expected {
				t.Errorf("expected: %x; got: %x", tc.expected, c8.PC)
			}
		})
	}
}

func TestSkipEqualReg(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  uint16
		expectErr bool
	}{
		{name: "Equal reg/reg.", inst: 0x51E0,
			reg: []treg{
				treg{reg: 0x01, value: 10},
				treg{reg: 0x0E, value: 10}},
			expected: PCInit + 1, expectErr: false},
		{name: "Unequal reg/reg", inst: 0x52E0,
			reg: []treg{
				treg{reg: 0x02, value: 10},
				treg{reg: 0x0E, value: 11}},
			expected: PCInit, expectErr: false},
		{name: "Invalid Instruction", inst: 0x2FFF,
			reg: []treg{
				treg{reg: 0x01, value: 10},
				treg{reg: 0x0E, value: 10}},
			expected: PCInit, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.SkipEqualReg(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute SkipEqualReg: %v", err)
			}
			if c8.PC != tc.expected {
				t.Errorf("expected: %x; got: %x", tc.expected, c8.PC)
			}
		})
	}
}

func TestLoadValue(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		expected  treg
		expectErr bool
	}{
		{name: "Load FF into VF", inst: 0x6BFF, expected: treg{reg: 0x0B, value: 0xFF}, expectErr: false},
		{name: "Load EB into VA", inst: 0x6AEB, expected: treg{reg: 0x0A, value: 0xEB}, expectErr: false},
		{name: "Invalid Instruction", inst: 0x2AEB, expected: treg{reg: 0x0F, value: 0xFF}, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			err := c8.LoadValue(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute LoadValue: %v", err)
			}
			if c8.V[tc.expected.reg] != tc.expected.value {
				t.Errorf("expected: %x; got: %x", tc.expected.value, c8.V[tc.expected.reg])
			}
		})
	}
}

func TestAddValue(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  treg
		expectErr bool
	}{
		{name: "Load FF into VB", inst: 0x7B0A,
			reg: []treg{
				treg{reg: 0x0B, value: 10}},
			expected: treg{reg: 0x0B, value: 0x14}, expectErr: false},
		{name: "Load EB into VA", inst: 0x7EEB,
			reg: []treg{
				treg{reg: 0x0E, value: 15}},
			expected: treg{reg: 0x0E, value: 0xFA}, expectErr: false},
		{name: "Invalid Instruction", inst: 0x2AEB,
			reg: []treg{
				treg{reg: 0x03, value: 10},
				treg{reg: 0x0E, value: 10}},
			expected: treg{reg: 0x03, value: 0xFF}, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.AddValue(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute AddValue: %v", err)
			}
			if c8.V[tc.expected.reg] != tc.expected.value {
				t.Errorf("expected: %x; got: %x", tc.expected.value, c8.V[tc.expected.reg])
			}
		})
	}
}

func TestLoadReg(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  treg
		expectErr bool
	}{
		{name: "Load VB into V7", inst: 0x87B0,
			reg: []treg{
				treg{reg: 0x07, value: 5},
				treg{reg: 0x0B, value: 10}},
			expected: treg{reg: 0x07, value: 0x0A}, expectErr: false},
		{name: "Load V2 into VA", inst: 0x8A20,
			reg: []treg{
				treg{reg: 0x0A, value: 5},
				treg{reg: 0x02, value: 15}},
			expected: treg{reg: 0x0A, value: 0x0F}, expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid LoadReg Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.LoadReg(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute LoadReg: %v", err)
			}
			if c8.V[tc.expected.reg] != tc.expected.value {
				t.Errorf("expected: %x; got: %x", tc.expected.value, c8.V[tc.expected.reg])
			}
		})
	}
}

func TestOr(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  treg
		expectErr bool
	}{
		{name: "Or VB | V7", inst: 0x87B1,
			reg: []treg{
				treg{reg: 0x07, value: 70},
				treg{reg: 0x0B, value: 25}},
			expected: treg{reg: 0x07, value: 95}, expectErr: false},
		{name: "Or V2 | VA", inst: 0x8A21,
			reg: []treg{
				treg{reg: 0x0A, value: 150},
				treg{reg: 0x02, value: 95}},
			expected: treg{reg: 0x0A, value: 223}, expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid Or Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.Or(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Or: %v", err)
			}
			if c8.V[tc.expected.reg] != tc.expected.value {
				t.Errorf("expected: %x; got: %x", tc.expected.value, c8.V[tc.expected.reg])
			}
		})
	}
}

func TestAnd(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  treg
		expectErr bool
	}{
		{name: "Or VB | V7", inst: 0x87B2,
			reg: []treg{
				treg{reg: 0x07, value: 70},
				treg{reg: 0x0B, value: 25}},
			expected: treg{reg: 0x07, value: 0}, expectErr: false},
		{name: "Or V2 | VA", inst: 0x8A22,
			reg: []treg{
				treg{reg: 0x0A, value: 150},
				treg{reg: 0x02, value: 95}},
			expected: treg{reg: 0x0A, value: 22}, expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid And Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.And(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Or: %v", err)
			}
			if c8.V[tc.expected.reg] != tc.expected.value {
				t.Errorf("expected: %x; got: %x", tc.expected.value, c8.V[tc.expected.reg])
			}
		})
	}
}

func TestXor(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  treg
		expectErr bool
	}{
		{name: "Xor VB ^ V7", inst: 0x87B3,
			reg: []treg{
				treg{reg: 0x07, value: 187},
				treg{reg: 0x0B, value: 25}},
			expected: treg{reg: 0x07, value: 162}, expectErr: false},
		{name: "Xor V2 ^ VA", inst: 0x8A23,
			reg: []treg{
				treg{reg: 0x0A, value: 150},
				treg{reg: 0x02, value: 95}},
			expected: treg{reg: 0x0A, value: 201}, expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid Xor Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.Xor(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Xor: %v", err)
			}
			if c8.V[tc.expected.reg] != tc.expected.value {
				t.Errorf("expected: %x; got: %x", tc.expected.value, c8.V[tc.expected.reg])
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  []treg
		expectErr bool
	}{
		{name: "Normal Add", inst: 0x87B4,
			reg: []treg{
				treg{reg: 0x07, value: 187},
				treg{reg: 0x0B, value: 25}},
			expected: []treg{
				{reg: 0x07, value: 212},
				{reg: 0x0F, value: 0}},
			expectErr: false},
		{name: "Add with carry", inst: 0x8A24,
			reg: []treg{
				treg{reg: 0x0A, value: 255},
				treg{reg: 0x02, value: 255}},
			expected: []treg{
				{reg: 0x0A, value: 0xFE},
				{reg: 0x0F, value: 1}},
			expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid Add Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.Add(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Add: %v", err)
			}
			for _, r := range tc.expected {
				if c8.V[r.reg] != r.value {
					t.Errorf("reg: V%x; expected: %x; got: %x", r.reg, r.value, c8.V[r.reg])
				}
			}
		})
	}
}

func TestSub(t *testing.T) {

	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  []treg
		expectErr bool
	}{
		{name: "Sub X > Y", inst: 0x87B5,
			reg: []treg{
				treg{reg: 0x07, value: 187},
				treg{reg: 0x0B, value: 25}},
			expected: []treg{
				{reg: 0x07, value: 162},
				{reg: 0x0F, value: 1}},
			expectErr: false},
		{name: "Sub X == Y", inst: 0x8A25,
			reg: []treg{
				treg{reg: 0x0A, value: 255},
				treg{reg: 0x02, value: 255}},
			expected: []treg{
				{reg: 0x0A, value: 0},
				{reg: 0x0F, value: 0}},
			expectErr: false},
		{name: "Sub X < Y", inst: 0x8A25,
			reg: []treg{
				treg{reg: 0x0A, value: 13},
				treg{reg: 0x02, value: 160}},
			expected: []treg{
				{reg: 0x0A, value: 109},
				{reg: 0x0F, value: 0}},
			expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid Sub Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.Sub(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Sub: %v", err)
			}
			for _, r := range tc.expected {
				if c8.V[r.reg] != r.value {
					t.Errorf("reg: V%x; expected: %x; got: %x", r.reg, r.value, c8.V[r.reg])
				}
			}
		})
	}
}

func TestShiftRight(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  []treg
		expectErr bool
	}{
		{name: "Shift 209 w/Carry", inst: 0x87B6,
			reg: []treg{
				treg{reg: 0x07, value: 209}},
			expected: []treg{
				{reg: 0x07, value: 104},
				{reg: 0x0F, value: 1}},
			expectErr: false},
		{name: "Shift 250 w/o Carry", inst: 0x8A26,
			reg: []treg{
				treg{reg: 0x0A, value: 250}},
			expected: []treg{
				{reg: 0x0A, value: 125},
				{reg: 0x0F, value: 0}},
			expectErr: false},
		{name: "Shift 13 w/Carry", inst: 0x8A26,
			reg: []treg{
				treg{reg: 0x0A, value: 13}},
			expected: []treg{
				{reg: 0x0A, value: 6},
				{reg: 0x0F, value: 1}},
			expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid ShiftRight Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.ShiftRight(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Sub: %v", err)
			}
			for _, r := range tc.expected {
				if c8.V[r.reg] != r.value {
					t.Errorf("reg: V%x; expected: %x; got: %x", r.reg, r.value, c8.V[r.reg])
				}
			}
		})
	}
}

func TestSubN(t *testing.T) {

	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  []treg
		expectErr bool
	}{
		{name: "SubN X > Y", inst: 0x87B7,
			reg: []treg{
				treg{reg: 0x07, value: 187},
				treg{reg: 0x0B, value: 25}},
			expected: []treg{
				{reg: 0x07, value: 0x5E},
				{reg: 0x0F, value: 0}},
			expectErr: false},
		{name: "SubN X == Y", inst: 0x8A27,
			reg: []treg{
				treg{reg: 0x0A, value: 255},
				treg{reg: 0x02, value: 255}},
			expected: []treg{
				{reg: 0x0A, value: 0},
				{reg: 0x0F, value: 0}},
			expectErr: false},
		{name: "SubN X < Y", inst: 0x8A27,
			reg: []treg{
				treg{reg: 0x0A, value: 13},
				treg{reg: 0x02, value: 160}},
			expected: []treg{
				{reg: 0x0A, value: 147},
				{reg: 0x0F, value: 1}},
			expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid SubN Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.SubN(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Sub: %v", err)
			}
			for _, r := range tc.expected {
				if c8.V[r.reg] != r.value {
					t.Errorf("reg: V%x; expected: %x; got: %x", r.reg, r.value, c8.V[r.reg])
				}
			}
		})
	}
}

func TestShiftLeft(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  []treg
		expectErr bool
	}{
		{name: "Shift 209 w/Carry", inst: 0x87BE,
			reg: []treg{
				treg{reg: 0x07, value: 209}},
			expected: []treg{
				{reg: 0x07, value: 0xA2},
				{reg: 0x0F, value: 1}},
			expectErr: false},
		{name: "Shift 50 w/o Carry", inst: 0x8A2E,
			reg: []treg{
				treg{reg: 0x0A, value: 50}},
			expected: []treg{
				{reg: 0x0A, value: 100},
				{reg: 0x0F, value: 0}},
			expectErr: false},
		{name: "Shift 255 w/Carry", inst: 0x8A2E,
			reg: []treg{
				treg{reg: 0x0A, value: 255}},
			expected: []treg{
				{reg: 0x0A, value: 0xFE},
				{reg: 0x0F, value: 1}},
			expectErr: false},
		{name: "Invalid Set Instruction", inst: 0x2AEB, expectErr: true},
		{name: "Invalid ShiftLeft Instruction", inst: 0x8AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.ShiftLeft(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Sub: %v", err)
			}
			for _, r := range tc.expected {
				if c8.V[r.reg] != r.value {
					t.Errorf("reg: V%x; expected: %x; got: %x", r.reg, r.value, c8.V[r.reg])
				}
			}
		})
	}
}

func TestSkipNotEqualReg(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  uint16
		expectErr bool
	}{
		{name: "Equal reg/reg.", inst: 0x91E0,
			reg: []treg{
				treg{reg: 0x01, value: 10},
				treg{reg: 0x0E, value: 10}},
			expected: PCInit, expectErr: false},
		{name: "Unequal reg/reg", inst: 0x92E0,
			reg: []treg{
				treg{reg: 0x02, value: 10},
				treg{reg: 0x0E, value: 11}},
			expected: PCInit + 1, expectErr: false},
		{name: "Invalid Instruction", inst: 0x2FFF,
			reg: []treg{
				treg{reg: 0x01, value: 10},
				treg{reg: 0x0E, value: 10}},
			expected: PCInit, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.SkipNotEqualReg(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute SkipNotEqualReg: %v", err)
			}
			if c8.PC != tc.expected {
				t.Errorf("expected: %x; got: %x", tc.expected, c8.PC)
			}
		})
	}
}

func TestLoadI(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		expected  uint16
		expectErr bool
	}{
		{name: "Load 0x0FFF", inst: 0xAFFF, expected: 0x0FFF, expectErr: false},
		{name: "Load 0x0AEB", inst: 0xAAEB, expected: 0x0AEB, expectErr: false},
		{name: "Invalid Instruction", inst: 0x2AEB, expected: 0, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			err := c8.LoadI(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute LoadValue: %v", err)
			}
			if c8.I != tc.expected {
				t.Errorf("expected: %x; got: %x", tc.expected, c8.I)
			}
		})
	}
}

func TestJumpWithReg(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		reg       []treg
		expected  uint16
		expectErr bool
	}{
		{name: "Jump 0x1E0 + 50.", inst: 0xB1E0,
			reg: []treg{
				treg{reg: 0x00, value: 50}},
			expected: 0x212, expectErr: false},
		{name: "Jump 0xFEF + 164", inst: 0xBFEF,
			reg: []treg{
				treg{reg: 0x00, value: 164}},
			expected: 0x1093, expectErr: false},
		{name: "Invalid Instruction", inst: 0x2FFF,
			reg: []treg{
				treg{reg: 0x01, value: 10},
				treg{reg: 0x0E, value: 10}},
			expected: PCInit, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			for _, r := range tc.reg {
				c8.V[r.reg] = r.value
			}
			err := c8.JumpWithReg(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute JumpWithReg: %v", err)
			}
			if c8.PC != tc.expected {
				t.Errorf("expected: %x; got: %x", tc.expected, c8.PC)
			}
		})
	}
}

func TestRandomAnd(t *testing.T) {
	tt := []struct {
		name      string
		inst      uint16
		expected  treg
		expectErr bool
	}{
		{name: "v5 - 155", inst: 0xC59B, expected: treg{reg: 0x05, value: 155}, expectErr: false},
		{name: "vA - 255", inst: 0xCAFF, expected: treg{reg: 0x0A, value: 255}, expectErr: false},
		{name: "Invalid Instruction", inst: 0x2AEB, expectErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c8 := setup()
			err := c8.RandomAnd(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute Or: %v", err)
			}
			if c8.V[tc.expected.reg] == tc.expected.value {
				t.Errorf("expected not to get: %x; got: %x", tc.expected.value, c8.V[tc.expected.reg])
			}
		})
	}
}
