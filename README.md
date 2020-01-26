# Dotori

![travisCI](https://secure.travis-ci.org/rd101/dotori.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/rd101/dotori)](https://goreportcard.com/report/github.com/rd101/dotori)

Dotori is web based asset library tool.

## 목 표
- VFX, 애니메이션, 게임에 사용되는 모든 타입의 에셋을 등록,관리
- 사용라이브러리: OpenColorIO, OpenImageIO, FFmpeg
- 주의사항: 회사 특이사항과 관련된 코드를 내부에 넣지말것. 셋팅영역으로 뺄 것
- 테스트서버: http://csi.lazypic.org:8089

## DB 설치, 실행
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

## Dotori 실행

```bash
$ sudo dotori -http :80
```

> 여러분이 macOS를 사용한다면 기본적으로 80포트는 아파치 서버가 사용중일 수 있습니다. `:80` 포트에 실행되는 아파치 서버를 종료하기 위해서 $ sudo apachectl stop 를 터미널에 입력해주세요.

## Download
- [Linux](https://github.com/RD101/dotori/releases/download/v0.0.1/dotori_linux_x86-64.tgz) 
- [macOS](https://github.com/RD101/dotori/releases/download/v0.0.1/dotori_darwin_x86-64.tgz)
- [Windows](https://github.com/RD101/dotori/releases/download/v0.0.1/dotori_windows_x86-64.tgz)

## Command-line

#### Item 추가
```bash
$ dotori -add -inputpath /project/path -outputpath /library/backup/path -author woong -tag asset,3D -description 설명 -type maya
```

#### Item 삭제
```bash
$ sudo dotori -rm -id "idstring"
```

#### 웹서버 실행
```bash
$ sudo dotori -http :80
```

## restAPI
RestAPI 생성후 기록 예정

## Infomation / History
- '19.9 RD101에서 오픈소스로 시작
- License: BSD 3-Clause License
