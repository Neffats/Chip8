package chip8

import "testing"

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
