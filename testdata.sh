#!/bin/sh
#DB에 아이템 추가 테스트를 위한 명령어를 모아둔 파일
#copyFile함수 테스트용 데이터
dotori -add -author bae -itemtype maya -inputpath /tmp/dotori/ -outputpath /Users/seoyoungbae/dotoriDB/ -description "copying folder test1" -tag "copy test"
dotori -add -author bae -itemtype maya -inputpath /tmp/dotori -outputpath /Users/seoyoungbae/dotoriDB -description "copying folder test2" -tag "copy test"
dotori -add -author bae -itemtype nuke -inputpath /tmp/dotori -outputpath /Users/seoyoungbae/dotoriDB -description "copying folder test3" -tag "copy test"
#
dotori -add -author bae -itemtype nuke -inputpath /library/asset -outputpath /library/asset
dotori -add -author liah -itemtype houdini -inputpath /library/asset -outputpath /library/asset
dotori -add -author woong -itemtype maya -inputpath /library/asset/01 -outputpath /library/asset -description "description1 about some details" -tag "나무 낙엽 item1"
dotori -add -author woong -itemtype maya -inputpath /library/asset/02 -outputpath /library/asset -description "description2 about some details" -tag "나무 낙엽 item2"
dotori -add -author woong -itemtype maya -inputpath /library/asset/03 -outputpath /library/asset -description "description3 about some details" -tag "나무 낙엽 item3"
dotori -add -author woong -itemtype maya -inputpath /library/asset/04 -outputpath /library/asset -description "description4 about some details" -tag "나무 낙엽 item4"
dotori -add -author woong -itemtype maya -inputpath /library/asset/05 -outputpath /library/asset -description "description5 about some details" -tag "나무 낙엽 item5"
dotori -add -author woong -itemtype maya -inputpath /library/asset/06 -outputpath /library/asset -description "description6 about some details" -tag "나무 낙엽 item6"
dotori -add -author woong -itemtype maya -inputpath /library/asset/07 -outputpath /library/asset -description "description7 about some details" -tag "나무 낙엽 item7"
dotori -add -author woong -itemtype maya -inputpath /library/asset/08 -outputpath /library/asset -description "description8 about some details" -tag "나무 낙엽 item8"
dotori -add -author woong -itemtype maya -inputpath /library/asset/09 -outputpath /library/asset -description "description9 about some details" -tag "나무 낙엽 item9"
dotori -add -author woong -itemtype maya -inputpath /library/asset/10 -outputpath /library/asset -description "description10 about some details" -tag "나무 낙엽 item10"
dotori -add -author woong -itemtype maya -inputpath /library/asset/11 -outputpath /library/asset -description "description11 about some details" -tag "나무 낙엽 item11"
dotori -add -author woong -itemtype maya -inputpath /library/asset/12 -outputpath /library/asset -description "description12 about some details" -tag "나무 낙엽 item12"
dotori -add -author woong -itemtype maya -inputpath /library/asset/13 -outputpath /library/asset -description "description13 about some details" -tag "나무 낙엽 item13"
dotori -add -author woong -itemtype maya -inputpath /library/asset/14 -outputpath /library/asset -description "description14 about some details" -tag "나무 낙엽 item14"
dotori -add -author woong -itemtype maya -inputpath /library/asset/15 -outputpath /library/asset -description "description15 about some details" -tag "나무 낙엽 item15"
dotori -add -author woong -itemtype maya -inputpath /library/asset/16 -outputpath /library/asset -description "description16 about some details" -tag "나무 낙엽 item16"
dotori -add -author woong -itemtype maya -inputpath /library/asset/17 -outputpath /library/asset -description "description17 about some details" -tag "나무 낙엽 item17"

#REST API
curl -X POST -d "author=bae&itemtype=nuke&inputpath=/library/asset&outputpath=/library/asset&thumbimg=/library/asset&thumbmov=/library/asset" http://127.0.0.1/api/item
#DELETE는 실행 전에 직접 id를 넣어준 후 테스트한다.
#curl -X DELETE "http://127.0.0.1/api/item?type=nuke&id="