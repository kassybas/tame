package stepblock

import "github.com/kassybas/tame/internal/step"

type StepBlock struct {
	steps []step.Step
}

func NewStepBlock(sb []step.Step) StepBlock {
	var newStepBlock StepBlock
	newStepBlock.steps = make([]step.Step, len(sb))
	for i := range sb {
		newStepBlock.steps[i] = sb[i]
	}
	return newStepBlock
}

func (sb *StepBlock) Get(index int) *step.Step {
	return &sb.steps[index]
}
func (sb *StepBlock) GetAll() []step.Step {
	return sb.steps
}

func (sb *StepBlock) Set(index int, s step.Step) {
	sb.steps[index] = s
}
