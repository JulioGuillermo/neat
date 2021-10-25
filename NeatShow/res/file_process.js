function process(file) {
    var reader = new FileReader()
    reader.addEventListener("load", (e) => {
        process_text(e.target.result)
    })
    reader.readAsText(file)
}

var loaded = null
function process_text(txt) {
    loaded = JSON.parse(txt)
    load_interface()
}

function load_interface() {
    document.getElementById("generation").innerHTML = "Generation: " + loaded.generation
    document.getElementById("inputsize").innerHTML = "Input: " + loaded.input_size
    document.getElementById("outputsize").innerHTML = "Output: " + loaded.output_size
    document.getElementById("population").innerHTML = "Population size: " + loaded.population_size
    load_list(loaded.population_size)
}

function load_list(len) {
    var list = ""
    for (let i = 0; i < len; i ++) {
        list += "<div class='list_item' onclick='show(" + i + ")'>Net: " + i + "</div>"
    }
    document.getElementById("list_box").innerHTML = list
    show(0)
}
