package chip8

import "fmt"

//CheckInst returns true if the msb of instruction matches expected.
func CheckInst(inst uint16, expected uint16) bool {
	got := (inst & 0xF000)
	return got == expected
}

//Jump PC to the specified address.
func (c *CPU) Jump(inst uint16) error {
	if check := CheckInst(inst, 0x1000); !check {
		return fmt.Errorf("received invalid Jump instruction: %x", inst)
	}
	var addr uint16
	addr = inst & 0x0FFF
	c.PC = addr
	return nil
}

/*
Call pushes current PC value onto the stack and
sets the PC to the specified address.
*/
func (c *CPU) Call(inst uint16) error {
	if check := CheckInst(inst, 0x2000); !check {
		return fmt.Errorf("received invalid Call instruction: %x", inst)
	}
	var addr uint16
	addr = inst & 0x0FFF
	err := c.Push(c.PC)
	if err != nil {
		return fmt.Errorf("could not push address: %v", err)
	}
	c.PC = addr
	return nil
}

//SkipEqualVal compare the specified register and value, increment PC if equal.
func (c *CPU) SkipEqualVal(inst uint16) error {
	if check := CheckInst(inst, 0x3000); !check {
		return fmt.Errorf("received invalid SkipEqualVal instruction: %x", inst)
	}
	reg := (inst & 0x0F00) >> 8
	//Convert to uint8 for comparison with reg.
	val := uint8(inst & 0x00FF)
	if c.V[reg] == val {
		// Only increment by one since we'll increment PC again at end of main loop?
		// TODO: Need to check this.
		c.PC++
	}
	return nil
}

//SkipNotEqualVal compare the specified register and value, increment PC if not equal.
func (c *CPU) SkipNotEqualVal(inst uint16) error {
	if check := CheckInst(inst, 0x4000); !check {
		return fmt.Errorf("received invalid SkipNotEqualVal instruction: %x", inst)
	}
	reg := (inst & 0x0F00) >> 8
	//Convert to uint8 for comparison with reg.
	val := uint8(inst & 0x00FF)
	if c.V[reg] != val {
		// Only increment by one since we'll increment PC again at end of main loop?
		// TODO: Need to check this.
		c.PC++
	}
	return nil
}

//SkipEqualReg compare both specified registers, increment PC if equal.
func (c *CPU) SkipEqualReg(inst uint16) error {
	if check := CheckInst(inst, 0x5000); !check {
		return fmt.Errorf("received invalid SkipEqualReg instruction: %x", inst)
	}
	regx := (inst & 0x0F00) >> 8
	regy := (inst & 0x00F0) >> 4

	if c.V[regx] == c.V[regy] {
		// Only increment by one since we'll increment PC again at end of main loop?
		// TODO: Need to check this.
		c.PC++
	}
	return nil
}

//LoadValue moves specified value into the register.
func (c *CPU) LoadValue(inst uint16) error {
	if check := CheckInst(inst, 0x6000); !check {
		return fmt.Errorf("received invalid LoadValue instruction: %x", inst)
	}
	r := (inst & 0x0F00) >> 8
	v := uint8(inst & 0x00FF)

	c.V[r] = v

	return nil
}

//AddValue adds specified value into the register.
func (c *CPU) AddValue(inst uint16) error {
	if check := CheckInst(inst, 0x7000); !check {
		return fmt.Errorf("received invalid AddValue instruction: %x", inst)
	}
	r := (inst & 0x0F00) >> 8
	v := uint8(inst & 0x00FF)

	c.V[r] += v

	return nil
}

//LoadReg sets the value of reg x to the value of reg y.
func (c *CPU) LoadReg(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid LoadReg instruction: %x", inst)
	}
	regx := (inst & 0x0F00) >> 8
	regy := (inst & 0x00F0) >> 4
	c.V[regx] = c.V[regy]
	return nil
}
