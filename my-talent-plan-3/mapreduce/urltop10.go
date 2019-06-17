package main

import (
	"strings"
	"strconv"
	"fmt"
	"sort"
	"bytes"
)

// URLTop10 .
//这里方法有两轮。
//第一轮将进行网址计数。
//第二个将对第一轮中生成的结果进行排序
//获取10个最常用的URL。
func URLTop10(nWorkers int) RoundsArgs {
	// YOUR CODE HERE :)
	// And don't forget to document your idea.
	var args RoundsArgs
	args = append(args, RoundArgs{
		MapFunc   : MyURLCountMap,
		ReduceFunc: MyURLCountReduce,
		NReduce   : nWorkers,
	})

	args = append(args, RoundArgs{
		MapFunc   : MyURLTop10Map,
		ReduceFunc: MyURLTop10Reduce,
		NReduce   : 1,
	})
	return args
}
//第一轮的map
func MyURLCountMap(filename string, contents string) []KeyValue {
	urls := strings.Split(string(contents), "\n")
	urlCnt := make(map[string] int)
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if len(url) == 0 {
			continue
		}
		urlCnt[url] ++
	}
	KeyValueMap := make([]KeyValue, 0, len(urlCnt))
	for url,cnt := range urlCnt {
		KeyValueMap = append(KeyValueMap, KeyValue{Key:url, Value: strconv.Itoa(cnt)})
	}
	return KeyValueMap
}
//第一轮的reduce
func MyURLCountReduce(key string, values []string) string {
	cnt := 0
	for _, value := range values {
		value, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		cnt += value
	}
	return fmt.Sprintf("%s %s\n",key, strconv.Itoa(cnt))
}
//第二轮的map
func MyURLTop10Map(filename string, contents string) []KeyValue {
	total := 10
	lines := strings.Split(string(contents), "\n")
	ucs := make([]*urlCount, 0, len(lines))
	KeyValuePairs := make([]KeyValue, 0, total)
	for _, value := range lines {
		tmp := strings.Split(value, " ")
		if tmp[0] != "" {
			cnt, err := strconv.Atoi(tmp[1])
			if err != nil {
				panic(err)
			}
			ucs = append(ucs, &urlCount{tmp[0], cnt})
		}
	}
	sort.Slice(ucs, func(i, j int) bool {
		if ucs[i].cnt == ucs[j].cnt {
			return ucs[i].url < ucs[j].url
		}
		return ucs[i].cnt > ucs[j].cnt
	})
	for i, u := range ucs {
		if i == total {
			break
		}
		KeyValuePairs = append(KeyValuePairs, KeyValue{Key:"", Value: fmt.Sprintf("%s %s", u.url, strconv.Itoa(u.cnt))})
	}
	return KeyValuePairs
}
//第二轮的reduce
func MyURLTop10Reduce(key string, values []string) string {
	cnts := make(map[string]int, len(values))
	for _, v := range values {
		v := strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		tmp := strings.Split(v, " ")
		n, err := strconv.Atoi(tmp[1])
		if err != nil {
			panic(err)
		}
		cnts[tmp[0]] = n
	}

	us, cs := TopN(cnts, 10)
	buf := new(bytes.Buffer)
	for i := range us {
		fmt.Fprintf(buf, "%s: %d\n", us[i], cs[i])
	}
	return buf.String()
}