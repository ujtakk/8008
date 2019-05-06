package main

import (
	"bufio"
	"fmt"
)

func showIndex(inst byte, r *bufio.Reader) {
	op_hi := (inst & 0xC0) >> 6
	op_mid := (inst & 0x38) >> 3
	op_low := (inst & 0x07) >> 0
	if op_hi == 03 {
		if op_mid == 07 {
			fmt.Printf("LMr %01x\n", op_low)
		} else if op_low == 07 {
			fmt.Printf("LrM %01x\n", op_mid)
		} else {
			fmt.Printf("Lrr %01x %01x\n", op_mid, op_low)
		}
	} else if (op_low&06)>>1 == 03 {
		imm, _ := r.ReadByte()
		if op_mid == 07 {
			fmt.Printf("LMI %02x\n", imm)
		} else {
			fmt.Printf("LrI %01x %02x\n", op_mid, imm)
		}
	} else {
		if (op_low&01)>>0 == 01 {
			fmt.Printf("DCr\n")
		} else {
			fmt.Printf("INr\n")
		}
	}
}

func showAccum(inst byte, r *bufio.Reader) {
	op_hi := (inst & 0xC0) >> 6
	op_mid := (inst & 0x38) >> 3
	op_low := (inst & 0x07) >> 0
	if op_hi == 02 {
		if op_low == 07 {
			switch op_mid {
			case 00:
				fmt.Printf("ADM\n")
			case 01:
				fmt.Printf("ACM\n")
			case 02:
				fmt.Printf("SUM\n")
			case 03:
				fmt.Printf("SBM\n")
			case 04:
				fmt.Printf("NDM\n")
			case 05:
				fmt.Printf("XRM\n")
			case 06:
				fmt.Printf("ORM\n")
			case 07:
				fmt.Printf("CPM\n")
			}
		} else {
			switch op_mid {
			case 00:
				fmt.Printf("ADr %01x\n", op_low)
			case 01:
				fmt.Printf("ACr %01x\n", op_low)
			case 02:
				fmt.Printf("SUr %01x\n", op_low)
			case 03:
				fmt.Printf("SBr %01x\n", op_low)
			case 04:
				fmt.Printf("NDr %01x\n", op_low)
			case 05:
				fmt.Printf("XRr %01x\n", op_low)
			case 06:
				fmt.Printf("ORr %01x\n", op_low)
			case 07:
				fmt.Printf("CPr %01x\n", op_low)
			}
		}
	} else {
		if op_low == 04 {
			imm, _ := r.ReadByte()
			switch op_mid {
			case 00:
				fmt.Printf("ADI %02x\n", imm)
			case 01:
				fmt.Printf("ACI %02x\n", imm)
			case 02:
				fmt.Printf("SUI %02x\n", imm)
			case 03:
				fmt.Printf("SBI %02x\n", imm)
			case 04:
				fmt.Printf("NDI %02x\n", imm)
			case 05:
				fmt.Printf("XRI %02x\n", imm)
			case 06:
				fmt.Printf("ORI %02x\n", imm)
			case 07:
				fmt.Printf("CPI %02x\n", imm)
			}
		} else {
			switch op_mid {
			case 00:
				fmt.Printf("RLC\n")
			case 01:
				fmt.Printf("RRC\n")
			case 02:
				fmt.Printf("RAL\n")
			case 03:
				fmt.Printf("RAR\n")
			}
		}
	}
}

func showPCStack(inst byte, r *bufio.Reader) {
	op_hi := (inst & 0xC0) >> 6
	op_mid := (inst & 0x38) >> 3
	op_low := (inst & 0x07) >> 0
	flag := (op_mid & 04) >> 2
	cond := (op_mid & 03) >> 0
	if op_hi == 01 {
		imm0, _ := r.ReadByte()
		imm1, _ := r.ReadByte()
		imm1 &= 0x3F
		switch op_low {
		case 04:
			fmt.Printf("JMP %02x %02x\n", imm0, imm1)
		case 00:
			if flag == 01 {
				fmt.Printf("JTc %01x %02x %02x\n", cond, imm0, imm1)
			} else {
				fmt.Printf("JFc %01x %02x %02x\n", cond, imm0, imm1)
			}
		case 06:
			fmt.Printf("CAL %02x %02x\n", imm0, imm1)
		case 02:
			if flag == 01 {
				fmt.Printf("CTc %01x %02x %02x\n", cond, imm0, imm1)
			} else {
				fmt.Printf("CFc %01x %02x %02x\n", cond, imm0, imm1)
			}
		}
	} else {
		switch op_low {
		case 07:
			fmt.Printf("RET\n")
		case 05:
			fmt.Printf("RST %01x\n", op_mid)
		case 03:
			if flag == 01 {
				fmt.Printf("RTc %01x\n", cond)
			} else {
				fmt.Printf("RFc %01x\n", cond)
			}
		}
	}
}

func showInOut(inst byte) {
	rr := (inst & 0x30) >> 4
	mmm := (inst & 0x0E) >> 1
	if rr == 00 {
		fmt.Printf("INP %01x\n", mmm)
	} else {
		fmt.Printf("OUT %01x %01x\n", rr, mmm)
	}
}

func showMachine(inst byte) {
	fmt.Printf("HLT\n")
}

func show(r *bufio.Reader) {
	for inst, err := r.ReadByte(); err == nil; inst, err = r.ReadByte() {
		op_hi := (inst & 0xC0) >> 6
		op_mid := (inst & 0x38) >> 3
		op_low := (inst & 0x07) >> 0
		switch {
		case
			op_hi == 03 && (op_mid != 07 || op_low != 07),
			op_hi == 00 && op_low == 06,
			op_hi == 00 && op_mid != 00 && (op_low&06)>>1 == 00:
			showIndex(inst, r)
		case
			op_hi == 02,
			op_hi == 00 && op_low == 04,
			op_hi == 00 && (op_mid&04)>>2 == 00 && op_low == 02:
			showAccum(inst, r)
		case
			op_hi == 01 && (op_low&01)>>0 == 00,
			op_hi == 00 && op_low == 07,
			op_hi == 00 && op_low == 03,
			op_hi == 00 && op_low == 05:
			showPCStack(inst, r)
		case
			op_hi == 01 && (op_low&01)>>0 == 01:
			showInOut(inst)
		case
			op_hi == 00 && op_mid == 00 && (op_low&06)>>1 == 00,
			op_hi == 03 && op_mid == 07 && op_low == 07:
			showMachine(inst)
		default:
		}
	}
}
