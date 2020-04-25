# 인증서 만드는 방법

### Letsencrypt 방법

#### macOS

macOS에 letsencrypt 를 설치합니다.

```bash
$ brew install letsencrypt
```

인증서를 생성합니다. -d 값에는 도메인을 넣어주세요.

```bash
$ sudo certbot certonly --standalone -d dotori.lazypic.org
```

> 참고: `-d lazypic.org -d www.lazypic.org` 형태로 -d 옵션을 이용해서 복수로 도메인을 등록할 수 있습니다.

인증서는 다음 경로에 생성됩니다.

```
/etc/letsencrypt/live/dotori.lazypic.org/fullchain.pem
/etc/letsencrypt/live/dotori.lazypic.org/privkey.pem
```

인증서 갱신

```bash
$ sudo certbot renew  --dry-run
```

Crontab 을 이용한 자동 인증서 갱신

```bash
$ crontab -e
* */12 * * * root /usr/bin/certbot renew >/dev/null 2>&1
```


#### CentOS7

certbot을 설치합니다.

```bash
$ sudo yum install epel-release
$ yum install certbot mod_ssl
```

인증서를 생성합니다. -d 값에는 도메인이 넣어주세요.
```bash
$ sudo certbot certonly --standalone -d dotori.lazypic.org
```

참고자료: https://www.rosehosting.com/blog/how-to-install-lets-encrypt-with-apache-on-centos-7/

### 자가 인증서 생성방법
스스로 사인한 인증서를 생성하면 접속시 에러가 발생하지만, https 보안프로토콜을 사용할 수 있습니다.
go를 설치하면 src/crypto/tls/generate_cert.go 파일을 이용해서 인증서를 생성할 수 있습니다.

```bash
$ go run /usr/local/go/src/crypto/tls/generate_cert.go -host="dotori.lazypic.org" -ca=true
```