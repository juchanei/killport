package main

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"strconv"

	"github.com/user/port-killer/internal"
)

func main() {
	var autoYes bool
	var verbose bool

	var rootCmd = &cobra.Command{
		Use:   "killport [port]",
		Short: "특정 포트를 점유한 프로세스를 종료하는 CLI 도구",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			portStr := args[0]
			port, err := strconv.Atoi(portStr)
			if err != nil || port < 1 || port > 65535 {
				fmt.Fprintf(os.Stderr, "[ERROR] 올바른 포트 번호를 입력하세요: %s\n", portStr)
				os.Exit(1)
			}
			proc, err := internal.FindProcessByPort(port)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] 포트 %d를 점유한 프로세스 탐색 실패: %v\n", port, err)
				os.Exit(1)
			}
			if proc.PID == 0 {
				fmt.Fprintf(os.Stdout, "[INFO] 포트 %d를 점유한 프로세스가 없습니다.\n", port)
				return
			}
			msg := fmt.Sprintf("포트 %d (PID %d, CMD: %s)를 종료하시겠습니까?", port, proc.PID, proc.Command)
			if autoYes {
				fmt.Fprintf(os.Stdout, "%s [자동 승인]\n", msg)
			} else {
				for {
					fmt.Fprintf(os.Stdout, "%s [y/N]: ", msg)
					if internal.AskForConfirmationWithIO(os.Stdin) {
						break
					} else {
						fmt.Fprintln(os.Stdout, "y 또는 n으로 입력하세요.")
					}
				}
			}
			err = internal.KillProcessByPID(proc.PID, os.Stdout, verbose)
			if err != nil {
				fmt.Fprintf(os.Stdout, "[ERROR] 프로세스 종료 실패: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "[SUCCESS] PID %d 종료 완료\n", proc.PID)

		},
	}
	rootCmd.Flags().BoolVarP(&autoYes, "yes", "y", false, "사용자 프롬프트 없이 자동 승인")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "상세 로그 출력")
	rootCmd.Execute()
}

