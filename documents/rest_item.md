# RestAPI Item

## Python 예제

### 에셋 검색하기
```python
#coding:utf8
import json
import urllib2

restURL = "http://172.16.101.230/search?itemtype=maya&searchword=나무"
try:
    data = json.load(urllib2.urlopen(restURL))
except:
    print("RestAPI에 연결할 수 없습니다.")
    #에러처리
if "error" in data:
    print(data["error"])
print(data["data"])
```
