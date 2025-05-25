package internal

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestParseLsofOutput(t *testing.T) {
	output := "COMMAND   PID USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME\nPython 12345 user   3u  IPv4 0x123456789abcdef0      0t0  TCP *:54321 (LISTEN)"
	proc, err := ParseLsofOutput(output)
	if err != nil {
		t.Fatalf("예상치 못한 에러: %v", err)
	}
	if proc.PID != 12345 {
		t.Errorf("PID 파싱 실패: got %d, want 12345", proc.PID)
	}
	if proc.Command != "Python" {
		t.Errorf("Command 파싱 실패: got %s, want Python", proc.Command)
	}

	// 여러 커맨드/포맷도 테스트 (예: node)
	output2 := "COMMAND   PID USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME\nnode 9999 user   3u  IPv4 0xabcdef      0t0  TCP *:3000 (LISTEN)"
	proc2, err2 := ParseLsofOutput(output2)
	if err2 != nil {
		t.Fatalf("예상치 못한 에러: %v", err2)
	}
	if proc2.PID != 9999 {
		t.Errorf("PID 파싱 실패: got %d, want 9999", proc2.PID)
	}
	if proc2.Command != "node" {
		t.Errorf("Command 파싱 실패: got %s, want node", proc2.Command)
	}

	// PID가 숫자가 아닌 경우 등은 예외 테스트에서 다루므로 여기선 생략
}

// 테스트용 임시 서버를 띄우고, 해당 포트를 점유한 프로세스 정보를 올바르게 찾는지 검증합니다.
func TestFindProcessByPort(t *testing.T) {
	// 테스트에 사용할 포트
	port := 54321

	// 임시로 Python HTTP 서버를 실행 (테스트 환경에 Python이 있다고 가정)
	cmd := exec.Command("python3", "-m", "http.server", strconv.Itoa(port))
	if err := cmd.Start(); err != nil {
		t.Fatalf("임시 서버 실행 실패: %v", err)
	}
	defer cmd.Process.Kill()

	// 잠시 대기 (서버가 완전히 올라올 때까지)
	tempWait()

	// 실제 함수 호출
	proc, err := FindProcessByPort(port)
	if err != nil {
		t.Fatalf("프로세스 탐색 실패: %v", err)
	}
	if proc.PID == 0 {
		t.Error("PID가 0입니다.")
	}
	if !strings.Contains(strings.ToLower(proc.Command), "python") {
		t.Errorf("명령어가 예상과 다름: %s", proc.Command)
	}
}

// 간단한 대기 함수 (1초)
func tempWait() {
	cmd := exec.Command("sleep", "1")
	_ = cmd.Run()
}
