package processors

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/expinc/melegraf/config"
	"github.com/expinc/melegraf/conveyor"
	"github.com/expinc/melegraf/metric"
	"github.com/expinc/melegraf/processor"
)

func TestDummyProcessorOrdinary(t *testing.T) {
	// prepare config string
	cfgStr := `
	{
		"name": "dummy_processor",
		"type": "dummy",
		"cronSpec": "@every 1s",
		"params": {
			"propValue": "dummy_value"
		}
	}`
	var cfg config.ProcessorConfig
	err := json.Unmarshal([]byte(cfgStr), &cfg)
	if err != nil {
		t.Error(err)
	}

	// create dummy processor
	proc, err := processor.NewProcessorRunner(ProcessorTypeDummy, &cfg)
	if err != nil {
		t.Error(err)
	}

	// set inputs & outputs
	input1 := conveyor.NewConveyor("input1", 10)
	input1.Put(metric.Metric{Name: "dummy_metric"})
	input2 := conveyor.NewConveyor("input2", 10)
	input2.Put(metric.Metric{Name: "dummy_metric"})
	input2.Put(metric.Metric{Name: "dummy_metric"})
	output1 := conveyor.NewConveyor("output1", 10)
	output2 := conveyor.NewConveyor("output2", 10)
	err = proc.AddInput(input1)
	if err != nil {
		t.Error(err)
	}
	err = proc.AddInput(input2)
	if err != nil {
		t.Error(err)
	}
	err = proc.AddOutput(output1)
	if err != nil {
		t.Error(err)
	}
	err = proc.AddOutput(output2)
	if err != nil {
		t.Error(err)
	}

	// start for some time
	err = proc.Start()
	if err != nil {
		t.Error(err)
	}
	runSeconds := 5
	time.Sleep(time.Duration(runSeconds) * time.Second)
	err = proc.Stop()
	if err != nil {
		t.Error(err)
	}

	// check if the processor has received and sent the metrics
	cntCronMetrics1 := 0
	cntCronMetrics2 := 0
	cntReceiveMetrics1 := 0
	cntReceiveMetrics2 := 0
	noMoreMetrics := false
	for {
		select {
		case mt := <-output1.GetChannel():
			dummyTag, err := mt.GetTag("dummy")
			if err != nil {
				t.Error(err)
			}
			if dummyTag == "cron" {
				cntCronMetrics1++
			} else if dummyTag == "received" {
				cntReceiveMetrics1++
			} else {
				t.Error("invalid metric received")
			}
		case mt := <-output2.GetChannel():
			dummyTag, err := mt.GetTag("dummy")
			if err != nil {
				t.Error(err)
			}
			if dummyTag == "cron" {
				cntCronMetrics2++
			} else if dummyTag == "received" {
				cntReceiveMetrics2++
			} else {
				t.Error("invalid metric received")
			}
		default:
			noMoreMetrics = true
		}
		if noMoreMetrics {
			break
		}
	}
	if cntCronMetrics1 != runSeconds {
		t.Errorf("invalid number of metrics received in outpu1: %d", cntCronMetrics1)
	}
	if cntCronMetrics2 != runSeconds {
		t.Errorf("invalid number of metrics received in outpu2: %d", cntCronMetrics2)
	}
	if cntReceiveMetrics1 != 3 {
		t.Errorf("invalid number of metrics received in outpu1: %d", cntReceiveMetrics1)
	}
	if cntReceiveMetrics2 != 3 {
		t.Errorf("invalid number of metrics received in outpu2: %d", cntReceiveMetrics2)
	}
}
