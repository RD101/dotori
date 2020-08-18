## Item 추가
### Maya, Houdini, Blender, Modo, Katana, OpenVDB, USD, Alembic
```bash
$ sudo dotori -add -itemtype modo -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -inputthumbimgpath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.jpg -inputthumbclippath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.mov -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/modo/data.lxo
```
- itemtype
- title
- author
- description
- tag
- inputthumbimgpath
- inputhumbclippath
- inputdatapath
### clip, unreal
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

### pdf, hwp
```bash
$ sudo dotori -add -itemtype pdf -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/pdf/지식재산권의기초.pdf
```
- itemtype
- title
- author
- description
- tag
- inputdatapath

### footage
```bash
$ sudo dotori -add -itemtype footage -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -fps 24 -incolorspace "ACES - ACES2065-1" -outcolorspace "Output - Rec.709" -inputdatapath "/Users/seoyoungbae/git/lazypic/tdcourse_examples/footage/exr_linear/A005C021_150831_R0D0.156404.exr /Users/seoyoungbae/git/lazypic/tdcourse_examples/footage/exr_linear/A005C021_150831_R0D0.156405.exr /Users/seoyoungbae/git/lazypic/tdcourse_examples/footage/exr_linear/A005C021_150831_R0D0.156406.exr"
```
- itemtype
- title
- author
- description
- fps
- incolorspace
- outcolorspace
- tag
- inputdatapath
