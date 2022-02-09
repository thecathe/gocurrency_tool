package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	sync_job_maker_count int = 10
)

func SyncCommunication() {

	fmt.Printf("\n\n\nSync Tests:\n\n")

	worker_contact := make(chan chan int)
	var all_jobs_complete [sync_job_maker_count]chan bool
	var finished_job_makers [sync_job_maker_count]bool

	var waitgroup_init sync.WaitGroup
	var waitgroup_finish sync.WaitGroup
	var lock_init sync.Mutex
	var lock_finish sync.Mutex

	lock_init.Lock()
	waitgroup_init.Add(sync_job_maker_count)
	lock_init.Unlock()

	lock_finish.Lock()
	waitgroup_finish.Add(sync_job_maker_count)
	lock_finish.Unlock()

	go SyncWorker(worker_contact)
	for i := 0; i < sync_job_maker_count; i++ {

		all_jobs_complete[i] = make(chan bool)
		finished_job_makers[i] = false

		go SyncJobMaker(i, worker_contact, i, all_jobs_complete[i], &waitgroup_init, &waitgroup_finish, &lock_init, &lock_finish)
		fmt.Printf("Sync: JobMaker%d: Created\n", i)
	}

	fmt.Printf("\nSync: Waitgroup: All JobMakers Setup\n")
	waitgroup_init.Wait()

	finished_waiting := make(chan bool)

	go func(waitgroup_finish *sync.WaitGroup, finished_waiting chan<- bool) {
		waitgroup_finish.Wait()
		fmt.Printf("\nSync: Waitgroup: All JobMakers Finished\n")
		finished_waiting <- true
	}(&waitgroup_finish, finished_waiting)

	wait_finished := true
	wait_checked := true
	for wait_finished || wait_checked {
		select {
		case <-finished_waiting:
			wait_finished = false
		case <-time.After(time.Millisecond * time.Duration(500)):
			all_passed, _finished_job_makers := func(all_jobs_complete [sync_job_maker_count]chan bool, finished_job_makers [sync_job_maker_count]bool) (bool, [sync_job_maker_count]bool) {
				all_passed := true
				for i, job_maker := range all_jobs_complete {

					// check if this has been closed
					if !finished_job_makers[i] {
						all_passed = false
						select {
						case <-job_maker:
							// time.Sleep(time.Millisecond * time.Duration(10))
							fmt.Printf("Sync: JobMaker%d: Finished\n", i)
							finished_job_makers[i] = true
						default:
							// fmt.Printf("JobMaker%d: Passed\n", i)
							continue
						}
					}
				}
				return all_passed, finished_job_makers
			}(all_jobs_complete, finished_job_makers)

			finished_job_makers = _finished_job_makers
			if all_passed {
				wait_checked = false
				fmt.Printf("\nSync: BoolArray: All JobMakers Finished\n")
			}
		}
	}

	fmt.Printf("\nSync: All JobMakers Finished\n")
	time.Sleep(time.Millisecond * time.Duration(100))
}

func SyncWorker(worker_contact <-chan chan int) {
	for {
		job_chan := <-worker_contact

		job_maker_id := <-job_chan
		job_id := <-job_chan

		time.Sleep(time.Millisecond * time.Duration(job_id))

		if job_maker_id >= job_id*2 {
			job_chan <- 0
		} else {
			job_chan <- 1
		}
	}
}

func SyncJobMaker(id int, worker_contact chan chan int, job_count int, jobs_complete chan<- bool, waitgroup_init *sync.WaitGroup, waitgroup_finish *sync.WaitGroup, lock_init *sync.Mutex, lock_finish *sync.Mutex) {
	job_chan := make(chan int)

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
			fmt.Printf("Sync: JobMaker%d, Job%d: %d, Ontime\n", id, i, response)
		case <-time.After(time.Millisecond * time.Duration(id)):
			response := <-job_chan
			fmt.Printf("Sync: JobMaker%d, Job%d: %d, Late\n", id, i, response)
		}
	}

	lock_finish.Lock()
	waitgroup_finish.Done()
	lock_finish.Unlock()

	jobs_complete <- true
	fmt.Printf("Sync: JobMaker%d: Completed, sent notice\n", id)
}
