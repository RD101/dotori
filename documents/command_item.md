## Item 추가

### Maya, Houdini, Blender, Modo, Katana, OpenVDB, USD, Alembic, Nuke, Fusion360, Max

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

### clip

```bash
$ sudo dotori -add -itemtype clip -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -fps 24 -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.mov
```
- itemtype
- title
- author
- description
- tag
- inputdatapath
- fps

### PDF, HWP, PPT, Sound, Texture, Unreal, IES

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
- tag
- inputdatapath
- fps
- incolorspace
- outcolorspace

### LUT

```bash
$ sudo dotori -add -itemtype lut -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -inputthumbimgpath /Users/seoyoungbae/git/fork/dotori/examples/maya/thumbnail.jpg -inputdatapath /Users/seoyoungbae/git/fork/dotori/examples/lut/ARRI_LogC2Video_709_adobe3d_33.cube
```

- itemtype
- title
- author
- description
- tag
- inputthumbimgpath
- inputdatapath

### HDRI

```bash
$ sudo dotori -add -itemtype hdri -title example -author woong -description "description1 about some details" -tag "나무 낙엽 item1" -incolorspace "ACES - ACES2065-1" -outcolorspace "Output - Rec.709" -inputdatapath /Users/seoyoungbae/git/lazypic/tdcourse_examples/hdri/night_city.hdr
```

- itemtype
- title
- author
- description
- tag
- inputdatapath
- incolorspace
- outcolorspace

## Item 삭제

```bash
$ sudo dotori -remove -itemid 5e89cef7cd1747fd5eacf256
```

## 시퀀스, 클립 검색

inputdatapath, seek 옵션을 붙히면 json 정보로 반환합니다.

```bash
dotori -inputdatapath /search/path -seek

[{"incolorspace":"","outcolorspace":"","renderin":0,"renderout":0,"searchpath":"/Users/woong/tdcourse_examples","convertext":".dpx","dir":"/Users/woong/tdcourse_examples/ACES_Plate","base":"LogC_ref_Isabella.%01d.dpx","ext":".dpx","digitnum":0,"framein":0,"frameout":0,"width":0,"height":0,"timecodein":"","timecodeout":"","length":1,"inputcodec":"","outputcodec":"","fps":0,"rollmedia":"","error":""},{"incolorspace":"","outcolorspace":"","renderin":156404,"renderout":156408,"searchpath":"/Users/woong/tdcourse_examples","convertext":".dpx","dir":"/Users/woong/tdcourse_examples/footage/dpx_alexaV3LogC_or_ocio_input_arri_v3LogC_EI800","base":"A005C021_150831_R0D0.%06d.dpx","ext":".dpx","digitnum":0,"framein":156404,"frameout":156408,"width":0,"height":0,"timecodein":"","timecodeout":"","length":5,"inputcodec":"","outputcodec":"","fps":0,"rollmedia":"","error":""},{"incolorspace":"","outcolorspace":"","renderin":0,"renderout":0,"searchpath":"/Users/woong/tdcourse_examples","convertext":".mp4","dir":"/Users/woong/tdcourse_examples/movs","base":"H264_1280x1280_framenum_25fps.mov","ext":".mov","digitnum":0,"framein":0,"frameout":0,"width":0,"height":0,"timecodein":"","timecodeout":"","length":0,"inputcodec":"","outputcodec":"","fps":0,"rollmedia":"","error":""},{"incolorspace":"","outcolorspace":"","renderin":0,"renderout":0,"searchpath":"/Users/woong/tdcourse_examples","convertext":".mp4","dir":"/Users/woong/tdcourse_examples/movs","base":"H264_2048x1556_24fps.mov","ext":".mov","digitnum":0,"framein":0,"frameout":0,"width":0,"height":0,"timecodein":"","timecodeout":"","length":0,"inputcodec":"","outputcodec":"","fps":0,"rollmedia":"","error":""}]
```



