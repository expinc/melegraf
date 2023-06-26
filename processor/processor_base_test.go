package processor

import (
	"testing"
	"time"

	"github.com/expinc/melegraf/config"
	"github.com/expinc/melegraf/conveyor"
	"github.com/expinc/melegraf/metric"
	"github.com/stretchr/testify/assert"
)

var cfg = &config.ProcessorConfig{
	Name: "proc",
	Type: "dummy",
}

func TestAddInput_Succeed(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	input1 := conveyor.NewConveyor("input1", 1)
	err := proc.AddInput(input1)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	input2 := conveyor.NewConveyor("input2", 1)
	err = proc.AddInput(input2)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	input3 := conveyor.NewConveyor("input3", 1)
	err = proc.AddInput(input3)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, 3, len(proc.inputs))
}

func TestAddInput_Fail_DuplicateName(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	input1 := conveyor.NewConveyor("input1", 1)
	err := proc.AddInput(input1)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	input1Dup := conveyor.NewConveyor("input1", 1)
	err = proc.AddInput(input1Dup)
	if !assert.Error(t, err) {
		assert.FailNow(t, "expected error")
	}
}

func TestRemoveInput_Succeed(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	input1 := conveyor.NewConveyor("input1", 1)
	proc.AddInput(input1)
	input2 := conveyor.NewConveyor("input2", 1)
	proc.AddInput(input2)
	input3 := conveyor.NewConveyor("input3", 1)
	proc.AddInput(input3)

	err := proc.RemoveInput("input2")
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
}

func TestRemoveInput_Fail_NotExist(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	input1 := conveyor.NewConveyor("input1", 1)
	proc.AddInput(input1)
	input2 := conveyor.NewConveyor("input2", 1)
	proc.AddInput(input2)
	input3 := conveyor.NewConveyor("input3", 1)
	proc.AddInput(input3)

	err := proc.RemoveInput("none")
	if !assert.Error(t, err) {
		assert.FailNow(t, "expected error")
	}
}

func TestAddOutput_Succeed(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	output1 := conveyor.NewConveyor("output1", 1)
	err := proc.AddOutput(output1)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	output2 := conveyor.NewConveyor("output2", 1)
	err = proc.AddOutput(output2)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	output3 := conveyor.NewConveyor("output3", 1)
	err = proc.AddOutput(output3)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, 3, len(proc.outputs))
}

func TestAddOutput_Fail_DuplicateName(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	output1 := conveyor.NewConveyor("output1", 1)
	err := proc.AddOutput(output1)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	output1Dup := conveyor.NewConveyor("output1", 1)
	err = proc.AddOutput(output1Dup)
	if !assert.Error(t, err) {
		assert.FailNow(t, "expected error")
	}
}

func TestRemoveOutput_Succeed(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	output1 := conveyor.NewConveyor("output1", 1)
	proc.AddOutput(output1)
	output2 := conveyor.NewConveyor("output2", 1)
	proc.AddOutput(output2)
	output3 := conveyor.NewConveyor("output3", 1)
	proc.AddOutput(output3)

	err := proc.RemoveOutput("output2")
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
}

func TestRemoveOutput_Fail_NotExist(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	output1 := conveyor.NewConveyor("output1", 1)
	proc.AddOutput(output1)
	output2 := conveyor.NewConveyor("output2", 1)
	proc.AddOutput(output2)
	output3 := conveyor.NewConveyor("output3", 1)
	proc.AddOutput(output3)

	err := proc.RemoveOutput("none")
	if !assert.Error(t, err) {
		assert.FailNow(t, "expected error")
	}
}

func TestStart_Succeed(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	isStarted := proc.IsStarted()
	if !assert.False(t, isStarted) {
		assert.FailNow(t, "expected not started")
	}

	err := proc.Start()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	isStarted = proc.IsStarted()
	assert.True(t, isStarted)
}

func TestStart_Fail_AlreadyStarted(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	isStarted := proc.IsStarted()
	if !assert.False(t, isStarted) {
		assert.FailNow(t, "expected not started")
	}

	err := proc.Start()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	isStarted = proc.IsStarted()
	if !assert.True(t, isStarted) {
		assert.FailNow(t, "expected started")
	}

	err = proc.Start()
	if !assert.Error(t, err) {
		assert.FailNow(t, "expected error")
	}
}

func TestStop_Succeed(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	err := proc.Start()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	err = proc.Stop()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	isStarted := proc.IsStarted()
	if !assert.False(t, isStarted) {
		assert.FailNow(t, "expected stopped")
	}
}

func TestStop_Fail_NotStarted(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	isStarted := proc.IsStarted()
	if !assert.False(t, isStarted) {
		assert.FailNow(t, "expected not started")
	}

	err := proc.Stop()
	if !assert.Error(t, err) {
		assert.FailNow(t, "expected error")
	}
}

func TestStop_Fail_AlreadyStopped(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	err := proc.Start()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	err = proc.Stop()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	isStarted := proc.IsStarted()
	if !assert.False(t, isStarted) {
		assert.FailNow(t, "expected stopped")
	}

	err = proc.Stop()
	if !assert.Error(t, err) {
		assert.FailNow(t, "expected error")
	}
}

func TestSend_Succeed(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	output1 := conveyor.NewConveyor("output1", 1)
	proc.AddOutput(output1)
	output2 := conveyor.NewConveyor("output2", 1)
	proc.AddOutput(output2)
	output3 := conveyor.NewConveyor("output3", 1)
	proc.AddOutput(output3)

	mt := metric.Metric{
		Name:   "test",
		Time:   time.Now(),
		Tags:   []metric.Tag{{Key: "tag1", Value: "value1"}},
		Fields: []metric.Field{{Key: "field1", Value: 1.0}},
	}
	err := proc.Send(mt)
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
}

func TestSend_Fail_NoOutput(t *testing.T) {
	proc := ProcessorBase{
		config: *cfg,
	}

	mt := metric.Metric{
		Name:   "test",
		Time:   time.Now(),
		Tags:   []metric.Tag{{Key: "tag1", Value: "value1"}},
		Fields: []metric.Field{{Key: "field1", Value: 1.0}},
	}
	err := proc.Send(mt)
	if !assert.Error(t, err) {
		assert.FailNow(t, "expected error")
	}
}
