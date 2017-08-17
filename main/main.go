package main

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
	"funcall"
	"unsafe"
)

func main() {
	// 1. Force CGO to produce these symbols.
	//    If we reference a function via "C.X" syntax,
	//    no metadata is accessible, we get plain "unsafe.Pointer".
	_ = C.add(1, 2)
	_ = C.mystrlen(nil)
	_ = C.foo()
	_ = C.sum(nil, 0)

	// 2. Because "X.C.int" and "Y.C.int" are incompatible,
	//    client code must map C types itself.
	fun := funcall.NewInvoker(func(x interface{}) interface{} {
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

	// 3. Now it is possible to use "_Cfunc<X>" syntax to call
	//    function via "reflect" package, thanks to additional
	//    info that is stored there.
	//    Arguments may be any Go types as long as they have
	//    convertion case defined.
	fmt.Println(fun.Int(_Cfunc_add, 1, 2))
	fmt.Println(fun.Int(_Cfunc_mystrlen, "1"))
	fmt.Println(fun.String(_Cfunc_foo))
	fmt.Println(fun.Int(_Cfunc_sum, []int{1, 2, 3}))
	fmt.Println(fun.Call(_Cfunc_sum, []int{1, 2, 3, 4}).(int))
}
