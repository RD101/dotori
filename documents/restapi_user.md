# RestAPI - User

## Put

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/user/autoplay | autoplay 설정 | value | `$ curl -k -X PUT -H "Authorization: Basic {TOKEN}" "https://dotori.lazypic.com/api/user/autoplay?value=true"` |
| /api/user/newsnum | news 갯수 설정 | value | `$ curl -k -X PUT -H "Authorization: Basic {TOKEN}" "https://dotori.lazypic.com/api/user/newsnum?value=4"` |
| /api/user/topnum | top 갯수 설정 | value | `$ curl -k -X PUT -H "Authorization: Basic {TOKEN}" "https://dotori.lazypic.com/api/user/topnum?value=4"` |
| /api/user/accesslevel | AccessLevel 설정 | json | `$ curl -k -X POST -H "Authorization: Basic {TOKEN}" -d '{"id":"userid","accesslevel":"admin"}' "https://dotori.lazypic.com/api/user/accesslevel"` |
