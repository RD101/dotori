
function addAttribute() {
    let childnum = document.getElementById("attributes").childElementCount;
    let e = document.createElement('div');
    e.className = "row";
    let html = `
        <div class="col pt-2">
            <div class="form-group p-0 m-0">
                <input type="text" class="form-control" placeholder="key" value="" id="key${childnum}" name="key${childnum}">
            </div>
        </div>
        <div class="col pt-2">
            <div class="form-group p-0 m-0">
                <input type="text" class="form-control" placeholder="value" value="" id="value${childnum}" name="value${childnum}">
            </div>			
        </div>
    `
    e.innerHTML = html;
    document.getElementById("attributes").appendChild(e);
    // 최종 생성된 Attributes 갯수를 attributeNum에 저장한다.
    document.getElementById("attributesNum").value = document.getElementById("attributes").childElementCount;
}

// Hotkey: http://gcctech.org/csc/javascript/javascript_keycodes.htm
document.onkeydown = function(e) {
    // 인풋창에서는 화살표를 움직였을 때 페이지가 이동되면 안된다.
    if (event.target.tagName === "INPUT") {
        return
    }
    if (e.which == 37) { // arrow left
        document.getElementById("previous").click();
    } else if (e.which == 39) { // arrow right
        document.getElementById("next").click();
    }
};



// setDetailViewModal 은 아이템을 선택했을 때 볼 수 있는 detailview 모달창에 어셋 정보를 세팅해주는 함수이다.
function setDetailViewModal(itemid) {
    // Detail View에 세팅할 아이템 정보를 RestAPI로 불러옴
    $.ajax({
        url: `/api/item?id=${itemid}`,
        type: "get",
        dataType: "json",
        success: function(response) {
            // thumbnail 세팅
            let itemtype = response["itemtype"];
            if (itemtype == "pdf" || itemtype == "ppt" || itemtype == "hwp" || itemtype == "sound" || itemtype == "ies") {
                document.getElementById("modal-detailview-thumbnail").innerHTML = `<img src="/assets/img/${itemtype}thumbnail.svg">`;
            } else if (itemtype == "hdri" || itemtype == "lut" || itemtype == "texture") {
                if (response["status"] == "done") {
                    document.getElementById("modal-detailview-thumbnail").innerHTML = `<img src="/mediadata?id=${itemid}&type=png">`
                } else {
                    document.getElementById("modal-detailview-thumbnail").innerHTML = `<img src="/assets/img/noimage.svg">`
                }
            } else if (itemtype == "clip") {
                let thumbnailHtml = `
                                    <video id="modal-detailview-video" controls autoplay loop>
                                        <source src="/mediadata?id=${itemid}&type=mp4" type="video/mp4">
                                        <source src="/mediadata?id=${itemid}&type=ogg" type="video/ogg">
                                        Your browser does not support the video tag.
                                    </video>
                                    `
                document.getElementById("modal-detailview-thumbnail").innerHTML = thumbnailHtml;
                if (response["status"] != "done") {
                    document.getElementById("modal-detailview-video").setAttribute("poster", "/assets/img/noimage.svg") 
                }
            } else {
                let thumbnailHtml = `
                                    <video id="modal-detailview-video" controls autoplay loop>
                                        <source src="/mediadata?id=${itemid}&type=mp4" type="video/mp4">
                                        <source src="/mediadata?id=${itemid}&type=ogg" type="video/ogg">
                                        Your browser does not support the video tag.
                                    </video>
                                    `
                document.getElementById("modal-detailview-thumbnail").innerHTML = thumbnailHtml;
                if (response["status"] == "done") {
                    document.getElementById("modal-detailview-video").setAttribute("poster", `/mediadata?id=${itemid}&type=png`) 
                } else {
                    document.getElementById("modal-detailview-video").setAttribute("poster", "/assets/img/noimage.svg") 
                }
            }

            // title, id, author, description 세팅
            document.getElementById("modal-detailview-title").innerHTML = response["title"] + `<button type="button" onclick="location.href='/edit${itemtype}?id=${itemid}'" class="btn btn-outline-warning float-right" id="modal-detailview-edit-button">Edit</span>`;
            document.getElementById("modal-detailview-itemid").innerHTML = itemid;
            document.getElementById("modal-detailview-itemtype").innerHTML = itemtype;
            document.getElementById("modal-detailview-author").innerHTML = response["author"];
            document.getElementById("modal-detailview-description").innerHTML = response["description"];
            
            // Tags 세팅
            let tagsHtml = `<strong>Tags</strong><br>`;
            for (let i=0; i<response["tags"].length;i++) {
                let tag = response["tags"][i];
                tagsHtml += `
                <a href="/search?searchword=tag:${tag}" class="tag badge badge-outline-darkmode">${tag}</a>
                `;
            }
            document.getElementById("modal-detailview-tags").innerHTML = tagsHtml
            
            // Attributes 세팅
            if (Object.keys(response["attributes"]).length != 0) {
                let attributesHtml = `<strong>Attributes</strong>
                                      <dl class="attributes">`;
                for (key in response["attributes"]) {
                    let value = response["attributes"][key];
                    attributesHtml += `
                                    <div class="row no-gutters">
                                        <dt>${key}</dt>
                                        <dd>${value}</dd>
                                    </div>
                                    `
                }
                attributesHtml += `</dl>`
                document.getElementById("modal-detailview-attributes").innerHTML = attributesHtml
            } else {
                document.getElementById("modal-detailview-attributes").innerHTML = ``
            }
            
            // buttons 세팅
            document.getElementById("modal-detailview-edit-button").href=`/edit${itemtype}?itemtype=${itemtype}&id=${itemid}`
            let outputdatapath=response["outputdatapath"]
            let footerHtml = `
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-download-button" onclick="location.href='/download-item?id=${itemid}'">Download</button>
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-copypath-button" onclick="copyPath('${outputdatapath}')">Copy Path</button>
            `
            let footerHtmlForAdmin=`
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-download-button" onclick="location.href='/download-item?id=${itemid}'">Download</button>
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-copypath-button" onclick="copyPath('${outputdatapath}')">Copy Path</button>
            <button type="button" class="btn btn-outline-danger" id="modal-detailview-delete-button" data-dismiss="modal" data-toggle="modal" data-target="#modal-rmitem">Delete</button>
            `
            if (document.getElementById("accesslevel").value == "admin") {
                console.log("test")
                document.getElementById("modal-rmitem-itemid").value = itemid;
                document.getElementById("modal-detailview-footer").innerHTML = footerHtmlForAdmin       // admin 계정일 때만 delete 버튼이 보인다.
            } else {
                document.getElementById("modal-detailview-footer").innerHTML = footerHtml
            }
            if (itemtype == "footage") {
                document.getElementById("modal-detailview-download-button").style.visibility="hidden"   // footage는 download 버튼이 보이지 않는다.  
            }
        },
        error: function() {
            alert("어셋 정보를 가져오는 데 실패했습니다");
            location.reload();
        }
    });

}

