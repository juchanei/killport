package internal

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// 테스트용으로 sleep 프로세스를 띄우고 종료 테스트
func TestKillProcessByPID(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("테스트용 프로세스 실행 실패: %v", err)
	}
	pid := cmd.Process.Pid
	t.Logf("테스트용 PID: %d", pid)

	var buf bytes.Buffer
	err = KillProcessByPID(pid, &buf, false)
	if err != nil {
		t.Errorf("KillProcessByPID 실패: %v", err)
	}

	// 종료된 프로세스의 리소스 정리 (좀비 방지)
	_ = cmd.Wait()

	// 프로세스가 실제로 종료됐는지 확인
	out, _ := exec.Command("ps", "-p", strconv.Itoa(pid)).Output()
	if strings.Contains(string(out), strconv.Itoa(pid)) {
		t.Errorf("프로세스가 종료되지 않음: PID %d", pid)
	}
}
