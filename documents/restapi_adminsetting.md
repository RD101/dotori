# RestAPI - Admin Setting

## Get

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/adminsetting | Admin Setting 가지고 오기 | . | `$ curl -X Get "https://dotori.lazypic.com/api/adminsetting"` |

## Post

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/rename | 파일명 변경 | path, find, replace, permission | `$ curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"path":"/dotori/62/5e/1d/8e/1f9107/ad8e/ad/a0/17/data/","find":"A00", "replace":"W00", ":permission":false}'  "https://dotori.lazypic.com/api/rename"` |
| /api/dbbackup | dbbackup | date | `$ curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"date":"20221111"}'  "https://dotori.lazypic.com/api/dbbackup"` |