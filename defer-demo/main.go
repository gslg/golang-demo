package main

import (
	"fmt"
	"runtime"
)

//https://www.ardanlabs.com/blog/2013/06/understanding-defer-panic-and-recover.html
func simulateError(key string) error {
	return fmt.Errorf("这是一个错误:%s", key)
}

func simulatePanic(key string) {
	panic("我要崩溃了:" + key)
}
//这里main中a=10
func test() int {
	fmt.Println("Start Test.......")

	a := 10
	fmt.Println("a的地址是:",&a)
	err := simulateError("1")

	defer func() {
		fmt.Println("Start Defer.....")
		if err != nil {
			fmt.Println("Defer Error:", err)
		}
		a = 11
		fmt.Println("a的地址是:",&a)
	}()

	fmt.Println("End Test...........")

	return a
}

//这里main中a会等于11
func test2() (a int) {
	fmt.Println("Start Test.......")

	a = 10

	err := simulateError("1")
	fmt.Println("err address:", &err)

	defer func() {
		fmt.Println("Start Defer.....")
		fmt.Println("err address:", &err)
		if err != nil {
			fmt.Println("Defer Error:", err)
		}
		a = 11

	}()

	err = simulateError("2")

	fmt.Println("End Test...........")

	return a
}

func test3() (err error) {
	defer func() {
		fmt.Println("Start Defer.....")
		fmt.Println("err address:", &err)
		if err != nil {
			fmt.Println("Defer Error:", err)
		}
	}()

	fmt.Println("Start Test.......")

	err = simulateError("3")

	fmt.Println("err address:", &err)
	fmt.Println("end Test.......")

	return err
}

func test4() {
	fmt.Println("start test4")
	simulatePanic("4")
	fmt.Println("end test4")
}

func test5() {

	defer func() {
		fmt.Println("test5 defer: ")
	}()
	fmt.Println("start test5")

	simulatePanic("5")

	fmt.Println("end test5")
}

func test6() {
	defer func() {
		fmt.Println("start test6 defer....")

		if r := recover(); r != nil {
			fmt.Println("Defer panic:", r)
		}

		fmt.Println("end test6 defer")

	}()

	fmt.Println("start test6")

	simulatePanic("6")

	fmt.Println("end test6")
}

//Go will execute these inline defer functions in the opposite order that are defined (First In - Last Out).
func test7() (err error) {
	defer func() {
		fmt.Println("panic defer start *************")
		if r := recover(); r != nil {
			fmt.Println("Defer panic: ", r)
		}
		fmt.Println("panic  defer end. ***************")
	}()

	defer func() {
		fmt.Println("error defer start ###########")

		if err != nil {
			fmt.Println("Defer error: ", err)
		}

		fmt.Println("error defer end ###########")
	}()
	fmt.Println("start test7")
	err = simulateError("哈哈错误")
	//simulatePanic("2")
	panic("我panic了")
}

//在上面test7的过程中我们在defer中处理panic,但是main中并不知道这点，我们需要将造成panic的异常告诉main
func test8() (err error) {
	defer func() {
		fmt.Println("panic defer start *************")
		if r := recover(); r != nil {
			fmt.Println("Defer panic: ", r)
			err = fmt.Errorf("%v", r)
		}
		fmt.Println("panic  defer end. ***************")
	}()

	defer func() {
		fmt.Println("error defer start ###########")

		if err != nil {
			fmt.Println("Defer error: ", err)
		}

		fmt.Println("error defer end ###########")
	}()
	fmt.Println("start test7")
	err = simulateError("哈哈错误")
	simulatePanic("2")
	fmt.Println("end test7")

	return err
}

//test8用两个defer处理了error和panic错误，并报告给了main,但是在实际中并不优雅，我们需要一个统一处理error和panic的地方
func catchPanic(err *error, functionName string) {
	//fmt.Println("Err Addr Defer:", &err)
	if r := recover(); r != nil {
		fmt.Printf("%s : PANIC Defered : %v.\n", functionName, r)

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		fmt.Printf("%s : Stack Trace : %s.\n", functionName, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		fmt.Printf("%s : ERROR : %v.\n", functionName, *err)
	}

}

func test9() (err error){
	defer catchPanic(&err,"Test9")
	fmt.Println("Start Test9")

	err = simulateError("1")

	fmt.Println("End Test9")

	return err
}

func test10() (err error){
	fmt.Println("Err Addr:", &err)
	defer catchPanic(&err,"Test10")
	fmt.Println("Start Test10")

	err = simulateError("10")
	if err != nil{
		catchPanic(&err,"Test10")
	}

	fmt.Println("End Test10")

	simulatePanic("我是一个panic")
	return err
}

//但是，最好还是只使用这个捕获panic,而将错误日志交给应用的逻辑去处理，
// 否则像上面那样的话可能会重复记录错误两次
func catchPanicOnly(err *error, functionName string) {
	if r := recover(); r != nil {
		fmt.Printf("%s : PANIC Defered : %v\n", functionName, r)

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

func main() {
	/*a := test()

	fmt.Println("a的值是:",a)*/
	//a := test2()
	//
	//fmt.Println("a的值是:",a)

	/*if err := test3(); err != nil {
		fmt.Printf("Main Error:%v\n", err)
	}*/

	/*test5()
	fmt.Println("main 测试完毕.....")*/
	//go test5() go routine中的panic也会导致main goroutine中断
	//time.Sleep(time.Second)
	//test6()
	//test7()
	//fmt.Println("main 测试完毕.....")

	//if err := test8(); err != nil {
	//	fmt.Printf("Main Error:%v\n", err)
	//}
	/*if err := test9(); err != nil {
		fmt.Printf("Main Error:%v\n", err)
	}*/
	if err := test10(); err != nil {
		fmt.Printf("Main Error:%v\n", err)
	}

}
