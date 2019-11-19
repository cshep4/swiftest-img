package score

import (
	"github.com/cshep4/swiftest-img/internal/img"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Scorer interface {
	Grade(text string) (*img.MarkedDocument, error)
}

type scorer struct {
}

func New() Scorer {
	return &scorer{}
}

func (s *scorer) Grade(text string) (*img.MarkedDocument, error) {
	text = strings.Replace(text, " ", "", -1)
	questions := strings.Split(text, "\n")

	var results []img.Question
	var score, total int
	for _, q := range questions {
		if s.isValidQuestion(q) {
			continue
		}

		parts := strings.Split(q, "=")

		first, operator, second, err := s.separateQuestion(parts[0])
		if err != nil {
			return nil, err
		}

		answer, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		expectedAnswer, err := s.calculateAnswer(first, operator, second)
		if err != nil {
			return nil, err
		}

		correct := answer == expectedAnswer
		if correct {
			score++
		}
		total++

		results = append(results, img.Question{
			Question: parts[0],
			Answer:   parts[1],
			Correct:  correct,
		})
	}

	return &img.MarkedDocument{
		Questions: results,
		Total:     total,
		Score:     score,
	}, nil
}

func (s *scorer) isValidQuestion(q string) bool {
	return q == "" || !strings.Contains(q, "=") || !s.containsOperator(q)
}

func (s *scorer) containsOperator(question string) bool {
	for k, _ := range operators {
		if strings.Contains(question, k) {
			return true
		}
	}

	return false
}

func (s *scorer) separateQuestion(question string) (string, string, string, error) {
	for k, v := range operators {
		if strings.Contains(question, k) {
			numbers := strings.Split(question, k)
			return numbers[0], v, numbers[1], nil
		}
	}

	return "", "", "", ErrUnsupportedOperator
}

func (s *scorer) calculateAnswer(first, operator, second string) (int, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return 0, err
	}

	first = reg.ReplaceAllString(first, "")
	firstInt, err := strconv.Atoi(first)
	if err != nil {
		return 0, err
	}

	second = reg.ReplaceAllString(second, "")
	secondInt, err := strconv.Atoi(second)
	if err != nil {
		return 0, err
	}

	switch operator {
	case "+":
		return firstInt + secondInt, nil
	case "-":
		return firstInt - secondInt, nil
	case "x":
		return firstInt * secondInt, nil
	case "/":
		return firstInt / secondInt, nil
	case "^":
		return int(math.Pow(float64(firstInt), float64(secondInt))), nil
	}

	return 0, ErrUnsupportedOperator
}
