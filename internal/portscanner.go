package internal

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ProcessInfo는 포트를 점유한 프로세스 정보를 담는다.
type ProcessInfo struct {
	PID     int
	Command string
}

// FindProcessByPort는 지정한 포트를 점유 중인 프로세스 정보를 반환한다.
// macOS 환경에서 lsof 사용
func FindProcessByPort(port int) (ProcessInfo, error) {
	cmd := exec.Command("lsof", "-i", ":"+strconv.Itoa(port), "-sTCP:LISTEN", "-n", "-P")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// lsof가 exit code 1이고 출력(표준/에러 모두)이 완전히 비었으면 "프로세스 없음" 상황으로 간주
		if _, ok := err.(*exec.ExitError); ok && strings.TrimSpace(string(out)) == "" {
			return ProcessInfo{}, nil
		}
		return ProcessInfo{}, fmt.Errorf("lsof 실행 실패: %w", err)
	}
	return ParseLsofOutput(string(out))
}

// ParseLsofOutput은 lsof 명령의 출력을 받아 ProcessInfo를 반환하는 순수함수다.
func ParseLsofOutput(output string) (ProcessInfo, error) {
	var proc ProcessInfo
	lines := strings.Split(output, "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[1]) == "" {
		// 프로세스 없음: 에러가 아니라 PID 0 반환
		return proc, nil
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return proc, fmt.Errorf("lsof 결과 파싱 실패: %v", lines[1])
	}
	pid, err := strconv.Atoi(fields[1])
	if err != nil {
		return proc, fmt.Errorf("PID 파싱 실패: %w", err)
	}
	proc.PID = pid
	proc.Command = fields[0]
	return proc, nil
}

