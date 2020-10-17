package vm

import (
	"testing"
	"github.com/holiman/uint256"
	"github.com/ethereum/go-ethereum/params"
)

func BenchmarkStackPush(b *testing.B) {
	var (
		env            = NewEVM(Context{}, nil, params.TestChainConfig, Config{})
		stack          = newstack()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	value := new(uint256.Int).SetUint64(0x1337)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.push(value)
	}
}

func BenchmarkStackPop(b *testing.B) {
	var (
		env            = NewEVM(Context{}, nil, params.TestChainConfig, Config{})
		stack          = newstack()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	value := new(uint256.Int).SetUint64(0x1337)

	// TODO how to make sure that b.N is limited to 1024?  why does this crash when the stack is massive?

	for i := 0; i < b.N; i++ {
		stack.push(value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.pop()
	}
}
