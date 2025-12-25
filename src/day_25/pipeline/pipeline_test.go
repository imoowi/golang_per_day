package pipeline

import "testing"

func TestPipeline(t *testing.T) {
	//生成打牌动作
	plays := genPlays()
	validatedPlays := stageValidate(plays)
	legalizedPlays := stageLegalize(validatedPlays)
	stageBroadcast(legalizedPlays)
}
