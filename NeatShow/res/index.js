function init() {
    var json_input = document.getElementById("json_input")
    json_input.addEventListener("change", (e) => {
        process(e.target.files[0])
        json_input.value = null
    })
}

window.addEventListener("load", init)
