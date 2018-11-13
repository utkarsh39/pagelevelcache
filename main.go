package main

import (
	"flag"
	"fmt"
	"github.com/Onmysofa/pagelevelcache/evaluate"
	"github.com/Onmysofa/pagelevelcache/parse"
	"time"
)

func main() {

	//for i := 4; i < 7; i++ {
	//	for j := 4; j <= 7; j++ {
	//		size := math.Pow10(i)
	//		num := math.Pow10(j)
	//		qps := evaluate.EvalCcachePage(int64(size), int(num), 100, time.Minute * 10, 10)
	//		fmt.Printf("%v ", qps);
	//	}
	//	fmt.Println("")
	//}

	//for t := 1; t < 6; t++ {
	//	qps :=evaluate.EvalCcache(1000000, 10000000, 10, time.Minute * 10, t)
	//	fmt.Printf("%v ", qps);
	//}


	//for i := 4; i < 7; i++ {
	//	for j := 4; j <= 7; j++ {
	//		size := math.Pow10(i)
	//		num := math.Pow10(j)
	//		qps := evaluate.EvalGcache(int(size), int(num), 10)
	//		fmt.Printf("%v ", qps);
	//	}
	//	fmt.Println("")
	//}

	//for t := 1; t < 6; t++ {
	//	qps := evaluate.EvalGcache(1000000, 10000000, 10)
	//	fmt.Printf("%v ", qps);
	//}

	//evaluate.EvalGcache(1000, 1000000, 10)

	//ch, err := parse.ParseFile("/home/ruogu/Desktop/capstone/data/first1000.json")
	//if err != nil {
	//	return
	//}
	//for i := 7; i < 10; i++ {
	//	size := math.Pow10(i)
	//	qps := evaluate.EvalCcacheTrace(ch, int64(size), 1000, 100, time.Minute * 10, 10)
	//	fmt.Printf("%v ", qps);
	//	fmt.Println("")
	//}

	//parse.ParititionFile("/home/ruogu/Desktop/capstone/data/trace_2018_03_06_24h.json", 8)

	//argsWithoutProg := os.Args[1:]
	//size, err := strconv.ParseInt(argsWithoutProg[1], 10, 64)
	//if err != nil {
	//	return
	//}
	//threads, err := strconv.ParseInt(argsWithoutProg[2], 10, 64)
	//funBenchTrace(argsWithoutProg[0], size, int(threads))
	//funCalcUniqueSize(argsWithoutProg[0])

	tracePtr := flag.String("t", "", "The trace to test against")
	sizePtr := flag.Int64("s", 5000000000, "The size of cache by byte")
	granularityPtr := flag.Int("g", 30000, "The granularity to report")

	ohrPtr := flag.Bool("o", false, "Calculate OHR")
	phrPtr := flag.Bool("p", true, "Calculate PHR")

	flag.Parse()

	if *ohrPtr {
		funBenchTraceOHR(*tracePtr, *granularityPtr, *sizePtr)
	}

	if *phrPtr {
		funBenchTracePHR(*tracePtr, *granularityPtr, *sizePtr)
	}

}

func funCalcSize(filename string) {
	chs, err := parse.ParseFile(filename, 1)
	if err != nil {
		return
	}

	res := calcSizeSum(chs[0])
	fmt.Print("Size sum: ", res)
}

func funCalcUniqueSize(filename string) {
	chs, err := parse.ParseFile(filename, 1)
	if err != nil {
		return
	}

	res := calcUniqueSize(chs[0])
	fmt.Print("Unique size sum: ", res)
}

func funCalcNum(filename string) {
	chs, err := parse.ParseFile(filename, 0)
	if err != nil {
		return
	}

	start := time.Now()

	res := calcNum(chs[0])

	duration := time.Now().Sub(start)
	qps := float64(res)/ duration.Seconds()

	fmt.Println("Number: ", res)
	fmt.Println("QPS: ", qps)
}

func funBenchTraceThroughtput(filename string, size int64, threads int) {
	chs, err := parse.ParseFileWithoutValue(filename, 1)
	if err != nil {
		return
	}

	fmt.Println("Size: ", size, " Threads: ", threads)
	num := calcNum(chs[0])
	fmt.Println("Num: ", num)

	chs, err = parse.ParseFileWithoutValue(filename, threads)
	if err != nil {
		return
	}

	fmt.Print("Wait 60s for parsing...")
	time.Sleep(60 * time.Second)
	fmt.Println("")

	qps := evaluate.EvalCcacheTrace(chs, size, num, 100, time.Minute * 10, threads)
	fmt.Printf("%v ", qps);
	fmt.Println("")

}

func funBenchTracePHR(filename string, granularity int, size int64) {

	fmt.Println("Trace:", filename)
	fmt.Println("Granularity:", granularity)
	fmt.Println("Cache size:", size)
	
	chs, err := parse.ParseFileWithoutValue(filename, 1)
	if err != nil {
		return
	}

	ratio := evaluate.EvalCcachePHR(chs[0], granularity, "LFU", size,100, time.Minute * 10)
	fmt.Printf("Ratio: %v\n ", ratio);
}

func funBenchTraceOHR(filename string, granularity int, size int64) {

	fmt.Println("Trace:", filename)
	fmt.Println("Granularity:", granularity)
	fmt.Println("Cache size:", size)

	chs, err := parse.ParseFileWithoutValue(filename, 1)
	if err != nil {
		return
	}

	ratio := evaluate.EvalCcacheOHR(chs[0], granularity, "LFU", size,100, time.Minute * 10)
	fmt.Printf("Ratio: %v\n ", ratio);
}

func calcSizeSum(ch chan *parse.PageReq) int64 {
	var sum int64 = 0
	for r := range ch {
		for _, o := range r.Objs {
			sum += int64(o.Size)
		}
	}

	return sum
}

func calcUniqueSize(ch chan *parse.PageReq) int64 {
	var sum int64 = 0
	m := make(map[string]bool)

	for r := range ch {
		for _, o := range r.Objs {
			k := fmt.Sprintf("%v:%v", o.Backend, o.Uri)
			_,ok := m[k]
			if !ok {
				sum += int64(o.Size)
				m[k] = true
			}
		}
	}

	return sum
}

func calcNum(ch chan *parse.PageReq) int {
	sum := 0
	for range ch {
		sum++
	}

	return sum
}

