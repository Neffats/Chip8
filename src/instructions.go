package chip8

import (
	"fmt"
	"math/rand"
	"time"
)

//CheckInst returns true if the msb of instruction matches expected.
func CheckInst(inst uint16, expected uint16) bool {
	got := (inst & 0xF000)
	return got == expected
}

//ClearScreen is the cpu wrapper for the graphics ClearScreen(). Sets each pixel to 0.
//Instruction Format: 00E0
func (c *CPU) ClearScreen(inst uint16) error {
	if check := CheckInst(inst, 0x0000); !check {
		return fmt.Errorf("received invalid ClearScreen instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x00E0 {
		return fmt.Errorf("received invalid ClearScreen instruction: %x", inst)
	}

	err := c.G.ClearScreen()
	if err != nil {
		return fmt.Errorf("could not clear the screen: %v", err)
	}

	return nil
}

//Jump PC to the specified address.
//Instruction Format: 1nnn
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
Instruction Format: 2nnn
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

//Return will set the program counter to the value at the top of the stack.
//Instruction Format: 00EE
func (c *CPU) Return(inst uint16) error {
	if check := CheckInst(inst, 0x0000); !check {
		return fmt.Errorf("received invalid Return instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x00EE {
		return fmt.Errorf("received invalid Return instruction: %x", inst)
	}

	addr, err := c.Pop()
	if err != nil {
		return fmt.Errorf("could not pop return address from stack: %v", err)
	}

	c.PC = addr

	return nil

}

//SkipEqualVal compare the specified register and value, increment PC if equal.
//Instruction Format: 3xkk
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
		c.PC += 2
	}
	return nil
}

//SkipNotEqualVal compare the specified register and value, increment PC if not equal.
//Instruction Format: 4xkk
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
		c.PC += 2
	}
	return nil
}

//SkipEqualReg compare both specified registers, increment PC if equal.
//Instruction Format: 5xy0
func (c *CPU) SkipEqualReg(inst uint16) error {
	if check := CheckInst(inst, 0x5000); !check {
		return fmt.Errorf("received invalid SkipEqualReg instruction: %x", inst)
	}
	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4

	if c.V[regX] == c.V[regY] {
		// Only increment by one since we'll increment PC again at end of main loop?
		// TODO: Need to check this.
		c.PC += 2
	}
	return nil
}

//LoadValue moves specified value into the register.
//Instruction Format: 6xkk
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
//Instruction Format: 7xkk
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
//Instruction Format: 8xy0
func (c *CPU) LoadReg(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid LoadReg instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 0 {
		return fmt.Errorf("received invalid LoadReg instruction: %x", inst)
	}
	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4
	c.V[regX] = c.V[regY]
	return nil
}

//Or - Set Vx = Vx OR Vy.
//Instruction Format: 8xy1
func (c *CPU) Or(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid Or instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 1 {
		return fmt.Errorf("received invalid Or instruction: %x", inst)
	}
	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4
	c.V[regX] = c.V[regX] | c.V[regY]
	return nil
}

//And - Set Vx = Vx AND Vy.
//Instruction Format: 8xy2
func (c *CPU) And(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid And instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 2 {
		return fmt.Errorf("received invalid And instruction: %x", inst)
	}
	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4
	c.V[regX] = c.V[regX] & c.V[regY]
	return nil
}

//Xor - Set Vx = Vx XOR Vy.
//Instruction Format: 8xy3
func (c *CPU) Xor(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid Xor instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 3 {
		return fmt.Errorf("received invalid Xor instruction: %x", inst)
	}
	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4
	c.V[regX] = c.V[regX] ^ c.V[regY]
	return nil
}

//Add - Set Vx = Vx + Vy, set VF = carry.
//Instruction Format: 8xy4
func (c *CPU) Add(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid Add instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 4 {
		return fmt.Errorf("received invalid Add instruction: %x", inst)
	}
	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4

	var result uint16
	result = uint16(c.V[regX]) + uint16(c.V[regY])

	c.V[regX] = uint8(result)
	if result > 255 {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}

	return nil
}

//Sub - Set Vx = Vx - Vy, set VF = NOT borrow.
//Instruction Format: 8xy5
func (c *CPU) Sub(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid Sub instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 5 {
		return fmt.Errorf("received invalid Sub instruction: %x", inst)
	}

	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4

	if c.V[regX] > c.V[regY] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[regX] = c.V[regX] - c.V[regY]

	return nil
}

//ShiftRight - Set Vx = Vx SHR 1.
//Instruction Format: 8xy6
func (c *CPU) ShiftRight(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid ShiftRight instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 6 {
		return fmt.Errorf("received invalid ShiftRight instruction: %x", inst)
	}
	regX := (inst & 0x0F00) >> 8

	c.V[0xF] = c.V[regX] & 1
	c.V[regX] >>= 1

	return nil
}

