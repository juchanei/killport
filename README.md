# killport

`killport`는 macOS 환경에서 특정 포트를 점유한 프로세스를 빠르게 종료할 수 있도록 도와주는 간단한 CLI 도구입니다.

개발 중 IDE나 터미널이 비정상 종료되어 포트가 해제되지 않을 경우, 서버를 재실행하려고 하면 충돌이 발생할 수 있습니다. `killport`는 이런 상황을 해결하기 위해 만들어졌습니다.

## 사전 준비

lsof 명령어가 설치되어 있어야 합니다.

## 설치

Go 환경이 구성된 경우:

```bash
go install github.com/juchanei/killport@latest
```

설치 후 생성된 바이너리가 `$HOME/go/bin`(Go 1.17 이상) 또는 `$GOPATH/bin`에 위치합니다. 해당 디렉터리가 PATH에 포함되어 있어야 터미널에서 `killport` 명령어를 사용할 수 있습니다.

예시:

```bash
export PATH="$PATH:$HOME/go/bin"
```

위 명령어를 `~/.zshrc` 또는 `~/.bash_profile` 등에 추가해두면, 터미널을 새로 열 때마다 자동으로 적용됩니다.

## 사용법

```bash
killport 3000             # 포트 3000을 점유한 프로세스를 종료할지 물음
killport 3000 -y          # 포트 3000을 점유한 프로세스를 묻지 않고 종료
killport 3000 --verbose   # 로그를 출력하며 실행
killport 3000 -v
```

- `-y` 옵션으로 사용자 확인 생략(자동 승인)
- `-v`, `--verbose` 옵션으로 로그 출력

## 주의사항
- 시스템/다른 사용자의 프로세스 종료 시 root 권한이 필요할 수 있습니다.
- 종료 후 리소스 정리를 위해 wait() 호출이 필요할 수 있습니다.
- 현재는 macOS만 지원됩니다.
