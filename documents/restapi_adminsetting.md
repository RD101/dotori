# RestAPI - Admin Setting

## Get

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/adminsetting | Admin Setting 가지고 오기 | . | `$ curl -X Get "https://dotori.lazypic.com/api/adminsetting"` |

## Post

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/rename | 파일명 변경 | path, find, replace, permission | `$ curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"path":"/asset/data/path","find":"A00", "replace":"W00", ":permission":false}'  "http://172.30.1.20/api/rename"` |