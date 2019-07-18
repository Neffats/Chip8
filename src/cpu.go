package chip8

import (
	"encoding/binary"
	"fmt"
)

const (
	//PCInit is the inital value for the PC.
	PCInit = uint16(512)
	//SPInit is the inital value for the SP.
	SPInit = uint8(16)
)

//CPU for Chip8 VM.
type CPU struct {
	PC uint16
	I  uint16

	G *Graphics

	Stack [16]uint16
	SP    uint8

	Memory *Memory

	// Registers
	V [16]uint8
}

//NewCPU returns a new CPU blank struct.
func NewCPU(m *Memory, g *Graphics) *CPU {
	return &CPU{
		PC:     PCInit,
		I:      0,
		G:      g,
		SP:     SPInit,
		Memory: m,
	}
}

//Fetch the instruction the PC is currently pointing at.
func (c *CPU) Fetch() (uint16, error) {
	var inst []byte

	//Retrieve first byte of instruction.
	i, err := c.Memory.Read(c.PC)
	if err != nil {
		return binary.BigEndian.Uint16(inst), fmt.Errorf("could not fetch first byte of instruction: %v", err)
	}
	inst = append(inst, i)

	//Retrieve last byte of instruction.
	i, err = c.Memory.Read(c.PC)
	if err != nil {
		return binary.BigEndian.Uint16(inst), fmt.Errorf("could not fetch last byte of instruction: %v", err)
	}
	inst = append(inst, i)

	return binary.BigEndian.Uint16(inst), nil
}

//Decode the instruction and return the relevant handler function.
func (c *CPU) Decode(inst uint16) (func() error, error) {
	switch op := inst & 0xF000; op {
	//Jump to location nnn.
	case 0x1000:
		return func() error { return c.Jump(inst) }, nil
	//Call subroutine at nnn.
	case 0x2000:
		return func() error { return c.Call(inst) }, nil
	//Skip next instruction if Vx = kk.
	case 0x3000:
		return func() error { return c.SkipEqualVal(inst) }, nil
	//Skip next instruction if Vx != kk.
	case 0x4000:
		return func() error { return c.SkipNotEqualVal(inst) }, nil
	//Skip next instruction if Vx = Vy.
	case 0x5000:
		return func() error { return c.SkipEqualReg(inst) }, nil
	//Set Vx = kk.
	case 0x6000:
		return func() error { return c.LoadValue(inst) }, nil
	//Set Vx = Vx + kk.
	case 0x7000:
		return func() error { return c.AddValue(inst) }, nil
	//Set instructions.
	case 0x8000:
		switch t := inst & 0x000F; t {
		//Set Vx = Vy.
		case 0x0000:
			return func() error { return c.LoadReg(inst) }, nil
		//Set Vx = Vx OR Vy.
		case 0x0001:
			return func() error { return c.Or(inst) }, nil
		//Set Vx = Vx AND Vy.
		case 0x0002:
			return func() error { return c.And(inst) }, nil
		//Set Vx = Vx XOR Vy.
		case 0x0003:
			return func() error { return c.Xor(inst) }, nil
		//Set Vx = Vx + Vy, set VF = carry.
		case 0x0004:
			return func() error { return c.Add(inst) }, nil
		//Set Vx = Vx - Vy, set VF = NOT borrow.
		case 0x0005:
			return func() error { return c.Sub(inst) }, nil
		//Set Vx = Vx SHR 1.
		case 0x0006:
			return func() error { return c.ShiftRight(inst) }, nil
		//Set Vx = Vy - Vx, set VF = NOT borrow.
		case 0x0007:
			return func() error { return c.SubN(inst) }, nil
		//Set Vx = Vx SHL 1.
		case 0x000E:
			return func() error { return c.ShiftLeft(inst) }, nil
		}
	//Skip next instruction if Vx != Vy.
	case 0x9000:
		return func() error { return c.SkipNotEqualReg(inst) }, nil
	//Set I = nnn.
	case 0xA000:
		return func() error { return c.LoadI(inst) }, nil
	//Jump to location nnn + V0.
	case 0xB000:
		return func() error { return c.JumpWithReg(inst) }, nil
	//Set Vx = random byte AND kk.
	case 0xC000:
		return func() error { return c.RandomAnd(inst) }, nil
	//Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
	case 0xD000:
		return func() error { return c.DrawSprite(inst) }, nil
	}
	return nil, fmt.Errorf("invalid instruction: %x", inst)
}

//Push data onto stack and decrement stack pointer.
func (c *CPU) Push(data uint16) error {
	c.SP--
	//Check if the stack is full.
	if int(c.SP) >= len(c.Stack) {
		return fmt.Errorf("reached stack limit")
	}
	c.Stack[c.SP] = data
	return nil
}

//Pop top value off stack and increment stack pointer.
func (c *CPU) Pop() (uint16, error) {
	if c.SP > SPInit {
		return 0, fmt.Errorf("no data in stack to pop")
	}
	data := c.Stack[c.SP]
	c.SP++

	return data, nil
}
