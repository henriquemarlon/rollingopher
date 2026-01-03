package main

import (
	"log"

	"github.com/henriquemarlon/rollingopher/pkg/rollup"
)

func handleAdvance(r *rollup.Rollup) bool {
	advance, err := r.ReadAdvanceState()
	if err != nil {
		log.Printf("[echo] failed to read advance: %v", err)
		return false
	}

	log.Printf("[echo] received advance from %s with %d bytes", advance.MsgSender.Hex(), len(advance.Payload))

	_, err = r.EmitNotice(advance.Payload)
	if err != nil {
		log.Printf("[echo] failed to emit notice: %v", err)
		return false
	}

	log.Printf("[echo] emitted notice with payload")
	return true
}

func handleInspect(r *rollup.Rollup) bool {
	inspect, err := r.ReadInspectState()
	if err != nil {
		log.Printf("[echo] failed to read inspect: %v", err)
		return false
	}

	log.Printf("[echo] received inspect with %d bytes", len(inspect.Payload))

	err = r.EmitReport(inspect.Payload)
	if err != nil {
		log.Printf("[echo] failed to emit report: %v", err)
		return false
	}

	log.Printf("[echo] emitted report with payload")
	return true
}

func main() {
	r, err := rollup.New()
	if err != nil {
		log.Fatalf("[echo] failed to create rollup: %v", err)
	}
	defer r.Close()

	accept := true
	for {
		reqType, _, err := r.Finish(accept)
		if err != nil {
			log.Printf("[echo] finish error: %v", err)
			continue
		}

		switch reqType {
		case rollup.RequestTypeAdvance:
			accept = handleAdvance(r)
		case rollup.RequestTypeInspect:
			accept = handleInspect(r)
		}
	}
}