//SubN - Set Vx = Vy - Vx, set VF = NOT borrow
//Instruction Format: 8xy7
func (c *CPU) SubN(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid SubN instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 7 {
		return fmt.Errorf("received invalid SubN instruction: %x", inst)
	}

	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4

	if c.V[regX] < c.V[regY] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[regX] = c.V[regY] - c.V[regX]

	return nil
}

//ShiftLeft - Set Vx = Vx SHL 1.
//Instruction Format: 8xyE
func (c *CPU) ShiftLeft(inst uint16) error {
	if check := CheckInst(inst, 0x8000); !check {
		return fmt.Errorf("received invalid ShiftLeft instruction: %x", inst)
	}
	if check := inst & 0x000F; check != 0xE {
		return fmt.Errorf("received invalid ShiftLeft instruction: %x", inst)
	}
	regX := (inst & 0x0F00) >> 8

	c.V[0xF] = c.V[regX] >> 7
	c.V[regX] <<= 1

	return nil
}

//SkipNotEqualReg - Skip next instruction if Vx != Vy.
//Instruction Format: 9xy0
func (c *CPU) SkipNotEqualReg(inst uint16) error {
	if check := CheckInst(inst, 0x9000); !check {
		return fmt.Errorf("received invalid SkipNotEqualReg instruction: %x", inst)
	}

	regX := (inst & 0x0F00) >> 8
	regY := (inst & 0x00F0) >> 4

	if c.V[regX] != c.V[regY] {
		// Only increment by one since we'll increment PC again at end of main loop?
		// TODO: Need to check this.
		c.PC += 2
	}

	return nil
}

//LoadI sets the I register to the specified address.
//Instruction Format: Annn
func (c *CPU) LoadI(inst uint16) error {
	if check := CheckInst(inst, 0xA000); !check {
		return fmt.Errorf("received invalid LoadI instruction: %x", inst)
	}

	addr := inst & 0x0FFF
	c.I = addr

	return nil
}

//JumpWithReg - Jump/Set PC to location nnn + V0.
//Instruction Format: Bnnn
func (c *CPU) JumpWithReg(inst uint16) error {
	if check := CheckInst(inst, 0xB000); !check {
		return fmt.Errorf("received invalid JumpWithReg instruction: %x", inst)
	}

	addr := inst & 0x0FFF
	c.PC = addr + uint16(c.V[0])

	return nil
}

