package batch

import (
	"fmt"
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {

	steps, addOn := calculateStartAndEnd(n, pool)

	var wg sync.WaitGroup

	ch := make(chan user, n+5)

	data := make([]user, 0)

	for i := int64(1); i <= pool; i++ {

		wg.Add(int(i))
		start := (i * steps) - steps
		end := i * steps

		if i == steps && addOn > 0 {
			end = i*steps + addOn
		}

		go getBatchUsers(start, end, ch, &wg)

	}

	wg.Wait()
	fmt.Println(len(ch))
	for v := range ch {
		fmt.Println(v)
		data = append(data, v)
	}

	return data
}

func getBatchUsers(start int64, end int64, ch chan user, wg *sync.WaitGroup) {

	defer wg.Done()
	for i := start; i < end; i++ {
		user := getOne(i)
		ch <- user
	}

}

func calculateStartAndEnd(n int64, pool int64) (steps int64, addOn int64) {

	steps = n / pool
	addOn = n % pool

	return steps, addOn
}
