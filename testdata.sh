#!/bin/sh
#DB에 아이템 추가 테스트를 위한 명령어를 모아둔 파일
dotori -add -itemtype maya -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1"
dotori -add -itemtype maya -title example -author woong -description "description2 about some details" -tag "나무 낙엽 item2"
dotori -add -itemtype maya -title example -author woong -description "description3 about some details" -tag "나무 낙엽 item3"
dotori -add -itemtype maya -title example -author woong -description "description4 about some details" -tag "자동차 트럭"
dotori -add -itemtype maya -title example -author bailey -description "description5 about some details" -tag "자동차 SUV"
dotori -add -itemtype maya -title example -author bailey -description "description6 about some details" -tag "하늘 노을"
dotori -add -itemtype maya -title example -author bailey -description "description7 about some details" -tag "하늘 밤"
dotori -add -itemtype maya -title example -author bailey -description "description8 about some details" -tag "하늘 아침"
dotori -add -itemtype maya -title example -author bailey -description "description9 about some details" -tag "기차 KTX"
dotori -add -itemtype maya -title example -author bailey -description "description10 about some details" -tag "기차 무궁화호"
dotori -add -itemtype usd -title example -author bailey -description "description10 about some details" -tag "기차 무궁화호"
dotori -add -itemtype sound -title example -author bailey -description "description10 about some details" -tag "bgm"

#REST API
curl -X POST -d "author=bae&itemtype=maya" http://127.0.0.1/api/item
#DELETE는 실행 전에 직접 id를 넣어준 후 테스트한다.
#curl -X DELETE "http://127.0.0.1/api/item?type=nuke&id="