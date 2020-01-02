
function addAttribute() {
    let att = document.getElementById("attributes").innerHTML;
    let childnum = document.getElementById("attributes").childElementCount
    att = att + `
    <div class="row">
        <div class="col">
            <div class="form-group p-0 m-0">
                <small class="form-text text-muted">Key</small>
                <input type="text" class="form-control" placeholder="key" value="" name="key${childnum}">
            </div>
        </div>
        <div class="col">
            <div class="form-group p-0 m-0">
                <small class="form-text text-muted">Value</small>
                <input type="text" class="form-control" placeholder="value" value="" name="value${childnum}">
            </div>			
        </div>
    </div>
    `
    document.getElementById("attributes").innerHTML = att;
}