// rmItemModal 은 삭제 modal창에서 Delete 버튼을 누르면 실행되는 아이템 삭제 함수이다.
function rmItemModal(itemId) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?id=${itemId}`,
        type: "delete",
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function() {
            alert("id: "+itemId+"\n\n아이템 삭제를 성공했습니다."); 
            location.reload();
        },
        error: function(){
            alert("id: "+itemId+"\n\n아이템 삭제를 실패했습니다.");  
        }
    });
}

// handlerNumCheck 은 숫자만 적히도록 하는 레귤러익스프레션이다.
function handlerNumCheck(element){
    $(element).val($(element).val().replace(/[^0-9]/g,""));
    
    if(element.name == "umask" || element.name == "folderpermission" || element.name == "filepermission"){
        $(element).val($(element).val().replace(/[^0-7]/g,""));
    }
}

// addPageBlankCheck 은 addPage에서 빈값을 체크하는 함수이다.
function addPageBlankCheck(){
    if(document.getElementById("addTitle").value == ''){
        alert("Title을 입력해주세요.");
        return false;
    }
    if(document.getElementById("addAuthor").value == ''){
        alert("Author를 입력해주세요.");
        return false;
    }
    if(document.getElementById("addDescription").value == ''){
        alert("Description을 입력해주세요.");
        return false;
    }
    if(document.getElementById("addTags").value == ''){
        alert("Tags를 입력해주세요.");
        return false;
    }
    return true;
}

// toggoleItems는 상단의 체크박스를 클릭하면 작동하는 함수로, 모든 체크박스를 선택 혹은 선택해제 한다.
function toggleItems(){
    // 기준이 되는 체크박스의 상태값을 가지고 온다.
    let status = document.getElementById("toggle-checkbox").checked
    // 가져온 상태값을 기준으로 모든 체크박스의 상태를 설정한다.
    let checkboxes = document.querySelectorAll('*[name^="checkbox"]');
    for (i=0;i<checkboxes.length;i++) {
        checkboxes[i].checked = status
    }
}

// recentlyClick 은 초기페이지에서 최근등록된 아이템의 next, prev 버튼을 눌렀을때 실행하는 함수이다.
function recentlyClick(totalItemNum, buttonState) {
    // totalItemNum 최근에셋의 전체 아이템 수
    let totalPageNum = Math.ceil(totalItemNum / 4); // 전체 페이지 수
    let clearItemNum = (totalPageNum * 4) - totalItemNum; // 마지막 페이지의 공백처리할 아이템 수
    let currentPageNum = parseInt(document.getElementById("recentlyPage").getAttribute('value')); // 현재 보고있는 페이지

    if (buttonState==="init"){
        currentPageNum = 1;
    }
    else if (buttonState==="next"){
        if (currentPageNum===totalPageNum){
            currentPageNum = 1;
        }else{
            currentPageNum++;
        }
    }else if (buttonState==="prev"){
        if (currentPageNum===1){
            currentPageNum = totalPageNum;
        }else{
            currentPageNum--;
        }
    }

    document.getElementById("recentlyPage").innerHTML = currentPageNum + " / " + totalPageNum;
    document.getElementById("recentlyPage").setAttribute('value', currentPageNum);

    if(totalPageNum !== 1 && clearItemNum!==0){
        // 마지막 페이지일 때
        if(currentPageNum===totalPageNum){
            for(let i = 3; clearItemNum!=0; i--, clearItemNum--){
                document.getElementById("recentCard"+i).style.visibility="hidden"
            }
        }
        // 마지막 페이지가 아닐 때
        if (currentPageNum != totalPageNum) {
            for(let i = 0; i<4; i++){
                document.getElementById("recentCard"+i).style.visibility="visible"
            }
        }
    }

    // Get Favorite Asset IDs
    let userid = document.getElementById("userid").value;
    let token = document.getElementById("token").value;
    let favoriteAssetIds = new Array();
    $.ajax({
        url: `/api/favoriteasset?userid=${userid}`,
        headers: {
            "Authorization": "Basic " + token
        },
        type: "get",
        dataType: "json",
        success: function(response) {
            favoriteAssetIds = response["favoriteAssetIds"];
        },
        complete: function() {
            // Recent Item 정보 가져온 후 세팅
            $.ajax({
                url: `/api/recentitem?recentlypage=${currentPageNum}`,
                type: "get",
                dataType: "json",
                success: function(data) {
                    let thumbnailwidth = document.getElementById("thumbnailwidth").value;
                    let thumbnailheight = document.getElementById("thumbnailheight");
                    let img = ""
                    for (let i = 0; i < data.length; i++){
                        let itemid = data[i].id;
                        // 썸네일 스위칭
                        let recentlyImageForm = document.getElementById("recentlyImageForm"+i)
                        if (data[i].itemtype == "pdf" || data[i].itemtype == "ppt" || data[i].itemtype == "hwp" || data[i].itemtype == "sound" || data[i].itemtype == "ies") {
                            img = '<img class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                            '" src="/assets/img/' + data[i].itemtype + 'thumbnail.svg">'
                        } else if(data[i].itemtype=="hdri" || data[i].itemtype == "lut" || data[i].itemtype=="texture"){
                            if(data[i].status == "done"){
                                img = '<img class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" src="/mediadata?id=' + itemid + '&type=png">'
                            }else{
                                img = '<img class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" src="/assets/img/noimage.svg">'
                            }
                        } else if (data[i].itemtype=="clip") {
                            if(data[i].status == "done"){
                                img = '<video class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" autoplay loop>' +
                                        '<source src="/mediadata?id=' + itemid + '&type=mp4" type="video/mp4">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=ogg" type="video/ogg">' +
                                        'Your browser does not support the video tag.'+
                                        '</video>'
                            }else{
                                img = '<video class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" controls poster="/assets/img/noimage.svg">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=mp4" type="video/mp4">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=ogg" type="video/ogg">' +
                                        'Your browser does not support the video tag.'+
                                        '</video>'
                            }
                        } else{
                            if(data[i].status == "done"){
                                img = '<video class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" controls poster="/mediadata?id=' + itemid + '&type=png">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=mp4" type="video/mp4">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=ogg" type="video/ogg">' +
                                        'Your browser does not support the video tag.'+
                                        '</video>'
                            }else{
                                img = '<video class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" controls poster="/assets/img/noimage.svg">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=mp4" type="video/mp4">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=ogg" type="video/ogg">' +
                                        'Your browser does not support the video tag.'+
                                        '</video>'
                            }
                        }
                        recentlyImageForm.innerHTML = img;
                        document.getElementById("recentCardBody"+i).onclick = function () { setDetailViewModal(itemid);}    // detail view 용 item id 스위칭
                        document.getElementById("recentlyCreateTime"+i).innerHTML = data[i].createtime.split('T')[0];           // create time 스위칭
                        document.getElementById("recentUsingRate"+i).innerHTML = data[i].usingrate;                             // using rate 스위칭
                        // 태그 스위칭
                        let tagsHtml = '';
                        for (let j=0;j<data[i].tags.length;j++) {
                            tagsHtml += '<a href="/search?searchword=tag:' + data[i].tags[j] + '" class="tag badge badge-outline-darkmode">' + data[i].tags[j] + '</a>';
                        }
                        tagsHtml += `<div class="tag finger badge badge-outline-warning" onclick="copyPath('${data[i].outputdatapath}')">copypath</div>`
                        document.getElementById("recentCardTags"+i).innerHTML = tagsHtml;                             
                        // 즐겨찾기 아이콘 스위칭
                        let fillBool = "unfilled";
                        if (favoriteAssetIds){
                            for (j=0;j<favoriteAssetIds.length;j++) {
                                if (favoriteAssetIds[j] == itemid) {
                                    fillBool = "filled";
                                }
                            }
                        }
                        titleHtml= data[i].title;
                        titleHtml +=    `<div class="bookmark-icon"> 
                                            <div class="bookmark-clicklistener" onclick="clickBookmarkIcon(this, '${fillBool}','${itemid}');event.stopPropagation();"></div>
                                            <object type="image/svg+xml" data="/assets/img/bookmark-${fillBool}.svg" class="bookmark-${fillBool}-icon"></object>
                                        </div>`;
                        document.getElementById("recentlyTitle"+i).innerHTML = titleHtml;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        },
        error: function(response) {
            alert(response["responseText"]);
        }
    })
}

// topUsingClick 은 초기페이지에서 가장 많이 사용되는 아이템의 next, prev 버튼을 눌렀을 때 실행하는 함수이다.
function topUsingClick(totalItemNum, buttonState) {
    // RecentlyTotalNum 가장 많이 사용된 에셋의 전체 아이템 수
    let totalPageNum = Math.ceil(totalItemNum / 4); // 전체 페이지 수
    let clearItemNum = (totalPageNum * 4) - totalItemNum; // 마지막 페이지의 공백처리할 아이템 수
    let currentPageNum = parseInt(document.getElementById("topUsingPage").getAttribute('value'));

    if (buttonState==="init"){
        currentPageNum = 1;
    }
    else if (buttonState==="next"){
        if (currentPageNum===totalPageNum){
            currentPageNum = 1;
        }else{
            currentPageNum++;
        }
    }else if (buttonState==="prev"){
        if (currentPageNum===1){
            currentPageNum = totalPageNum;
        }else{
            currentPageNum--;
        }
    }
    
    document.getElementById("topUsingPage").innerHTML = currentPageNum + " / " + totalPageNum;
    document.getElementById("topUsingPage").setAttribute('value', currentPageNum);

    if(totalPageNum !== 1 && clearItemNum!=0){
        // 마지막 페이지일 때
        if(currentPageNum==totalPageNum){
            for(let i = 3; clearItemNum!=0; i--, clearItemNum--){
                document.getElementById("topUsingCard"+i).style.visibility="hidden";
            }
        }
        // 마지막 페이지가 아닐 때
        if (currentPageNum != totalPageNum) {
            for(let i = 0; i<4; i++){
                document.getElementById("topUsingCard"+i).style.visibility="visible"
            }
        }
    }

    // Get Favorite Asset IDs
    let userid = document.getElementById("userid").value;
    let token = document.getElementById("token").value;
    let favoriteAssetIds = new Array();
    $.ajax({
        url: `/api/favoriteasset?userid=${userid}`,
        headers: {
            "Authorization": "Basic " + token
        },
        type: "get",
        dataType: "json",
        success: function(response) {
            favoriteAssetIds = response["favoriteAssetIds"];
        },
        complete: function() {
            // Top Using Item 정보 가져온 후 세팅
            $.ajax({
                url: `/api/topusingitem?usingpage=${currentPageNum}`,
                type: "get",
                dataType: "json",
                success: function(data) {
                    let thumbnailwidth = document.getElementById("thumbnailwidth").value;
                    let thumbnailheight = document.getElementById("thumbnailheight").value;
                    let img = ""
                    for (let i = 0; i < data.length; i++){
                        let itemid = data[i].id;
                        let topUsingImageForm = document.getElementById("topUsingImageForm"+i)
                        if (data[i].itemtype == "pdf" || data[i].itemtype == "ppt" || data[i].itemtype == "hwp" || data[i].itemtype == "sound" || data[i].itemtype == "ies") {
                            img = '<img class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" src="/assets/img/' + data[i].itemtype + 'thumbnail.svg">'
                        } else if(data[i].itemtype=="hdri" || data[i].itemtype=="lut" || data[i].itemtype=="texture"){
                            if(data[i].status == "done"){
                                img = '<img class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" src="/mediadata?id=' + itemid + '&type=png">'
                            }else{
                                img = '<img class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" src="/assets/img/noimage.svg">'
                            }
                        } else if (data[i].itemtype=="clip") {
                            if(data[i].status == "done"){
                                img = '<video class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" autoplay loop>' +
                                        '<source src="/mediadata?id=' + itemid + '&type=mp4" type="video/mp4">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=ogg" type="video/ogg">' +
                                        'Your browser does not support the video tag.'+
                                        '</video>'
                            }else{
                                img = '<video class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" controls poster="/assets/img/noimage.svg">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=mp4" type="video/mp4">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=ogg" type="video/ogg">' +
                                        'Your browser does not support the video tag.'+
                                        '</video>'
                            }
                        } else {
                            if(data[i].status == "done"){
                                img = '<video class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" controls poster="/mediadata?id=' + itemid + '&type=png">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=mp4" type="video/mp4">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=ogg" type="video/ogg">' +
                                        'Your browser does not support the video tag.'+
                                        '</video>'
                            }else{
                                img = '<video class="card-img" width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                        '" controls poster="/assets/img/noimage.svg">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=mp4" type="video/mp4">' +
                                        '<source src="/mediadata?id=' + itemid + '&type=ogg" type="video/ogg">' +
                                        'Your browser does not support the video tag.'+
                                        '</video>'
                            }
                        }
                        topUsingImageForm.innerHTML = img;
                        document.getElementById("topUsingCardBody"+i).onclick = function () { setDetailViewModal(itemid);}  // detail view 용 item id 스위칭
                        document.getElementById("topUsingRate"+i).innerHTML = data[i].usingrate;                            // using rate 스위칭
                        // 태그 스위칭
                        let tagsHtml = '';
                        for (let j=0;j<data[i].tags.length;j++) {
                            tagsHtml += '<a href="/search?searchword=tag:' + data[i].tags[j] + '" class="tag badge badge-outline-darkmode">' + data[i].tags[j] + '</a>';
                        }
                        tagsHtml += `<div class="tag finger badge badge-outline-warning" onclick="copyPath('${data[i].outputdatapath}')">copypath</div>`
                        document.getElementById("topUsingCardTags"+i).innerHTML = tagsHtml;  
                        // 즐겨찾기 아이콘 스위칭
                        let fillBool = "unfilled";

                        if (favoriteAssetIds){
                            for (j=0;j<favoriteAssetIds.length;j++) {
                                if (favoriteAssetIds[j] == itemid) {
                                    fillBool = "filled";
                                }
                            }
                        }
                        titleHtml= data[i].title;
                        titleHtml +=    `<div class="bookmark-icon"> 
                                            <div class="bookmark-clicklistener" onclick="clickBookmarkIcon(this, '${fillBool}','${itemid}');event.stopPropagation();"></div>
                                            <object type="image/svg+xml" data="/assets/img/bookmark-${fillBool}.svg" class="bookmark-${fillBool}-icon"></object>
                                        </div>`;
                        document.getElementById("topUsingTitle"+i).innerHTML = titleHtml;                             
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        },
        error: function(response) {
            alert("Failed to Get Favorite Asset IDs.");
        }

    })

}


// clickBookmarkIcon 은 즐겨찾기 버튼을 눌렀을 때 실행되는 함수다. 
function clickBookmarkIcon(target, fillBool, itemid) { 
    let parentNode = target.parentNode;
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;

    // 즐겨찾기에 추가
    if (fillBool == "unfilled") {
        $.ajax({
            url: "/api/favoriteasset",
            headers: {
                "Authorization": "Basic " + token
            },
            type: "post",
            data: {"itemid":itemid,"userid":userid},
            success: function(data) {
                parentNode.innerHTML = `<div class="bookmark-clicklistener" onclick="clickBookmarkIcon(this, 'filled','${itemid}');event.stopPropagation();"></div>
                                        <object type="image/svg+xml" data="/assets/img/bookmark-filled.svg" class="bookmark-filled-icon"></object>`
            },
            error: function(response) {
                alert(response["responseText"]);
            }
        })
    // 즐겨찾기에서 제거
    } else if (fillBool == "filled") {
        
        $.ajax({
            url: `/api/favoriteasset?itemid=${itemid}&userid=${userid}`,
            headers: {
                "Authorization": "Basic " + token
            },
            type: "delete",
            success: function() {
                parentNode.innerHTML = `<div class="bookmark-clicklistener" onclick="clickBookmarkIcon(this, 'unfilled','${itemid}');event.stopPropagation();"></div>
                                        <object type="image/svg+xml" data="/assets/img/bookmark-unfilled.svg" class="bookmark-unfilled-icon"></object>`
            },
            error: function(response) {
                alert(response["responseText"]);
            }
        })
    }    
}

// initPasword 함수는 비밀번호 초기화 버튼을 눌렀을 때 실행되는 함수이다.
function initPasword() {
    let token = document.getElementById("token").value;
    let checkboxes = document.querySelectorAll('*[name^="checkbox"]:checked');

    for (i = 0; i < checkboxes.length; i++) {
        let userID = checkboxes[i].value;
        $.ajax({
            url: "/api/initpassword",
            headers: {
                "Authorization": "Basic " + token
            },
            type: "post",
            data: {
                id: userID,
            },
            success: function() {
                alert(userID+" 사용자의 패스워드가 초기화 되었습니다");
                location.reload();
            },
            error: function() {
                alert(userID+" 사용자의 패스워드를 초기화하는 데 실패하였습니다.")
            }
        })
    }
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}


function SetAutoplay() {
    fetch('/api/user?autoplay='+document.getElementById("autoplay").checked, {
        method: 'PUT',
        headers: {
            "Authorization": "Basic "+ getCookie("SessionToken"),
        },
    })
    .then((response) => {
        if (!response.ok) {
            return response.text().then((data) => {
                alert(data);
                return data;
            })
        }
        if (response.ok) {
            return response.json().then((data) => {
                return
            })
        }
    })
    .catch((err) => {
        alert(err);
        return err
    });
}

// copyClipboard 는 value 값을 받아서, 클립보드로 복사하는 기능이다.
function copyClipboard(value) {
    let id = document.createElement("input");   // input요소를 만듬
    id.setAttribute("value", value);            // input요소에 값을 추가
    document.body.appendChild(id);              // body에 요소 추가
    id.select();                                // input요소를 선택
    document.execCommand("copy");               // 복사기능 실행
    document.body.removeChild(id);              // body에 요소 삭제

    // Toast 띄우기
    tata.success('Copy Clipboard', "Data path copyed!", {
        position: 'tr',
        duration: 1000,
        onClose: null,
    })
}

// copyPath 함수는 아이디값을 받아서, 클립보드로 복사하는 기능이다.
function copyPath(path) {
    let windowsUNCPrefix = "";
    $.ajax({
        url: `/api/adminsetting`,
        type: "get",
        headers: {
            "Authorization": "Basic " + document.getElementById("token").value
        },
        dataType: "json",
        async: false,
        success: function(data) {
            windowsUNCPrefix = data.windowsuncprefix;
        },
        error: function(){
            alert("admin 셋팅에서 Windows UNC Prefix 값을 가지고 올 수 없습니다.");  
        }
    });
    if (navigator.userAgent.indexOf("Win") != -1) { // windows 경우
        path = windowsUNCPrefix + path.replace(/\//g, "\\")
    }
    copyClipboard(path)
}