//RandomAnd takes a byte (kk) AND's it with a random byte and stores in specified reg (x).
//Instruction Format: Cxkk
func (c *CPU) RandomAnd(inst uint16) error {
	if check := CheckInst(inst, 0xC000); !check {
		return fmt.Errorf("received invalid RandomAnd instruction: %x", inst)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	v := inst & 0x00FF
	regX := (inst & 0x0F00) >> 8

	c.V[regX] = uint8(v) & uint8(rand.Int31())

	return nil
}

//DrawSprite displays n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
//Instruction Format: Dxyn
func (c *CPU) DrawSprite(inst uint16) error {
	if check := CheckInst(inst, 0xD000); !check {
		return fmt.Errorf("received invalid RandomAnd instruction: %x", inst)
	}

	x := int32((inst & 0x0F00) >> 8)
	y := int32((inst & 0x00F0) >> 4)
	size := uint8(inst & 0x000F)

	collision, err := c.G.Draw(x, y, size, c.I)
	if err != nil {
		return fmt.Errorf("could not draw sprite onto screen: %v", err)
	}

	if collision {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}

	return nil
}

//SkipIfKey skips the next instruction if specified key is pressed.
//Instruction Format: Ex9E
func (c *CPU) SkipIfKey(inst uint16) error {
	if check := CheckInst(inst, 0xE000); !check {
		return fmt.Errorf("received invalid SkipIfKey instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x009E {
		return fmt.Errorf("received invalid SkipIfKey instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8

	k, err := c.Input.IsPressed(c.V[reg])
	if err != nil {
		return fmt.Errorf("could not read key state: %v", err)
	}
	if k {
		c.PC += 2
	}

	return nil
}

//SkipIfNotKey skips the next instruction if specified key is not pressed.
//Instruction Format: ExA1
func (c *CPU) SkipIfNotKey(inst uint16) error {
	if check := CheckInst(inst, 0xE000); !check {
		return fmt.Errorf("received invalid SkipIfNotKey instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x00A1 {
		return fmt.Errorf("received invalid SkipIfNotKey instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8

	k, err := c.Input.IsPressed(c.V[reg])
	if err != nil {
		return fmt.Errorf("could not read key state: %v", err)
	}
	if !k {
		c.PC += 2
	}

	return nil
}

//WaitForKey will wait for a key to be pressed, and store the key in the specified x register.
//Instruction Format: Fx0A
func (c *CPU) WaitForKey(inst uint16) error {
	if check := CheckInst(inst, 0xF000); !check {
		return fmt.Errorf("received invalid WaitForKey instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x000A {
		return fmt.Errorf("received invalid WaitForKey instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8

	key, err := c.Input.WaitForKey()
	if err != nil {
		return fmt.Errorf("could not wait for key: %v", err)
	}

	c.V[reg] = key

	return nil
}

//SetRegDT will set the register x to the delay timers current value.
//Instruction Format: Fx07
func (c *CPU) SetRegDT(inst uint16) error {
	if check := CheckInst(inst, 0xF000); !check {
		return fmt.Errorf("received invalid SetRegDT instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x0007 {
		return fmt.Errorf("received invalid SetRegDT instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8
	t, err := c.DT.Get()
	if err != nil {
		return fmt.Errorf("could not get delay timer: %v", err)
	}

	c.V[reg] = uint8(t)
	return nil
}

//SetDT will set the DT to the value of register x.
//Instruction Format: Fx15
func (c *CPU) SetDT(inst uint16) error {
	if check := CheckInst(inst, 0xF000); !check {
		return fmt.Errorf("received invalid SetDT instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x0015 {
		return fmt.Errorf("received invalid SetDT instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8

	err := c.DT.Set(uint(c.V[reg]))
	if err != nil {
		return fmt.Errorf("could not set DT: %v", err)
	}

	return nil
}

//AddIReg will add the I register with the x register and store the result in I register.
//Instruction Format: Fx1E
func (c *CPU) AddIReg(inst uint16) error {
	if check := CheckInst(inst, 0xF000); !check {
		return fmt.Errorf("received invalid AddIReg instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x001E {
		return fmt.Errorf("received invalid AddIReg instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8

	c.I += uint16(c.V[reg])

	return nil
}

//SetISprite sets I register to the address of the hex sprite of the value in register x.
//Instruction Format: Fx29
func (c *CPU) SetISprite(inst uint16) error {
	if check := CheckInst(inst, 0xF000); !check {
		return fmt.Errorf("received invalid SetISprite instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x0029 {
		return fmt.Errorf("received invalid SetISprite instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8
	if c.V[reg] > 0xF {
		return fmt.Errorf("value in V[%x] is bigger than 0xF: %d", reg, c.V[reg])
	}
	c.I = uint16(c.V[reg] * 5)

	return nil
}

//SplitDecimal will take the decimal value from register x, split it by "column" (BCD) e.g. 137 > 1 | 3 | 7.
//Then store the hundreds digit at I, tens digit at I+1 and the ones digit at I+2.
//Instruction Format: Fx33
func (c *CPU) SplitDecimal(inst uint16) error {
	if check := CheckInst(inst, 0xF000); !check {
		return fmt.Errorf("received invalid SplitDecimal instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x0033 {
		return fmt.Errorf("received invalid SplitDecimal instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8
	dec := c.V[reg]

	hundred := (dec / 100) % 10
	ten := (dec / 10) % 10
	one := dec & 10

	c.Memory.Write(byte(hundred), c.I)
	c.Memory.Write(byte(ten), c.I+1)
	c.Memory.Write(byte(one), c.I+2)

	return nil
}

//StoreRegs will save the values of registers 0 through x to memory starting at address I.
//Instruction Format: Fx55
func (c *CPU) StoreRegs(inst uint16) error {
	if check := CheckInst(inst, 0xF000); !check {
		return fmt.Errorf("received invalid StoreRegs instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x0055 {
		return fmt.Errorf("received invalid StoreRegs instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8
	e := c.V[reg]

	for i := uint8(0); i <= e; i++ {
		err := c.Memory.Write(byte(c.V[i]), c.I+uint16(i))
		if err != nil {
			return fmt.Errorf("could not write reg[%x] to memory address %x: %v", reg, c.I+uint16(i), err)
		}
	}

	return nil
}

//LoadRegs will load the values starting at memory address I into registers 0 through x.
//Instruction Format: Fx65
func (c *CPU) LoadRegs(inst uint16) error {
	if check := CheckInst(inst, 0xF000); !check {
		return fmt.Errorf("received invalid LoadRegs instruction: %x", inst)
	}
	if check := inst & 0x00FF; check != 0x0065 {
		return fmt.Errorf("received invalid LoadRegs instruction: %x", inst)
	}

	reg := (inst & 0x0F00) >> 8
	e := c.V[reg]
	fmt.Printf("reg: %d", e)

	var err error

	for i := uint16(0); i <= uint16(e); i++ {
		addr := c.I + i
		c.V[i], err = c.Memory.Read(addr)
		if err != nil {
			return fmt.Errorf("could not read memory address %x into register[%x]: %v", c.I+uint16(i), i, err)
		}
	}

	return nil
}

//NotImplemented is a placeholder while the instructions are finished. Allows the program to emulator.
func (c *CPU) NotImplemented(inst uint16) error {
	//fmt.Printf("Instruction not implemented: %x\n", inst)
	return nil
}
