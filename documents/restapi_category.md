# RestAPI Tags

Tags Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/tags/{id}| tags 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://dotori.lazypic.com/api/tags/{id}"

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/category/{id} | 기존 Category 정보를 수정합니다 | name, parentname |curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"name":"env","parentname":""}' "https://dotori.lazypic.com/api/category/{id}"


## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/category | 카테고리를 생성합니다. | name, parentname | curl -X POST -H "Authorization: Basic {TOKEN}" -d '{"name":"env","parentname":""}' "https://dotori.lazypic.com/api/category"


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
