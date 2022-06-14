# RestAPI Tags

Tags Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/tags/{id}| tags 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://dotori.lazypic.com/api/tags/{id}"

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/tags/{id}|기존 tags 정보를 수정합니다|tags|curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"tags":["tag1","tag2","tag3"]}' "https://dotori.lazypic.com/api/tags/{id}"

## Option 체크

```bash
curl https://dotori.lazypic.com/api/tags -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,PUT,OPTIONS
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
