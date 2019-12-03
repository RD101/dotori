#!/bin/sh
#DB에 아이템 추가 테스트를 위한 명령어를 모아둔 파일
dotori -add -author bae -type nuke -inputpath /library/asset -outputpath /library/asset
dotori -add -author woong -type maya -inputpath /library/asset -outputpath /library/asset
dotori -add -author liah -type houdini -inputpath /library/asset -outputpath /library/asset

#REST API
curl -d "author=bae&type=nuke&inputpath=/library/asset&outputpath=/library/asset&thumbimg=/library/asset&thumbmov=/library/asset" http://127.0.0.1/api/add