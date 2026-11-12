#include "textflag.h"

// func add(a, b uint32) uint32
TEXT ·add(SB),NOSPLIT,$0-12
	MOVL a+0(FP), AX
	ADDL b+4(FP), AX
	MOVL AX, ret+8(FP)
	RET
