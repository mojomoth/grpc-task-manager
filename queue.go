package grpctaskmanager

import "log"

type Queue []interface{}

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *Queue) Enqueue(data interface{}) {
	*q = append(*q, data)
	log.Printf("enqueue: %v\n", data)
}

func (q *Queue) Dequeue() interface{} {
	if q.IsEmpty() {
		log.Printf("queue is emplty\n")
		return nil
	}

	// get first data
	data := (*q)[0]

	// remove first data
	*q = (*q)[1:]

	log.Printf("dequeue: %v\n", data)
	return data
}
