# Dotori

[![Go Report Card](https://goreportcard.com/badge/github.com/rd101/dotori)](https://goreportcard.com/report/github.com/rd101/dotori)

Dotori is web-based asset library tool.

### 목 표
- VFX, 애니메이션, 게임, 사운드 작업에 사용되는 모든 타입의 에셋을 등록,관리
- 사용라이브러리: OpenColorIO, OpenImageIO, FFmpeg
- 주의사항: 회사 특이사항과 관련된 코드를 내부에 넣지말것. 셋팅영역으로 뺄 것
- 테스트서버: https://csi.lazypic.org:8089

#### 서버권장사항
- 동시접속자 처리를 위한 OS: Linux, macOS, Windows Server
- 메모리 32기가 이상. (데이터를 많이 처리할 때 DB는 약 14기가의 메모리를 사용한다.)

### DB 설치, 실행
mongoDB 설치

CentOS

```bash
$ sudo yum install mongodb mongodb-server
$ sudo service mongod start
```

macOS

```bash
$ brew uninstall mongodb
$ brew tap mongodb/brew
$ brew install mongodb-community
$ brew services start mongodb-community
```

### 기타 library 및 명령어 설치
- [Library 설치](documents/setlibrary.md)

### Dotori 실행

```bash
$ sudo dotori -http :80
```

> 여러분이 macOS를 사용한다면 기본적으로 80포트는 아파치 서버가 사용중일 수 있습니다. `:80` 포트에 실행되는 아파치 서버를 종료하기 위해서 $ sudo apachectl stop 를 터미널에 입력해주세요.

### Download
- [Linux](https://github.com/RD101/dotori/releases/download/v0.0.1/dotori_linux_x86-64.tgz) 
- [macOS](https://github.com/RD101/dotori/releases/download/v0.0.1/dotori_darwin_x86-64.tgz)
- [Windows](https://github.com/RD101/dotori/releases/download/v0.0.1/dotori_windows_x86-64.tgz)

### Command-line

#### Item 추가
```bash
$ dotori -add -inputpath /project/path -author woong -tag asset,3D -description 설명 -type maya
```

#### Item 삭제
```bash
$ sudo dotori -remove -itemtype maya -itemid 5e89cef7cd1747fd5eacf256
```

#### 웹서버 실행
```bash
$ sudo dotori -http :80
```

- [인증서 만드는 방법](documents/how_to_make_certification.md)

### REST API
Dotori는 REST API를 지원합니다. Python, Go, Java, Javascript, node.JS, C++, C, C# 등 수많은 언어를 통해 Dotori를 이용할 수 있습니다.
아래는 Dotori restAPI reference 문서입니다.
- [item](documents/restapi_item.md)
- [admin setting](documents/restapi_adminsetting.md)

### 개발환경셋팅
Go에서 컴파일된 파일이 생성되는 경로를 설정하기 위해 GOBIN 환경변수 셋팅이 필요합니다.

리눅스라면 .bashrc에 선언해주세요.
macOS이고 zsh쉘을 사용한다면 `.zshenv` 에 bash쉘을 사용한다면 `.bashrc`에 아래 설정을 추가해주세요.

```bash
export GOBIN=$HOME/bin
export PATH=$PATH:$GOBIN
```

dotori는 sudo로 실행해야 합니다. 그러나 linux의 경우, sudo가 현재 계정의 PATH를 다 가져오지 못하는 경우가 있습니다. 그럴 때는 /etc/visudoers 파일을 아래처럼 변경해주세요.

```bash
$ sudo visudo
...

#Default secure_path="/usr/local/sbin:/usr/local/bin:/usr/bin" # 기존 부분 주석 처리
Default env_keep=PATH # 새로 추가
```

### 예제파일
- 에셋 라이브러리 개발에 사용된 예제 파일은 `examples` 폴더에 들어있습니다.
- footage 데이터
    - footage 데이터는 95메가 정도의 용량을 가지고 있습니다.
    - 리포지터리에는 최대한 가벼운 파일, 코드만 올리기 위해 위 폴더에 footage 데이터는 들어가 있지 않습니다.
    - footage 예제파일은 https://github.com/lazypic/tdcourse_examples/tree/master/footage 에서 다운받을 수 있습니다.

### Infomation / History
- '19.9 RD101에서 오픈소스로 시작
- License: BSD 3-Clause License
