
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

// copyButton 은 아이디값을 받아서, 클립보드로 복사하는 기능이다.
function copyButton(elementId) {
    let id = document.createElement("input");                       // input요소를 만듬
    id.setAttribute("value", elementId);                            // input요소에 값을 추가
    document.getElementById("modal-detailview").appendChild(id);    // modal에 요소 추가
    id.select();                                                    // input요소를 선택
    document.execCommand("copy");                                   // 복사기능 실행
    document.getElementById("modal-detailview").removeChild(id);    // modal에서 요소 삭제
}

// setDetailViewModal 은 아이템을 선택했을 때 볼 수 있는 detailview 모달창에 detail 정보를 세팅해주는 함수이다.
function setDetailViewModal(itemid) {
    $.ajax({
        url: `/api/item?id=${itemid}`,
        type: "get",
        dataType: "json",
        success: function(response) {
            document.getElementById("modal-detailview-title").innerHTML = response["title"];
            document.getElementById("modal-detailview-itemid").innerHTML = itemid;
            document.getElementById("modal-detailview-author").innerHTML = response["author"];
            document.getElementById("modal-detailview-description").innerHTML = response["description"];
            let outputdatapath=response["outputdatapath"]
            let footerHtml = `
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-download-button" onclick="location.href='/download-item?id=${itemid}'">Download</button>
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-copypath-button" onclick="copyButton('${outputdatapath}')">Copy Path</button>
            `
            let footerHtmlForAdmin=`
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-download-button" onclick="location.href='/download-item?id=${itemid}'">Download</button>
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-copypath-button" onclick="copyButton('${outputdatapath}')">Copy Path</button>
            <button type="button" class="btn btn-outline-danger" id="modal-detailview-delete-button" data-dismiss="modal" data-toggle="modal" data-target="#modal-rmitem">Delete</button>
            `
            if (document.getElementById("accesslevel").value == "admin") {
                document.getElementById("modal-rmitem-itemid").value = itemid;
                document.getElementById("modal-detailview-footer").innerHTML = footerHtmlForAdmin
            } else {
                document.getElementById("modal-detailview-footer").innerHTML = footerHtml
            }
            if (response["itemtype"] == "footage") {
                document.getElementById("modal-detailview-download-button").style.visibility="hidden"        
            }
        },
        error: function(result) {
            alert(result);
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
    
    if (buttonState=="next"){
        if (currentPageNum===totalPageNum){
            currentPageNum = 1;
        }else{
            currentPageNum++;
        }
    }else{
        if (currentPageNum===1){
            currentPageNum = totalPageNum;
        }else{
            currentPageNum--;
        }
    }

    document.getElementById("recentlyPage").innerHTML = currentPageNum + " / " + totalPageNum;
    if(clearItemNum!==0 && currentPageNum===totalPageNum){
        for(let i = 3; clearItemNum!=0; i--, clearItemNum--){
            document.getElementById("recentlyImageForm"+i).innerHTML = ""
            document.getElementById("recentlyTitle"+i).innerHTML = ""
            document.getElementById("recentlyAuthor"+i).innerHTML = ""
            document.getElementById("recentlyCreateTime"+i).innerHTML = ""
        }
    }
    document.getElementById("recentlyPage").setAttribute('value', currentPageNum);
    $.ajax({
        url: `/api/recentitem?recentlypage=${currentPageNum}`,
        type: "get",
        dataType: "json",
        success: function(data) {
            let thumbnailwidth = document.getElementById("thumbnailwidth").value;
            let thumbnailheight = document.getElementById("thumbnailheight").value;
            let img = ""
            for (let i = 0; i < data.length; i++){
                let recentlyImageForm = document.getElementById("recentlyImageForm"+i)
                if (data[i].itemtype=="pdf"){
                    img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/assets/img/pdfthumbnail.svg">'
                }else if(data[i].itemtype=="hwp"){
                    img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/assets/img/hwpthumbnail.svg">'
                }else if(data[i].itemtype=="sound"){
                    img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/assets/img/soundthumbnail.svg">'
                }else if(data[i].itemtype=="hdri" || data[i].itemtype=="texture"){
                    if(data[i].status == "done"){
                        img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/mediadata?id=' + data[i].id + '&type=png">'
                    }else{
                        img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/assets/img/noimage.svg">'
                    }
                }else{
                    if(data[i].status == "done"){
                        img = '<video width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" controls poster="/mediadata?id=' + data[i].id + '&type=png">' +
                                '<source src="/mediadata?id=' + data[i].id + '&type=mp4" type="video/mp4">' +
                                '<source src="/mediadata?id=' + data[i].id + '&type=ogg" type="video/ogg">' +
                                'Your browser does not support the video tag.'+
                                '</video>'
                    }else{
                        img = '<video width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" controls poster="/assets/img/noimage.svg">' +
                                '<source src="/mediadata?id=' + data[i].id + '&type=mp4" type="video/mp4">' +
                                '<source src="/mediadata?id=' + data[i].id + '&type=ogg" type="video/ogg">' +
                                'Your browser does not support the video tag.'+
                                '</video>'
                    }
                recentlyImageForm.innerHTML = img;
                }
                document.getElementById("recentlyTitle"+i).innerHTML = "Title: " + data[i].title;
                document.getElementById("recentlyAuthor"+i).innerHTML = "Author: " + data[i].author;
                document.getElementById("recentlyCreateTime"+i).innerHTML = "CreateTime: " + data[i].createtime.split('T')[0];
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// recentlyClick 은 초기페이지에서 가장많이 사용되는 아이템의 next, prev 버튼을 눌렀을때 실행하는 함수이다.
function topUsingClick(totalItemNum, buttonState) {
    // RecentlyTotalNum 가장많이 사용된 에셋의 전체 아이템 수
    let totalPageNum = Math.ceil(totalItemNum / 4); // 전체 페이지 수
    let clearItemNum = (totalPageNum * 4) - totalItemNum; // 마지막 페이지의 공백처리할 아이템 수
    let currentPageNum = parseInt(document.getElementById("topUsingPage").getAttribute('value'));

    if (buttonState=="next"){
        if (currentPageNum===totalPageNum){
            currentPageNum = 1;
        }else{
            currentPageNum++;
        }
    }else{
        if (currentPageNum===1){
            currentPageNum = totalPageNum;
        }else{
            currentPageNum--;
        }
    }
    
    document.getElementById("topUsingPage").innerHTML = currentPageNum + " / " + totalPageNum;
    if(clearItemNum!=0 && currentPageNum==totalPageNum){
        for(let i = 3; clearItemNum!=0; i--, clearItemNum--){
            document.getElementById("topUsingImageForm"+i).innerHTML = ""
            document.getElementById("topUsingTitle"+i).innerHTML = ""
            document.getElementById("topUsingAuthor"+i).innerHTML = ""
            document.getElementById("topUsingRate"+i).innerHTML = ""
        }
    }
    document.getElementById("topUsingPage").setAttribute('value', currentPageNum);
    $.ajax({
        url: `/api/topusingitem?usingpage=${currentPageNum}`,
        type: "get",
        dataType: "json",
        success: function(data) {
            let thumbnailwidth = document.getElementById("thumbnailwidth").value;
            let thumbnailheight = document.getElementById("thumbnailheight").value;
            let img = ""
            for (let i = 0; i < data.length; i++){
                let topUsingImageForm = document.getElementById("topUsingImageForm"+i)
                if (data[i].itemtype=="pdf"){
                    img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/assets/img/pdfthumbnail.svg">'
                }else if(data[i].itemtype=="hwp"){
                    img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/assets/img/hwpthumbnail.svg">'
                }else if(data[i].itemtype=="sound"){
                    img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/assets/img/soundthumbnail.svg">'
                }else if(data[i].itemtype=="hdri" || data[i].itemtype=="texture"){
                    if(data[i].status == "done"){
                        img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/mediadata?id=' + data[i].id + '&type=png">'
                    }else{
                        img = '<img width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" src="/assets/img/noimage.svg">'
                    }
                }else{
                    if(data[i].status == "done"){
                        img = '<video width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" controls poster="/mediadata?id=' + data[i].id + '&type=png">' +
                                '<source src="/mediadata?id=' + data[i].id + '&type=mp4" type="video/mp4">' +
                                '<source src="/mediadata?id=' + data[i].id + '&type=ogg" type="video/ogg">' +
                                'Your browser does not support the video tag.'+
                                '</video>'
                    }else{
                        img = '<video width="' + thumbnailwidth + '" height="'+ thumbnailheight +
                                '" controls poster="/assets/img/noimage.svg">' +
                                '<source src="/mediadata?id=' + data[i].id + '&type=mp4" type="video/mp4">' +
                                '<source src="/mediadata?id=' + data[i].id + '&type=ogg" type="video/ogg">' +
                                'Your browser does not support the video tag.'+
                                '</video>'
                    }
                    topUsingImageForm.innerHTML = img;
                }
                document.getElementById("topUsingTitle"+i).innerHTML = "Title: " + data[i].title;
                document.getElementById("topUsingAuthor"+i).innerHTML = "Author: " + data[i].author;
                document.getElementById("topUsingRate"+i).innerHTML = "topUsingRate: " + data[i].usingrate;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}
