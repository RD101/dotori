# RestAPI Category

Category Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/category/{id} | Category 정보를 가져옵니다 | id | curl -X GET -H "Authorization: Basic {TOKEN}" "https://dotori.lazypic.com/api/category/{id}"
| /api/rootcategories | 메인 Category 정보를 가져옵니다 | . | curl -X GET -H "Authorization: Basic {TOKEN}" "https://dotori.lazypic.com/api/rootcategories"
| /api/subcategories/{parentname} | 서브 Category 정보를 가져옵니다 | parentname | curl -X GET -H "Authorization: Basic {TOKEN}" "https://dotori.lazypic.com/api/subcategories/{parentname}"

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/category/{id} | 기존 Category 정보를 수정합니다 | name, parentname |curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"name":"env","parentname":""}' "https://dotori.lazypic.com/api/category/{id}"

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/category | 카테고리를 생성합니다. | name, parentname | curl -X POST -H "Authorization: Basic {TOKEN}" -d '{"name":"env","parentname":""}' "https://dotori.lazypic.com/api/category"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/category/{id} | 기존 Category 정보를 삭제합니다 | id |curl -X DELETE -H "Authorization: Basic {TOKEN}“ "https://dotori.lazypic.com/api/category/{id}"

## Option 체크

```bash
curl https://dotori.lazypic.com/api/category -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,PUT,DELETE,POST,OPTIONS
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
