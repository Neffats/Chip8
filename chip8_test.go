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
	return NewCPU(m)
}

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
				t.Fatalf("failed execute jump: %v", err)
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
				t.Fatalf("failed execute jump: %v", err)
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
			err := c8.LoadReg(tc.inst)
			if err != nil {
				if tc.expectErr == true {
					return
				}
				t.Fatalf("failed execute jump: %v", err)
			}
			if c8.V[tc.expected.reg] != tc.expected.value {
				t.Errorf("expected: %x; got: %x", tc.expected.value, c8.V[tc.expected.reg])
			}
		})
	}
}
