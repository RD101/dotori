#!/bin/sh

example=$(pwd)/examples
etcExample="/home/chaeyun.bae/app/tdcourse_examples"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNoYWV5dW4uYmFlIiwiYWNjZXNzbGV2ZWwiOiJhZG1pbiJ9.tI7mnGPCoYh4HTJun3tyOCd6FMGkjyQg5Z9nTfE9cGA"
serverIp=$1

echo ""
echo "----alembic-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/abc/abc_thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/abc/abc_thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/abc/data.abc;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"alembic\",\"title\":\"abc restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----blender-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/blender/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/blender/thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/blender/data.blend;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"blender\",\"title\":\"blender restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----blender-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/blender/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/blender/thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/blender/data.blend;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"blender\",\"title\":\"blender restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----clip-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/blender/thumbnail.mov;type=video/quicktime" \
-F "iteminfo={\"itemtype\":\"clip\",\"title\":\"clip restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
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
-F "iteminfo={\"itemtype\":\"footage\",\"title\":\"footage restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----fusion360-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/fusion360/thumbnail.png;type=image/png" \
-F "file2=@$example/fusion360/thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/fusion360/example.f3d;type=application/octet-stream" \
-F "file4=@$example/fusion360/example.step;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"fusion360\",\"title\":\"fusion360 restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----hdri-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$etcExample/hdri/night_city.hdr;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"hdri\",\"title\":\"hdri restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----houdini-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/houdini/thumbnail.png;type=image/png" \
-F "file2=@$example/houdini/houdini_thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/houdini/houdini_scene.hip;type=application/octet-stream" \
-F "file4=@$example/houdini/houdini_scene.hda;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"houdini\",\"title\":\"houdini restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----hwp-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-F "file1=@$example/hwp/2020-표준취업규칙.hwp;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"hwp\",\"title\":\"hwp restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----ies-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/ies/ies_example1.ies;type=application/octet-stream" \
-X POST -F "file2=@$example/ies/ies_example2.ies;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"ies\",\"title\":\"ies restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----katana-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/katana/project_spacepod.katana;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"katana\",\"title\":\"katana restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----lut-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/lut/ARRI_LogC2Video_709_adobe3d_33.cube;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"lut\",\"title\":\"lut restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----maya-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/maya/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/maya/thumbnail.mov;type=video/quicktime" \
-F "file3=@$example/maya/maya_scene.mb;type=application/octet-stream" \
-F "file4=@$example/maya/maya_scene.ma;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"maya\",\"title\":\"maya restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----modo-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/modo/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/modo/thumbnail.mov;type=video/quicktime" \
-F "file4=@$example/modo/data.lxo;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"modo\",\"title\":\"modo restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----nuke-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/nuke/nuke_scene.jpg;type=image/jpeg" \
-F "file2=@$example/nuke/nuke_scene.mov;type=video/quicktime" \
-F "file4=@$example/nuke/nuke_scene.nk;type=application/octet-stream" \
-F "file4=@$example/nuke/nuke_gizmo.gizmo;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"nuke\",\"title\":\"nuke restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----openvdb-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/openvdb/vdb_thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/openvdb/vdb_thumbnail.mp4;type=video/mp4" \
-F "file4=@$example/openvdb/data.vdb;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"openvdb\",\"title\":\"openvdb restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----pdf-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/pdf/지식재산권의기초.pdf;type=application/pdf" \
-F "iteminfo={\"itemtype\":\"pdf\",\"title\":\"pdf restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----ppt-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/ppt/powerpoint.pptx;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"ppt\",\"title\":\"ppt restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----sound-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/sound/sample.mp3;type=audio/mp3" \
-F "iteminfo={\"itemtype\":\"sound\",\"title\":\"sound restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----texture-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/texture/texture_example.tif;type=image/tiff" \
-X POST -F "file1=@$example/texture/texutre_example.jpg;type=image/jpeg" \
-F "iteminfo={\"itemtype\":\"texture\",\"title\":\"texture restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item


echo ""
echo "----unreal-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/unreal/example.cpp;type=application/octet-stream" \
-X POST -F "file2=@$example/unreal/NewBlueprint.uasset;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"unreal\",\"title\":\"unreal restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item

echo ""
echo "----usd-----"
echo ""
curl -H "Authorization: Basic $TOKEN" \
-X POST -F "file1=@$example/usd/thumbnail.jpg;type=image/jpeg" \
-F "file2=@$example/usd/thumbnail.mov;type=video/quicktime" \
-F "file4=@$example/usd/data.usdc;type=application/octet-stream" \
-F "iteminfo={\"itemtype\":\"usd\",\"title\":\"usd restapi test\",\"author\":\"dchecheb\",\"description\":\"3\",\"tags\":\"test\",\"attribute\":\"key1:value1,key2:value2\"}" \
http://$serverIp/api/item