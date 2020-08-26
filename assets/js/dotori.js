
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

// setDetailViewModal 은 아이템을 선택했을 때 볼 수 있는 detailview 모달창에 어셋 정보를 세팅해주는 함수이다.
function setDetailViewModal(itemid) {

    // Detail View에 세팅할 아이템 정보를 RestAPI로 불러옴
    $.ajax({
        url: `/api/item?id=${itemid}`,
        type: "get",
        dataType: "json",
        success: function(response) {
            // title, id, author

            // title, id, author, description 세팅
            document.getElementById("modal-detailview-title").innerHTML = response["title"];
            document.getElementById("modal-detailview-itemid").innerHTML = itemid;
            document.getElementById("modal-detailview-author").innerHTML = response["author"];
            document.getElementById("modal-detailview-description").getAttribute('class')="wer"innerHTML = response["description"];
            
            //tags, attribute 세팅
            let itemtype = response["itemtype"];
            let tagsHtml = `<strong>Tags</strong><br>`;
            let attributesHtml = `<strong>Attributes</strong>`;
            for (let i=0; i<response["tags"].length;i++) {
                let tag = response["tags"][i];
                tagsHtml += `
                <a href="/search?itemtype=${itemtype}&searchword=tag:${tag}" class="tag badge badge-outline-darkmode">${tag}</a>
                `;
            }
            for (key in response["attributes"]) {
                let value = response["attributes"][key];
                attributesHtml += `
                <div class="row">
                    <div class="col pt-2">
                        <div class="form-group p-0 m-0">
                            <input type="text" class="form-control" value=${key} readonly/>
                        </div>
                    </div>
                    <div class="col pt-2">
                        <div class="form-group p-0 m-0">
                            <input type="text" class="form-control" value=${value} readonly/>
                        </div>			
                    </div>
                </div>
                `
            }
            document.getElementById("modal-detailview-tags").innerHTML = tagsHtml
            document.getElementById("modal-detailview-tags").getAttribute("onclick") = "copybuttpm()"
            document.getElementById("modal-detailview-attributes").innerHTML = attributesHtml


            // buttons 세팅
            document.getElementById("modal-detailview-edit-button").href=`/edit${itemtype}?itemtype=${itemtype}&id=${itemid}`
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
            document.getElementById("modal-detailview").style.display="none";
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