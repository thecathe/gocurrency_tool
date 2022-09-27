package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	async_job_maker_count int = 10
)

func AsyncCommunication() {

	fmt.Printf("\n\n\nAsync Tests:\n\n")

	var worker_contact chan chan int = make(chan chan int, async_job_maker_count)
	var all_jobs_complete chan = make(chan bool, async_job_maker_count)

	var waitgroup_init sync.WaitGroup
	var waitgroup_finish sync.WaitGroup
	var lock_init sync.Mutex

	lock_init.Lock()
	waitgroup_init.Add(async_job_maker_count)
	lock_init.Unlock()

	waitgroup_finish.Add(async_job_maker_count)

	go AsyncWorker(worker_contact)
	for i := 0; i < async_job_maker_count; i++ {

		go AsyncJobMaker(i, worker_contact, i, all_jobs_complete, &waitgroup_init, &lock_init)
		fmt.Printf("Async: JobMaker%d: Created\n", i)
	}

	fmt.Printf("\nAsync: Waitgroup: All JobMakers Setup\n")
	waitgroup_init.Wait()

	finished_waiting := make(chan bool)

	go func(waitgroup_finish *sync.WaitGroup, finished_waiting chan<- bool) {
		waitgroup_finish.Wait()
		fmt.Printf("\nAsync: Waitgroup: All JobMakers Finished\n")
		finished_waiting <- true
	}(&waitgroup_finish, finished_waiting)

	wait_finished := true
	c := make(chan int)
	for wait_finished {
		select {
		case <-finished_waiting:
			wait_finished = false
		case <-all_jobs_complete:
			waitgroup_finish.Done()
		case worker_contact <- c:
			fmt.Printf("Async: Bored, Send Worker Special Channel\n")
			select {
			case c <- -1:
				fmt.Printf("Async: Bored, Send Worker Special Signal\n")
				select {
				case <-c:
					fmt.Printf("Async: Bored, Recv Worker Special Signal\n")
				case <-time.After(time.Millisecond * time.Duration(100)):
					fmt.Printf("Async: Bored, Worker too slow.\n")
					<-c
					fmt.Printf("Async: Bored, Worker caught up.\n")
				}
			default:
				fmt.Printf("Async: Bored, Worker too busy.\n")
			}
		case <-time.After(time.Millisecond * time.Duration(50)):
			time.Sleep(time.Millisecond * time.Duration(10))
		}
	}

	fmt.Printf("\nAsync: All JobMakers Finished\n")
	time.Sleep(time.Millisecond * time.Duration(100))
}

func AsyncWorker(worker_contact <-chan chan int) {
	for {
		job_chan := <-worker_contact

		select {
		case job_maker_id := <-job_chan:
			select {
			case job_id := <-job_chan:

				if job_maker_id < 0 {
					job_chan <- 0
					continue
				}

				time.Sleep(time.Millisecond * time.Duration(job_id))

				if job_maker_id >= job_id*2 {
					job_chan <- 0
				} else {
					job_chan <- 1
				}
			case <-time.After(time.Millisecond * time.Duration(500)):
				continue
			}
		case <-time.After(time.Millisecond * time.Duration(500)):
			continue
		}

	}
}

func AsyncJobMaker(id int, worker_contact chan chan int, job_count int, jobs_complete chan<- bool, waitgroup_init *sync.WaitGroup, lock_init *sync.Mutex) {
	job_chan := make(chan int, id)

	lock_init.Lock()
	waitgroup_init.Done()
	lock_init.Unlock()

	waitgroup_init.Wait()
	for i := 0; i < job_count; i++ {
		worker_contact <- job_chan
		job_chan <- id
		job_chan <- i

		select {
		case response := <-job_chan:
			fmt.Printf("Async: JobMaker%d, Job%d: %d, Ontime\n", id, i, response)
		case <-time.After(time.Millisecond * time.Duration(id)):
			response := <-job_chan
			fmt.Printf("Async: JobMaker%d, Job%d: %d, Late\n", id, i, response)
		}
	}

	jobs_complete <- true
	fmt.Printf("Async: JobMaker%d: Completed, sent notice\n", id)
}
