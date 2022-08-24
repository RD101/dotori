#!/bin/sh

example=$(pwd)/examples
etcExample="/home/chaeyun.bae/app/tdcourse_examples"
TOKEN=""
serverIp=$1

echo ""
echo "----alembic-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/abc/abc_thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/abc/abc_thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/abc/data.abc;type=application/octet-stream" \
-F "itemtype=alembic" \
-F "title=abc restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----blender-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/blender/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/blender/thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/blender/data.blend;type=application/octet-stream" \
-F "itemtype=blender" \
-F "title=blender restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----clip-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/blender/thumbnail.mov;type=video/quicktime" \
-F "itemtype=clip" \
-F "title=clip restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----footage-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$etcExample/footage/exr_aces_2065_1/A005C021_150831_R0D0.156404.exr;type=application/octet-stream" \
-F "file2=@$etcExample/footage/exr_aces_2065_1/A005C021_150831_R0D0.156404.exr;type=application/octet-stream" \
-F "file3=@$etcExample/footage/exr_aces_2065_1/A005C021_150831_R0D0.156405.exr;type=application/octet-stream" \
-F "file4=@$etcExample/footage/exr_aces_2065_1/A005C021_150831_R0D0.156406.exr;type=application/octet-stream" \
-F "file5=@$etcExample/footage/exr_aces_2065_1/A005C021_150831_R0D0.156407.exr;type=application/octet-stream" \
-F "file6=@$etcExample/footage/exr_aces_2065_1/A005C021_150831_R0D0.156407.exr;type=application/octet-stream" \
-F "itemtype=footage" \
-F "title=footage restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----fusion360-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/fusion360/thumbnail.png;type=image/png" \
-F "file2=@$example/fusion360/thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/fusion360/example.f3d;type=application/octet-stream" \
-F "file4=@$example/fusion360/example.step;type=application/octet-stream" \
-F "itemtype=fusion360" \
-F "title=fusion360 restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----hdri-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$etcExample/hdri/night_city.hdr;type=application/octet-stream" \
-F "itemtype=hdri" \
-F "title=hdri restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----houdini-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/houdini/thumbnail.png;type=image/png" \
-F "file2=@$example/houdini/houdini_thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/houdini/houdini_scene.hip;type=application/octet-stream" \
-F "file4=@$example/houdini/houdini_scene.hda;type=application/octet-stream" \
-F "itemtype=houdini" \
-F "title=houdini restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----hwp-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-F "file1=@$example/hwp/2020-표준취업규칙.hwp;type=application/octet-stream" \
-F "itemtype=hwp" \
-F "title=hwp restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----ies-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/ies/ies_example1.ies;type=application/octet-stream" \
-X POST -F "file2=@$example/ies/ies_example2.ies;type=application/octet-stream" \
-F "itemtype=ies" \
-F "title=ies restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----katana-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/katana/project_spacepod.katana;type=application/octet-stream" \
-F "itemtype=katana" \
-F "title=katana restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----lut-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/lut/ARRI_LogC2Video_709_adobe3d_33.cube;type=application/octet-stream" \
-F "itemtype=lut" \
-F "title=lut restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----maya-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/maya/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/maya/thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/maya/maya_scene.mb;type=application/octet-stream" \
-F "file4=@$example/maya/maya_scene.ma;type=application/octet-stream" \
-F "itemtype=maya" \
-F "title=maya restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----modo-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/modo/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/modo/thumbnail.mov;type=video/quicktime" \
-F "file4=@$example/modo/data.lxo;type=application/octet-stream" \
-F "itemtype=modo" \
-F "title=modo restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----nuke-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST \
-F "file1=@$example/nuke/nuke_scene.mov;type=video/quicktime" \
-F "file2=@$example/nuke/nuke_scene.nk;type=application/octet-stream" \
-F "file3=@$example/nuke/nuke_gizmo.gizmo;type=application/octet-stream" \
-F "itemtype=nuke" \
-F "title=nuke restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----openvdb-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/openvdb/vdb_thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/openvdb/vdb_thumbnail.mp4;type=video/mp4" \
-F "file4=@$example/openvdb/data.vdb;type=application/octet-stream" \
-F "itemtype=openvdb" \
-F "title=openvdb restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----pdf-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/pdf/지식재산권의기초.pdf;type=application/pdf" \
-F "itemtype=pdf" \
-F "title=pdf restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----ppt-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/ppt/powerpoint.pptx;type=application/octet-stream" \
-F "itemtype=ppt" \
-F "title=ppt restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----sound-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/sound/sample.mp3;type=audio/mp3" \
-F "itemtype=sound" \
-F "title=sound restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----texture-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/texture/texture_example.tif;type=image/tiff" \
-X POST -F "file1=@$example/texture/texutre_example.jpg;type=image/jpeg" \
-F "itemtype=texture" \
-F "title=texture restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item


echo ""
echo "----unreal-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/unreal/example.cpp;type=application/octet-stream" \
-X POST -F "file2=@$example/unreal/NewBlueprint.uasset;type=application/octet-stream" \
-F "itemtype=unreal" \
-F "title=unreal restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item

echo ""
echo "----usd-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/usd/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/usd/thumbnail.mov;type=video/quicktime" \
-F "file4=@$example/usd/data.usdc;type=application/octet-stream" \
-F "itemtype=usd" \
-F "title=usd restapi test" \
-F "author=dchecheb" \
-F "description=3" \
-F "tags=test" \
-F "attribute=key1:value1,key2:value2" \
http://$serverIp/api/item