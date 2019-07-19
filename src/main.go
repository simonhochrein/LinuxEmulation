package main

import (
	linuxemulation "LinuxEmulation"
	"debug/elf"
	"os"
)

// const (
// 	stackSize = 0x00800000
// 	stackBase = 0xbf800000
// )

// func getBytes(n uint64) []byte {
// 	buf := make([]byte, binary.MaxVarintLen64)
// 	binary.LittleEndian.PutUint64(buf, n)
// 	return buf
// }

// type application struct {
// 	f  *elf.File
// 	uc unicorn.Unicorn
// 	cs *capstr.Engine
// }

// func (a *application) pushBytes(bytes []byte) uint64 {
// 	sp, _ := a.uc.RegRead(unicorn.X86_REG_RSP)
// 	sp -= uint64(len(bytes))
// 	a.uc.RegWrite(unicorn.X86_REG_RSP, sp)
// 	a.uc.MemWrite(sp, bytes)
// 	return sp
// }

// func (a *application) addHooks() {
// 	a.uc.HookAdd(unicorn.HOOK_CODE, func(engine unicorn.Unicorn, addr uint64, size uint32) {
// 		code, _ := engine.MemRead(addr, uint64(size))
// 		instrs, _ := a.cs.Dis(code, addr, uint64(size))
// 		for _, ins := range instrs {
// 			fmt.Printf("%#x: %s %s\n", ins.Addr(), ins.Mnemonic(), ins.OpStr())
// 		}
// 	}, 1, 0)
// 	a.uc.HookAdd(unicorn.HOOK_INSN, func(engine unicorn.Unicorn) {
// 		syscalls.HandleSyscall(engine)
// 	}, 1, 0, unicorn.X86_INS_SYSCALL)
// 	a.uc.HookAdd(unicorn.HOOK_MEM_INVALID, func(engine unicorn.Unicorn, access int, address uint64, size int, value int64) bool {
// 		fmt.Printf("%x -- %d\n", address, size)
// 		return false
// 	}, 1, 0)
// }

// func (a *application) loadFile(file string) {

// 	f, _ := elf.Open(file)
// 	uc, _ := unicorn.NewUnicorn(unicorn.ARCH_X86, unicorn.MODE_64)
// 	a.cs, _ = capstr.New(capstr.ARCH_X86, capstr.MODE_64)
// 	a.f = f
// 	a.uc = uc

// 	a.addHooks()
// 	a.mapStack()

// 	for _, prog := range f.Progs {
// 		if prog.Flags&elf.PF_W != 0 {
// 			addr := f.Entry + prog.Vaddr + prog.Memsz
// 			if addr > syscalls.BRK {
// 				syscalls.BRK = addr
// 			}
// 		}

// 		if prog.Type == elf.PT_LOAD {
// 			data := make([]byte, prog.Filesz)
// 			prog.Open().Read(data)

// 			alignedStart := utils.PageStart(prog.Vaddr)
// 			alignedSize := utils.PageEnd((prog.Vaddr - alignedStart) + prog.Memsz)

// 			// fmt.Printf("Aligning Vaddr %x and Memsz %x to %x & %x\n", prog.Vaddr, prog.Memsz, alignedStart, )

// 			err := uc.MemMap(alignedStart, alignedSize)
// 			if err != nil {
// 				panic(err)
// 			}
// 			err = uc.MemWrite(prog.Vaddr, data)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}

// 	if syscalls.BRK > 0 {
// 		mask := uint64(4096)
// 		syscalls.BRK = (syscalls.BRK + mask) &^ (mask - 1)
// 	}

// 	something := make([]byte, binary.MaxVarintLen64)
// 	binary.LittleEndian.PutUint64(something, 0)
// 	a.pushBytes(something)

// 	argvp := a.pushBytes([]byte("a.out"))
// 	a.alignStack()
// 	a.pushBytes(([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
// 	a.pushBytes(something)
// 	a.pushBytes(something)
// 	a.pushBytes(getBytes(argvp))
// 	a.pushBytes(getBytes(0))
// 	a.pushBytes(getBytes(1))
// }

// func (a *application) alignStack() {
// 	sp, _ := a.uc.RegRead(unicorn.X86_REG_RSP)
// 	sp = sp &^ 15
// }

// func (a *application) mapStack() {
// 	err := a.uc.MemMap(stackBase, stackSize)
// 	if err != nil {
// 		panic(err)
// 	}

// 	a.uc.RegWrite(unicorn.X86_REG_RSP, stackBase+stackSize)
// }

type nativeFS struct {
}

type nativeInet struct {
}

func main() {
	// app := new(application)

	// app.loadFile("./a.out")

	// err := app.uc.Start(app.f.Entry, 0)
	// if err != nil {
	// 	panic(err)
	// }

	emu := linuxemulation.NewEmulator(nativeInet{}, nativeFS{})
	if len(os.Args) < 2 {
		print("Usage: ./start.sh [program]\n")
		os.Exit(1)
	}
	file, _ := elf.Open(os.Args[1])
	emu.Run(file)
}
