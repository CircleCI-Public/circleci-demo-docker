package canarygo

import "testing"

func TestGutterBalls(t *testing.T) {
	t.Log("Rolling all gutter balls... (expected score: 0)")
	game := NewGame()
	game.rollMany(20, 0)

	if score := game.Score(); score != 0 {
		t.Errorf("Expected score of 0, but it was %d instead.", score)
	}
}

func TestOnePinOnEveryThrow(t *testing.T) {
	t.Log("Each throw knocks down one pin... (expected score: 20)")
	game := NewGame()
	game.rollMany(20, 1)

	if score := game.Score(); score != 20 {
		t.Errorf("Expected score of 20, but it was %d instead.", score)
	}
}

func TestSingleSpare(t *testing.T) {
	t.Log("Rolling a spare, then a 3, then all gutters... (expected score: 16)")
	game := NewGame()
	game.rollSpare()
	game.Roll(3)
	game.rollMany(17, 0)

	if score := game.Score(); score != 16 {
		t.Errorf("Expected score of 16, but it was %d instead.", score)
	}
}

func TestSingleStrike(t *testing.T) {
	t.Log("Rolling a strike, then 3, then 7, then all gutters... (expected score: 24)")
	game := NewGame()
	game.rollStrike()
	game.Roll(3)
	game.Roll(4)
	game.rollMany(16, 0)

	if score := game.Score(); score != 24 {
		t.Errorf("Expected score of 24, but it was %d instead.", score)
	}
}

func TestPerfectGame(t *testing.T) {
	t.Log("Rolling all strikes... (expected score: 300)")
	game := NewGame()
	game.rollMany(21, 10)

	if score := game.Score(); score != 300 {
		t.Errorf("Expected score of 300, but it was %d instead.", score)
	}
}