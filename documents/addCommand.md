## Item 추가
### maya, houdini, blender
```bash
$ sudo dotori -add -itemtype maya -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -inputthumbimgpath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.jpg -inputthumbclippath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.mov -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/maya/maya_scene.ma
```
- itemtype
- title
- author
- description
- tag
- inputthumbimgpath
- inputhumbclippath
- inputdatapath
### clip
```bash
$ sudo dotori -add -itemtype clip -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -fps 24 -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.mov
```
- itemtype
- title
- author
- description
- fps
- tag
- inputdatapath
