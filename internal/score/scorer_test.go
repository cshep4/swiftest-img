package score

import (
	"github.com/cshep4/swiftest-img/internal/img"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScorer_Grade(t *testing.T) {
	scorer := New()

	text := "1 + 4 = 5\n2 + 3 = 5\n4 - 8 = -10\n5 x 6 = 30\n"

	expectedResult := &img.MarkedDocument{
		Questions: []img.Question{
			{
				Question: "1+4",
				Answer:   "5",
				Correct:  true,
			},
			{
				Question: "2+3",
				Answer:   "5",
				Correct:  true,
			},
			{
				Question: "4-8",
				Answer:   "-10",
				Correct:  false,
			},
			{
				Question: "5x6",
				Answer:   "30",
				Correct:  true,
			},
		},
		Total: 4,
		Score: 3,
	}

	res, err := scorer.Grade(text)
	require.NoError(t, err)

	assert.Equal(t, expectedResult, res)
}
