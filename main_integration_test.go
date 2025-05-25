package main

import (
	"net"
	"os/exec"
	"strings"
	"testing"
)

// 테스트용 미사용 포트 번호
const port = "54321"
const binaryPath = "./bin/port-killer" // 빌드된 바이너리 경로

// 종료할 프로세스가 없는 포트 입력 시나리오
func TestKillPort_NoProcess(t *testing.T) {
	cmd := exec.Command(binaryPath, port)
	output, err := cmd.CombinedOutput()

	outStr := string(output)

	// 종료 코드가 0(성공)인지 확인
	if exitErr, ok := err.(*exec.ExitError); ok {
		t.Fatalf("프로세스가 없을 때 종료 코드가 0이 아님. exit code: %d, 출력: %s", exitErr.ExitCode(), outStr)
	} else if err != nil {
		t.Fatalf("실행 에러: %v, 출력: %s", err, outStr)
	}

	if !strings.Contains(outStr, "를 점유한 프로세스가 없습니다") {
		t.Errorf("'를 점유한 프로세스가 없습니다' 메시지가 출력되지 않음. 실제 출력: %s", outStr)
	}
}

// 잘못된 포트 입력 시나리오
func TestKillPort_InvalidPort(t *testing.T) {
	invalidPorts := []string{"abc", "70000", "0", "65536"}
	for _, p := range invalidPorts {
		t.Run(p, func(t *testing.T) {
			cmd := exec.Command(binaryPath, p)
			output, err := cmd.CombinedOutput()

			outStr := string(output)
			// 반드시 비정상 종료여야 함
			if err == nil {
				t.Fatalf("잘못된 포트 입력 시 비정상 종료가 아님. 입력: %s, 출력: %s", p, outStr)
			}
			// exit code가 0이 아니어야 함
			if exitErr, ok := err.(*exec.ExitError); ok {
				if exitErr.ExitCode() == 0 {
					t.Errorf("잘못된 포트 입력 시 exit code가 0임. 입력: %s, 출력: %s", p, outStr)
				}
			} else {
				t.Errorf("ExitError 타입이 아님. 입력: %s, err: %v, 출력: %s", p, err, outStr)
			}
			// 에러 메시지 확인 (실제 출력에 맞게 수정)
			if !strings.Contains(outStr, "올바른 포트 번호를 입력하세요") {
				t.Errorf("적절한 에러 메시지가 출력되지 않음. 입력: %s, 실제 출력: %s", p, outStr)
			}
		})
	}
}

// 정상 종료 시나리오: macOS 기본 nc(netcat)로 임시 서버 실행
func TestKillPort_NormalWithNc(t *testing.T) {

	// 1. nc로 임시 서버를 백그라운드에서 실행
	ncCmd := exec.Command("nc", "-l", port)
	ncCmd.Stdout = nil
	ncCmd.Stderr = nil

	if err := ncCmd.Start(); err != nil {
		t.Fatalf("nc 임시 서버 실행 실패: %v", err)
	}

	defer func() {
		_ = ncCmd.Process.Kill()
		_ = ncCmd.Wait()
	}()

	// 2. nc가 실제로 포트를 listening할 때까지 대기 (0.5초)
	_ = exec.Command("sleep", "0.5").Run()

	out, _ := exec.Command("lsof", "-i", ":"+port).CombinedOutput()
	if !strings.Contains(string(out), "LISTEN") {
		t.Fatalf("nc가 %s 포트에서 정상적으로 listen하지 않음", port)
	}

	// 3. port-killer로 해당 포트를 종료

	cmd := exec.Command(binaryPath, port)
	stdin, errPipe := cmd.StdinPipe()
	if errPipe != nil {
		t.Fatalf("port-killer StdinPipe 생성 실패: %v", errPipe)
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte("y\n"))
	}()

	output, err := cmd.CombinedOutput()

	// 4. 종료 코드와 출력 검증
	if err != nil {
		t.Fatalf("port-killer 실행 실패: %v, 출력: %s", err, string(output))
	}
	if !strings.Contains(string(output), "[SUCCESS]") {
		t.Errorf("정상 종료 메시지 누락. 실제 출력: %s", string(output))
	}

	// 5. nc 프로세스가 종료되었는지 확인
	conn, dialErr := net.Dial("tcp", "127.0.0.1:"+port)
	if dialErr == nil {
		_ = conn.Close()
	}

	_ = exec.Command("sleep", "0.1").Run()

	if err := ncCmd.Process.Signal(nil); err == nil {
		t.Errorf("nc 임시 서버가 종료되지 않음 (아직 살아있음)")
	}
}
