package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Job struct {
	Execution func(workerId int, data interface{}) interface{}
	Data interface{}
}

func (j Job) execute(workerId int) Result {
	return j.Execution(workerId, j.Data)
}

type Result interface{}

type WorkerPool struct {
	workersCount int
	jobs chan Job
	results chan Result
}

func NewWorkerPool(workersCount int, jobs chan Job) WorkerPool {
	return WorkerPool{
		workersCount: workersCount,
		jobs:         jobs,
		results:      make(chan Result, workersCount),
	}
}

func (wp WorkerPool) Run() {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go worker(i+1, &wg, wp.jobs, wp.results)
	}

	go func () {
		wg.Wait()
		close(wp.results)
		fmt.Println("Worker Pool finished all jobs!")
	}()
}

func (wp WorkerPool) GetResults() <-chan Result {
	return wp.results
}

func worker(workerId int, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for job := range jobs {
		results <- job.execute(workerId)
	}
}


func main() {
	workersCount := 3
	jobsCount := 10
	jobs := make(chan Job, jobsCount)

	for i := 0; i < jobsCount; i++ {
		jobs <- Job{
			Execution: func(workerId int, data interface{}) interface{} {
				val, ok := data.(int)
				if !ok {
					return "error"
				}
				return fmt.Sprintf("worker %v: %v power of 2 is: %v", workerId, val, val * val)
			},
			Data: rand.Intn(10000),
		}
	}
	close(jobs)

	wp := NewWorkerPool(workersCount, jobs)
	wp.Run()
	for i := range wp.GetResults() {
		fmt.Println(i)
	}
}



