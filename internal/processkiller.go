package internal

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// KillProcessByPID는 주어진 PID의 프로세스를 종료한다.
// 우선 SIGTERM을 보내고, 실패 시 SIGKILL을 시도한다.
// 종료 결과와 에러를 반환한다.
// KillProcessByPID는 주어진 PID의 프로세스를 종료한다.
// 로그 메시지는 []string으로 반환하고, 에러는 error로 반환한다.
import "io"

// KillProcessByPID는 주어진 PID의 프로세스를 종료한다.
// 로그는 writer로 바로 출력하며, 에러만 반환한다.
func KillProcessByPID(pid int, w io.Writer, verbose bool) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintf(w, "[ERROR] PID %d: 프로세스 조회 실패: %v\n", pid, err)
		return fmt.Errorf("[ERROR] PID %d: 프로세스 조회 실패: %w", pid, err)
	}

	// SIGTERM 시도
	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		if verbose {
			fmt.Fprintf(w, "[WARN] PID %d: SIGTERM 실패: %v\n", pid, err)
		}
	} else {
		if verbose {
			fmt.Fprintf(w, "[INFO] PID %d: SIGTERM 전송됨\n", pid)
		}
	}

	// 프로세스가 종료됐는지 확인 (간단히 0.5초 대기)
	time.Sleep(500 * time.Millisecond) // 0.5초

	// 종료되지 않았으면 SIGKILL
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		if verbose {
			fmt.Fprintf(w, "[WARN] PID %d: 아직 살아있음, SIGKILL 시도\n", pid)
		}
		err = proc.Signal(syscall.SIGKILL)
		if err != nil {
			fmt.Fprintf(w, "[ERROR] PID %d: SIGKILL 실패: %v\n", pid, err)
			return fmt.Errorf("[ERROR] PID %d: SIGKILL 실패: %w", pid, err)
		}
		if verbose {
			fmt.Fprintf(w, "[INFO] PID %d: SIGKILL 전송됨\n", pid)
		}
	} else {
		if verbose {
			fmt.Fprintf(w, "[INFO] PID %d: 정상적으로 종료됨\n", pid)
		}
	}
	return nil
}
