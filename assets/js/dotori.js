
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

// setRmItemModal 은 아이템 삭제 버튼을 누르면 id값을 받아 modal창에 보여주는 함수이다.
function setRmItemModal(itemtype, itemId) {
    document.getElementById("modal-rmitem-itemtype").value = itemtype;
    document.getElementById("modal-rmitem-itemid").value = itemId;
}

// setRmItemModal 은 아이템 삭제 버튼을 누르면 id값을 받아 modal창에 보여주는 함수이다.
function setDetailViewModal(itemtype, itemid) {
    $.ajax({
        url: `/api/item?itemtype=${itemtype}&id=${itemid}`,
        type: "get",
        dataType: "json",
        success: function(response) {
            document.getElementById("modal-detailview-title").innerHTML = response["title"];
            document.getElementById("modal-detailview-itemid").innerHTML = itemid;
            document.getElementById("modal-detailview-author").innerHTML = response["author"];
            document.getElementById("modal-detailview-description").innerHTML = response["description"];
            let outputdatapath=response["outputdatapath"]
            let footerHtml=`
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-download-button" onclick="location.href='/download-item?itemtype=${itemtype}&id=${itemid}'">Download</a>
            <button type="button" class="btn btn-outline-darkmode" id="modal-detailview-copypath-button" onclick="copyButton('${outputdatapath}')">Copy Path</a>
            <button type="button" class="btn btn-outline-darkmode" data-dismiss="modal">Close</button>
            `
            document.getElementById("modal-detailview-footer").innerHTML = footerHtml
            if (itemtype == "footage") {
                document.getElementById("modal-detailview-download-button").type="hidden"        
            }
        },
        error: function(result) {
            alert(result);
        }
    });

}

// rmItemModal 은 삭제 modal창에서 Delete 버튼을 누르면 실행되는 아이템 삭제 함수이다. 
function rmItemModal(itemtype,itemId) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?itemtype=${itemtype}&id=${itemId}`,
        type: "delete",
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function() {
            alert("itemtype: "+itemtype+"\nid: "+itemId+"\n\n아이템 삭제를 성공했습니다."); 
            location.reload();
        },
        error: function(){
            alert("itemtype: "+itemtype+"\nid: "+itemId+"\n\n아이템 삭제를 실패했습니다.");  
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