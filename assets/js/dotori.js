
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

// 버튼을 누르면 아이템에 대한 정보가 저장된다.
function editItem(){
    console.log("debug1")
    let author = document.getElementById("author");
    let description = document.getElementById("description")
    let tag = document.getElementById("tag")

    console.log("debug2")
    $.ajax({
        url:"/api/item",
        type:"post",
        data:{
            type: "maya",
            author: author,
        },
        dataType: "json",
        success: function(data){
            console.log("success")
            alert("success");
        },
        eorr: function(request, status, erroir){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}