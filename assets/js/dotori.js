
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
            document.getElementById("modal-detailview-title").innerHTML = response["title"] + `<button type="button" onclick="location.href='/edit${itemtype}?id=${itemid}'" class="btn btn-sm btn-outline-warning float-right" id="modal-detailview-edit-button">Edit</span>`;
            document.getElementById("modal-detailview-itemid").innerHTML = itemid;
            document.getElementById("modal-detailview-itemtype").innerHTML = itemtype;
            document.getElementById("modal-detailview-author").innerHTML = response["author"];
            document.getElementById("modal-detailview-description").innerHTML = response["description"].replace(/(\r\n|\n|\r)/g,"<br />");
            
            // Tags 세팅
            let tagsHtml = `<strong>Tags</strong><br>`;
            for (let i=0; i<response["tags"].length;i++) {
                let tag = response["tags"][i];
                tagsHtml += `
                <a href="/search?searchword=tag:${tag}" class="tag badge badge-outline-darkmode">${tag}</a>
                `;
            }
            document.getElementById("modal-detailview-tags").innerHTML = tagsHtml
            if (response["categories"] !== null) {
                document.getElementById("modal-detailview-categories").innerHTML = response["categories"].join(" > ")
            }
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
            <button type="button" class="btn btn-sm btn-outline-darkmode" id="modal-detailview-download-button" onclick="location.href='/download-item?id=${itemid}'">Download</button>
            `
            let footerHtmlForAdmin=`
            <button type="button" class="btn btn-sm btn-outline-darkmode" id="modal-detailview-download-button" onclick="location.href='/download-item?id=${itemid}'">Download</button>
            <button type="button" class="btn btn-sm btn-outline-danger" id="modal-detailview-download-button" onclick="location.href='/rename/${itemid}'">Rename</button>
            <button type="button" class="btn btn-sm btn-outline-danger" id="modal-detailview-delete-button" data-dismiss="modal" data-toggle="modal" data-target="#modal-rmitem">Delete</button>
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
            location.reload();
        }
    });

}

