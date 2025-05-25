package internal

import (
	"strings"
	"testing"
)

func TestParseConfirmationInput(t *testing.T) {
	cases := map[string]bool{
		"y": true,
		"Y": true,
		"yes": true,
		"YES": true,
		"n": false,
		"no": false,
		"": false,
		"  y  ": true,
		"  no  ": false,
	}
	for input, want := range cases {
		got := ParseConfirmationInput(input)
		if got != want {
			t.Errorf("입력 파싱 실패: input=%q got=%v want=%v", input, got, want)
		}
	}
}

func TestAskForConfirmationWithIO(t *testing.T) {
	// autoYes true면 무조건 true
	ok := AskForConfirmationWithIO(strings.NewReader("n\n"))
	if ok {
		t.Error("입력이 n인데 true 반환")
	}
	// 입력 y
	ok = AskForConfirmationWithIO(strings.NewReader("y\n"))
	if !ok {
		t.Error("입력이 y인데 false 반환")
	}
	// 잘못된 입력 후 y (여러 입력 중 마지막이 y)
	ok = AskForConfirmationWithIO(strings.NewReader("foo\ny\n"))
	if !ok {
		t.Error("여러 입력 중 y가 마지막인데 false 반환")
	}
}
