package internal

import (
	"bufio"
	"io"
	"strings"
)


// ParseConfirmationInput은 입력값을 받아 종료 여부를 bool로 반환하는 순수함수입니다.
func ParseConfirmationInput(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

// AskForConfirmationWithIO는 io.Reader/Writer를 활용해 프롬프트 및 입력 확인을 수행합니다.
func AskForConfirmationWithIO(in io.Reader) bool {
	reader := bufio.NewReader(in)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		if ParseConfirmationInput(input) {
			return true
		} else if input == "n" || input == "no" || strings.TrimSpace(input) == "" {
			return false
		}
		// 잘못된 입력이면 반복 (메시지 출력은 main.go에서 담당)
	}
}