// rmItemModal 은 삭제 modal창에서 Delete 버튼을 누르면 실행되는 아이템 삭제 함수이다.
function rmItemModal(id) {
    fetch("/api/item?id="+ id, {
        method: 'DELETE',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        let elem = document.getElementById(id);
        elem.parentNode.removeChild(elem);
        tata.success('Remove', id + " Asset has been removed.", {position: 'tr',duration: 5000,onClose: null})
    })
    .catch((err) => {
        alert(err)
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
function toggleCheckboxs(){
    // 기준이 되는 체크박스의 상태값을 가지고 온다.
    let status = document.getElementById("toggle-checkbox").checked
    // 가져온 상태값을 기준으로 모든 체크박스의 상태를 설정한다.
    let checkboxes = document.querySelectorAll('*[name^="checkbox"]');
    for (i=0;i<checkboxes.length;i++) {
        checkboxes[i].checked = status
    }
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
function initPasword(userid) {
    let user = new Object()
    user.id = userid
    fetch('/api/initpassword', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(user),
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        tata.success('InitPassword', data.id+" 사용자의 패스워드가 초기화 되었습니다", {position: 'tr', duration: 5000, onClose: null})
    })
    .catch((err) => {
        alert(err)
    });
}

function changeUserAccessLevel(userid, level) {
    let user = new Object()
    user.id = userid
    user.accesslevel = level
    fetch('/api/user/accesslevel', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(user),
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        tata.success('Change AccessLevel', "AccessLevel of " + data.id + " user have been changed.", {position: 'tr', duration: 5000, onClose: null})
    })
    .catch((err) => {
        alert(err)
    });
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}


function SetUserAutoplay() {
    fetch('/api/user/autoplay?value='+document.getElementById("autoplay").checked, {
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

function SetUserNewsNum() {
    fetch('/api/user/newsnum?value='+document.getElementById("newsnum").value, {
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

function SetUserTopNum() {
    fetch('/api/user/topnum?value='+document.getElementById("topnum").value, {
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
    tata.success('Copy Clipboard', "Data path copied!", {
        position: 'tr',
        duration: 1000,
        onClose: null,
    })
}

function SelectAll() {
    let checkboxs = document.querySelectorAll('.select-item')    
    for (let i = 0; i < checkboxs.length; i += 1) {        
        checkboxs[i].checked = true
    }
}

function SelectNone() {
    let checkboxs = document.querySelectorAll('.select-item')    
    for (let i = 0; i < checkboxs.length; i += 1) {        
        checkboxs[i].checked = false
    }
}

// copyPath 함수는 경로를 받아서, 클립보드로 복사하는 기능이다.
function copyPath(path) {
    let admin = new Object()
    fetch('/api/adminsetting', {
        method: 'GET',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        admin = data;
        
        if (navigator.userAgent.indexOf("Win") != -1) { // windows 경우
            path = path.replace(admin.rootpath,"") // linux의 루트경로를 지운다.
            path = admin.windowsuncprefix.replace(/\//g, "\\") + path.replace(/\//g, "\\")
        }
        copyClipboard(path)
        
    })
    .catch((err) => {
        alert(err)
    });
    
}


/** copyNukePath 함수는 ID값을 받아서, 클립보드로 복사하는 기능이다.*/
function copyNukePath(id) {
    let admin = new Object()
    let paths = []
    let getAdmin = new Promise( function( resolve, reject ) {
        fetch('/api/adminsetting', {
            method: 'GET',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            admin = data;
            resolve();
        })
        .catch((err) => {
            reject(err)
        });
    });

    let getNukePaths = new Promise( function( resolve, reject ) {
        fetch('/api/nukepath/'+id, {
            method: 'GET',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
        })
        .then((response) => {
            if (!response.ok) {
                throw Error(response.statusText + " - " + response.url);
            }
            return response.json()
        })
        .then((data) => {
            paths = data.nukepath
            resolve();
        })
        .catch((err) => {
            reject(err)
        });
    });
    Promise.all( [ getAdmin, getNukePaths ] )
    .then(function () {
        // 모든 RestAPI가 성공하면 아래 항목을 실행한다.
        if (paths.length != 1) {
            tata.error('Copy Clipboard', "There are several .nk files in the data folder.", {
                position: 'tr',
                duration: 1000,
                onClose: null,
            })
            return
        }
        let path = paths.pop()
        if (navigator.userAgent.indexOf("Win") != -1) { // windows 경우
            path = admin.windowsuncprefix + path.replace(/\//g, "\\")
        }
        copyClipboard(path)
    })
    .catch(function ( reason ) {
        console.log( reason ); // rejected되었기때문에 여기가 출력
    });    
}


function rvlink(path) {
    let admin = new Object()
    fetch('/api/adminsetting', {
        method: 'GET',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        admin = data;
        if (navigator.userAgent.indexOf("Win") != -1) { // windows 경우
            path = path.replace(admin.rootpath,"") // linux의 루트경로를 지운다.
            path = admin.windowsuncprefix.replace(/\//g, "\\") + path.replace(/\//g, "\\")
        }
        let obj = document.createElement("a");   // input요소를 만듬
        obj.href = "rvlink://" + path
        document.body.appendChild(obj);
        obj.click()
        document.body.removeChild(obj);
    })
    .catch((err) => {
        alert(err)
    });
}

function EditTags() {
    let tags = new Object()
    tags.tags = string2array(document.getElementById("modal-edittags-tags").value)
    if (tags.tags.length === 0) {
        tata.error('Error', "태그를 입력해주세요.", {position: 'tr',duration: 5000,onClose: null})
        return
    }
    // 선택한 ID를 출력한다.
    let idcheckboxs = document.querySelectorAll("[id^='idcheckbox-']")
    let ids = []
    for (let i = 0; i < idcheckboxs.length; i += 1) {
        if (!idcheckboxs[i].checked) {
            continue
        }
        ids.push(idcheckboxs[i].value)
    }
    if (ids.length === 0) {
        tata.error('Error', "Item을 선택해주세요.", {position: 'tr',duration: 5000,onClose: null})
        return
    }
    for (let i = 0; i < ids.length; i += 1) {
        let id = ids[i]
        fetch('/api/tags/'+id, {
            method: 'PUT',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
            body: JSON.stringify(tags),
        })
        .then((response) => {
            if (!response.ok) {
                response.text().then(function (text) {
                    tata.error('Error', text, {position: 'tr',duration: 5000,onClose: null})
                    return
                });
            }
            return response.json()
        })
        .then((obj) => {
            // 내부를 비운다.
            let e = document.getElementById("tags-"+id)
            e.innerHTML = ""
            // UX를 업데이트 한다.
            for (let t = 0; t < obj.tags.length; t += 1) {
                let html = `<a href="/search?searchword=tag:${obj.tags[t]}" class="tag badge badge-outline-darkmode">${obj.tags[t]}</a>`
                e.innerHTML = html + e.innerHTML
            }            
        })
        .catch((err) => {
            console.log(err)
        });
    }
    tata.success('Edit', "태그가 편집되었습니다.", {position: 'tr',duration: 5000,onClose: null})
}

function SetItemCategory() {
    let rootCategory = document.getElementById("rootcategory").options[document.getElementById("rootcategory").selectedIndex].text
    let subCategory = document.getElementById("subcategory").options[document.getElementById("subcategory").selectedIndex].text
    
    // 선택한 ID를 출력한다.
    let idcheckboxs = document.querySelectorAll("[id^='idcheckbox-']")
    let ids = []
    for (let i = 0; i < idcheckboxs.length; i += 1) {
        if (!idcheckboxs[i].checked) {
            continue
        }
        ids.push(idcheckboxs[i].value)
    }
    if (ids.length === 0) {
        tata.error('Error', "Please select assets", {position: 'tr',duration: 5000,onClose: null})
        return
    }
    for (let i = 0; i < ids.length; i += 1) {
        let id = ids[i]
        fetch('/api/item?id='+id, {
            method: 'GET',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
        })
        .then((response) => {
            if (!response.ok) {
                response.text().then(function (text) {
                    tata.error('Error', text, {position: 'tr',duration: 5000,onClose: null})
                    return
                });
            }
            return response.json()
        })
        .then((data) => {
            console.log(data)
            console.log(rootCategory, subCategory)
            data.categories = [rootCategory, subCategory]
            // data 업데이트
            fetch('/api/item/'+data.id, {
                method: 'PUT',
                headers: {
                    "Authorization": "Basic "+ document.getElementById("token").value,
                },
                body: JSON.stringify(data),
            })
            .then((response) => {
                if (!response.ok) {
                    response.text().then(function (text) {
                        tata.error('Error', text, {position: 'tr',duration: 5000,onClose: null})
                        return
                    });
                }
                return response.json()
            })
            .then((data) => {
                console.log(data)
                
            })
            .catch((err) => {
                console.log(err)
            });

        })
        .catch((err) => {
            console.log(err)
        });
    }
    tata.success('Success', ids.length + " Items<br>Set Categories "+rootCategory+" > "+subCategory, {position: 'tr',duration: 5000,onClose: null})
}

  
function CopyPaths() {
    GetOutputDataPaths().then(function(data) { // Promise 타입은 then을 이용해서 값을 가지고 와야한다.
        document.getElementById("modal-copypaths-text").value = data.join("\n")
    });    
}

async function GetOutputDataPaths() {
    let idcheckboxs = document.querySelectorAll("[id^='idcheckbox-']")
    let ids = []
    for (let i = 0; i < idcheckboxs.length; i += 1) {
        if (!idcheckboxs[i].checked) {
            continue
        }
        ids.push(idcheckboxs[i].value)
    }
    if (ids.length === 0) {
        tata.error('Error', "Item을 선택해주세요.", {position: 'tr',duration: 5000,onClose: null})
        return
    }
    let paths = []
    for (let i = 0; i < ids.length; i += 1) {
        let id = ids[i]
        await fetch('/api/item?id='+id, {
            method: 'GET',
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value,
            },
        })
        .then((response) => {
            if (!response.ok) {
                response.text().then(function (text) {
                    tata.error('Error', text, {position: 'tr',duration: 5000,onClose: null})
                    return
                });
            }
            return response.json()
        })
        .then((obj) => {
            paths.push(obj.outputdatapath)
        })
        .catch((err) => {
            console.log(err)
        });
    }
    return paths
}

function string2array(str) {
    let newArr = [];
    if (str === "") {
        return newArr
    }
    let arr = str.split(",");
    for (let i = 0; i < arr.length; i += 1) {
        newArr.push(arr[i].trim())
    }
    return newArr;
}

function AddCategory() {
    let category = new Object()
    category.name = document.getElementById("modal-addcategory-name").value
    category.parentid = document.getElementById("modal-addcategory-parentid").value
    fetch('/api/category', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(category),
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        let body = `<div id="category-${data.id}" class="finger category border border-dark p-2 ps-3 m-1 d-block align-items-center" onclick="selectCategory('${data.id}')">${data.name}
                    <img src="/assets/img/delete.svg" class="mt-1 icon finger" onclick="setRmCategoryID('${data.id}')" data-bs-toggle="modal" data-bs-target="#modal-rmcategory">
                    </div>`
        if (data.parentid == "") {
            document.getElementById("listofrootcategory").innerHTML +=  body
        } else {
            document.getElementById("listofsubcategory").innerHTML +=  body
        }
        
        tata.success('Add', "A category has been added.", {position: 'tr',duration: 5000,onClose: null})
    })
    .catch((err) => {
        alert(err)
    });    
}

function BackupDB() {
    let today = new Date();
    let yyyy = today.getFullYear();
    let mm = String(today.getMonth() + 1).padStart(2, '0'); //January is 0!
    let dd = String(today.getDate()).padStart(2, '0');
    let obj = new Object()
    obj.date = yyyy+mm+dd
    fetch('/api/dbbackup', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(obj),
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        tata.success('Backup', "Backup finished", {position: 'tr',duration: 5000,onClose: null})
    })
    .catch((err) => {
        alert(err)
    });
}


function setRmCategoryID(id) {
    document.getElementById("modal-rmcategory-id").value = id
}

function selectCategory(id) {
    // 기존에 선택된 요소를 삭제한다.
    let selectlist = document.getElementsByClassName("border-warning")
    for (let i = 0; i < selectlist.length; i+=1) {
        let e = selectlist[i]
        e.classList.remove("border-warning")
        e.classList.add("border-dark")
    }
    // 클릭한 것만 테두리를 표시한다.
    let e = document.getElementById("category-"+id)
    e.classList.remove("border-dark");
    e.classList.add("border-warning");

    // 선택된 ID로 Parentid를 설정한다.
    if (e.parentNode.id == "listofrootcategory") {
        document.getElementById("modal-addcategory-parentid").value = id
    } else {
        document.getElementById("modal-addcategory-parentid").value = ""
    }
    
    if (e.parentNode.id != "listofrootcategory") {
        return
    }
    // 선택된 자식을 서브 카테고리에 띄운다.
    fetch('/api/subcategories/'+id, {
        method: 'GET',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        console.log(data)
        if (data == null) {
            // 값이 없다면 subcategory를 비운다.
            document.getElementById("listofsubcategory").innerHTML = ""
            return
        }
        document.getElementById("listofsubcategory").innerHTML = ""
        // subcategory를 그린다.
        for (let i = 0; i < data.length; i+=1) {
            let body = `<div id="category-${data[i].id}" class="finger category border border-dark p-2 ps-3 m-1 d-block align-items-center" onclick="selectCategory('${data[i].id}')">${data[i].name}
                    <img src="/assets/img/delete.svg" class="mt-1 icon finger" onclick="setRmCategoryID('${data[i].id}')" data-bs-toggle="modal" data-bs-target="#modal-rmcategory">
                    </div>`
            document.getElementById("listofsubcategory").innerHTML +=  body
        }
    })
    .catch((err) => {
        alert(err)
    });
}

function RmCategory() {    
    fetch('/api/category/'+document.getElementById("modal-rmcategory-id").value, {
        method: 'DELETE',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        let elem = document.getElementById("category-" + data.id);
        let elemParentID = elem.parentNode.id
        elem.parentNode.removeChild(elem);
        if (elemParentID != "listofrootcategory") { // sub카테고리를 삭제할 때
            tata.success('Remove', "A category has been removed.", {position: 'tr',duration: 5000,onClose: null})
        } else {
            location.reload(); // root카테고리가 삭제될 때는 페이지를 리로드한다.
        }
    })
    .catch((err) => {
        alert(err)
    });    
}

function changeRootCategory() {
    document.getElementById("rootcategory-string").value = document.getElementById("rootcategory").options[document.getElementById("rootcategory").selectedIndex].text
    document.getElementById("subcategory-string").value = ""
    let id = document.getElementById("rootcategory").value
    if (id == "") {
        return
    }
    fetch('/api/subcategories/'+id, {
        method: 'GET',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        if (data == null) {
            // 값이 없다면 subcategory를 비운다.
            document.getElementById("subcategory").value = ""
            return
        }
        let select = document.getElementById("subcategory")

        // 기존 subcategory를 비운다.
        select.innerHTML = `<option value=""></option>`

        // subcategory를 그린다.
        for (let i = 0; i < data.length; i+=1) {
            let opt = document.createElement('option');
            opt.value = data[i].id;
            opt.innerHTML = data[i].name;
            select.appendChild(opt);
        }
    })
    .catch((err) => {
        console.log(err)
    });
}

function changeSearchboxRootCategory() {
    let id = document.getElementById("searchbox-rootcategory-id").value
    if (id == "") {
        return
    }
    fetch('/api/subcategories/'+id, {
        method: 'GET',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        if (data == null) {
            // 값이 없다면 subcategory를 비운다.
            document.getElementById("searchbox-subcategory-id").value = ""
            return
        }
        let select = document.getElementById("searchbox-subcategory-id")

        // 기존 subcategory를 비운다.
        select.innerHTML = `<option value="">Sub Category</option>`

        // subcategory를 그린다.
        for (let i = 0; i < data.length; i+=1) {
            let opt = document.createElement('option');
            opt.value = data[i].id;
            opt.innerHTML = data[i].name;
            select.appendChild(opt);
        }
    })
    .catch((err) => {
        console.log(err)
    });
}

function changeSubCategory() {
    document.getElementById("subcategory-string").value = document.getElementById("listofsubcategory").options[document.getElementById("listofsubcategory").selectedIndex].text
}