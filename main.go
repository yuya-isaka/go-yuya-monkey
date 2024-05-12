package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/yuya-isaka/go-yuya-monkey/repl"
)

func main() {
	// getpwnam_r() と getpwuid_r() 関 数 は 、 そ れ ぞ れ getpwnam() と getpwuid() と 同 じ 情 報 を 取 得 す る が 、 取 得 し た passwd 構 造 体 を pwd が 指 す 領 域 に 格 納 す る 。 passwd 構 造 体 の メ ン バ ー が 指 す 文 字 列 は 、 サ イ ズ buflen の バ ッ フ ァ ー buf に 格 納 さ れ る 。 成 功 し た 場 合 *result に は 結 果 へ の ポ イ ン タ ー が 格 納 さ れ る 。 エ ン ト リ ー が 見 つ か ら な か っ た 場 合 や エ ラ ー が 発 生 し た 場 合 に は *result に は NULL が 入 る 。 呼 び 出 し
	// Current()
	// current()
	// lookupUnixUid(syscall.Getuid())

	// getuidシステムコールを利用して現在のプロセスのUIDを取得し、整数値として返しています。
	// 									rawSyscall(abi.FuncPCABI0(libc_getuid_trampoline), 0, 0, 0)

	// rawSyscall:
	// rawSyscallは、システムコールを直接呼び出すためのGoランタイムの内部関数です。
	// この関数を使うことで、システムコールを介してオペレーティングシステムの低レベルの機能にアクセスできます。

	// libc_getuid_trampoline:
	// libc_getuid_trampolineは、C言語の標準ライブラリにあるgetuidシステムコールをラップするためのトランポリン関数（中間関数）です。
	// abi.FuncPCABI0でこのトランポリン関数のアドレスを取得し、rawSyscallの第一引数として渡します。

	// rawSyscallの第二引数、第三引数、第四引数にはすべて0が渡されます。これはgetuidシステムコールに引数が不要なためです。

	// トランポリン関数とは、システムコールをラップしてくれている
	// rawSyscall(abi.FuncPCABI0(libc_getuid_trampoline), 0, 0, 0)
	// libc_getuid_trampolineトランポリン関数で、getuidシステムコールをラップしている
	// abi.FuncPCABI0でトランポリン関数のエントリポイント（PC）を取得
	// システムコールに、そのエントリポイントとなるアドレスと引数を渡して呼び出し
	// これでシステムコール呼び出しを実現している
	// システムコールのラッパー（トランポリン）、そのシステムコールのラッパーのアドレスを取得する関数、システムコールを呼び出す関数
	// _C_getpwuid_r
	// unix.Getpwuid(uid, &pwd, buf, size, &result)
	// syscall_syscall6(abi.FuncPCABI0(libc_getpwuid_r_trampoline),
	// 		uintptr(uid),
	// 		uintptr(unsafe.Pointer(pwd)),
	// 		uintptr(unsafe.Pointer(buf)),
	// 		size,
	// 		uintptr(unsafe.Pointer(result)),
	// 		0)

	user, err := user.Current()
	if err != nil {
		// deferを実行する
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
