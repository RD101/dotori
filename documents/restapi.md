# RestAPI

## Get
| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/item | 아이템 가지고 오기 | itemtype, id | `$ curl "http://192.168.219.104/api/item?itemtype=maya&id=5e24742f901da0498519f7a7"` |


## Post
| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/search | 검색하기 | itemtype, searchword | `$ curl -X POST -d "itemtype=maya&searchword=나무" http://192.168.219.104/api/search` |


## Python example
### asset 가지고 오기 

```python
#!/usr/bin/python
#coding:utf-8
import urllib2
import json
try:
    request = urllib2.Request("http://192.168.219.104/api/item?itemtype=maya&id=5e24742f901da0498519f7a7")
    result = urllib2.urlopen(request)
    data = json.load(result)
except:
    print("RestAPI에 연결할 수 없습니다.")
    # 이후 에러처리 할 것
print(data)
```

### 검색
```python
#!/usr/bin/python
#coding:utf-8
import urllib2
import json
try:
    request = urllib2.Request("http://192.168.219.104/api/search?itemtype=maya&searchword=나무")
    result = urllib2.urlopen(request)
    data = json.load(result)
except:
    print("RestAPI에 연결할 수 없습니다.")
print(data)
```