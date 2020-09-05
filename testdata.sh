#!/bin/sh
#DB에 아이템 추가 테스트를 위한 명령어를 모아둔 파일
sudo dotori -add -itemtype maya -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -inputthumbimgpath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.jpg -inputthumbclippath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.mov -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/maya/maya_scene.ma
sudo dotori -add -itemtype maya -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -inputthumbimgpath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.jpg -inputthumbclippath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.mov -inputdatapath "'/Users/seoyoungbae/git/fork/dotori/examples/maya/maya_scene.ma' '/Users/seoyoungbae/git/fork/dotori/examples/maya/maya_scene.mb'"
sudo dotori -add -itemtype houdini -title example -author woong -description "description2 about some details" -tag "나무 낙엽 item2" -inputthumbimgpath /Users/seoyoungbae/git/fork/dotori/examples/houdini/thumbnail.png -inputthumbclippath /Users/seoyoungbae/git/fork/dotori/examples/houdini/houdini_thumbnail.mov -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/houdini/houdini_scene.hda
sudo dotori -add -itemtype houdini -title example -author woong -description "description3 about some details" -tag "나무 낙엽 item3" -inputthumbimgpath /Users/seoyoungbae/git/fork/dotori/examples/houdini/thumbnail.png -inputthumbclippath /Users/seoyoungbae/git/fork/dotori/examples/houdini/houdini_thumbnail.mov -inputdatapath "'/Users/seoyoungbae/git/fork/dotori/examples/houdini/houdini_scene.hda' '/Users/seoyoungbae/git/fork/dotori/examples/houdini/houdini_scene.hip'"
sudo dotori -add -itemtype blender -title example -author woong -description "description4 about some details" -tag "자동차 트럭" -inputthumbimgpath /Users/seoyoungbae/git/fork/dotori/examples/blender/thumbnail.png -inputthumbclippath /Users/seoyoungbae/git/fork/dotori/examples/blender/thumbnail.mov -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/blender/data.blend
dotori -add -itemtype maya -title example -author bailey -description "description5 about some details" -tag "자동차 SUV" -inputthumbimgpath -inputthumbclippath -inputdatapath
dotori -add -itemtype maya -title example -author bailey -description "description6 about some details" -tag "하늘 노을" -inputthumbimgpath -inputthumbclippath -inputdatapath
dotori -add -itemtype maya -title example -author bailey -description "description7 about some details" -tag "하늘 밤" -inputthumbimgpath -inputthumbclippath -inputdatapath
dotori -add -itemtype maya -title example -author bailey -description "description8 about some details" -tag "하늘 아침" -inputthumbimgpath -inputthumbclippath -inputdatapath
dotori -add -itemtype maya -title example -author bailey -description "description9 about some details" -tag "기차 KTX" -inputthumbimgpath -inputthumbclippath -inputdatapath
dotori -add -itemtype maya -title example -author bailey -description "description10 about some details" -tag "기차 무궁화호" -inputthumbimgpath -inputthumbclippath -inputdatapath
dotori -add -itemtype usd -title example -author bailey -description "description10 about some details" -tag "기차 무궁화호" -inputthumbimgpath -inputthumbclippath -inputdatapath
dotori -add -itemtype sound -title example -author bailey -description "description10 about some details" -tag "bgm" -inputthumbimgpath -inputthumbclippath -inputdatapath
dotori -add -itemtype footage -title example -author woong -description "footage wildcard test" -fps 24 -tag "wildcard" -incolorspace "ACES - ACES2065-1" -outcolorspace "Output - Rec.709" -inputdatapath "/Users/woong/tdcourse_examples/footage/exr_aces_2065_1/A005C021_150831_R0D0.*.exr"

#아이템 삭제 커맨드를 실행 전에 직접 id를 입력해준 후 테스트한다.
#sudo dotori -remove -itemid 

#사용자의 accesslevel을 수정하는 testdata. 실행 전에 userid를 입력해준 후 테스트한다.
#sudo dotori -accesslevel admin -userid

#REST API
curl -X POST -d "author=bae&itemtype=maya" http://127.0.0.1/api/item
#DELETE는 실행 전에 직접 id를 넣어준 후 테스트한다.
#curl -X DELETE "http://127.0.0.1/api/item?type=nuke&id="