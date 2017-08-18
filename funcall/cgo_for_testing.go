package funcall

// int add(int a, int b) { return a + b; }
// int mystrlen(char *s) {
//   if (!s) { return 0; }
//   int ret = 0;
//   const char* p = s;
//   while (*p++);
//   return p - s - 1;
// }
// char* foo(void) { return "foo"; }
// int sum(int* xs, int len) {
//   int ret = 0;
//   for (int i = 0; i < len; ++i) { ret += xs[i]; }
//   return ret;
// }
import "C"
import (
	"fmt"
	"unsafe"
)

// Required to do benchmarks.
// Thou can not use CGo in "*_test" files.

var (
	counter int // For side-effects to ban optimizations

	add      FuncInt
	mystrlen FuncInt
	foo      FuncString
	sum      FuncInt
)

func init() {
	_ = C.add(1, 2)
	_ = C.mystrlen(nil)
	_ = C.foo()
	_ = C.sum(nil, 0)

	invoker := NewInvoker(func(x interface{}) interface{} {
		switch x := x.(type) {
		case int:
			return C.int(x)
		case string:
			return C.CString(x)
		case []int:
			y := make([]C.int, len(x))
			for i := range x {
				y[i] = C.int(x[i])
			}
			ptr := (*C.int)(unsafe.Pointer(&y[0]))
			// Note that 1 Go arguments maps to 2 positional C args.
			return []interface{}{ptr, C.int(len(x))}
		default:
			panic(fmt.Sprintf("Go2C: unexpected %#v", x))
		}
	}, func(x interface{}) interface{} {
		switch x := x.(type) {
		case C.int:
			return int(x)
		case *C.char:
			return C.GoString(x)

		default:
			panic(fmt.Sprintf("C2Go: unexpected %#v", x))
		}
	})

	add = WrapInt(invoker, _Cfunc_add)
	mystrlen = WrapInt(invoker, _Cfunc_mystrlen)
	foo = WrapString(invoker, _Cfunc_foo)
	sum = WrapInt(invoker, _Cfunc_sum)
}

func cgoadd(x, y int) int {
	return int(C.add(C.int(x), C.int(y)))
}

func cgomystrlen(s string) int {
	return int(C.mystrlen(C.CString(s)))
}

func cgofoo() string {
	return C.GoString(C.foo())
}

func cgosum(xs []int) int {
	ys := make([]C.int, len(xs))
	for i := range xs {
		ys[i] = C.int(xs[i])
	}
	return int(C.sum((*C.int)(unsafe.Pointer(&ys[0])), C.int(len(ys))))
}
