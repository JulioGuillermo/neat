var index = 0
function show(i) {
    if (loaded != null) {
        if (i < 0) {
            i = loaded.population.length - 1
        } else if (i >= loaded.population.length) {
            i = 0
        }
        index = i
        var net = loaded.population[i]
        net = process_net(net)
        prerender_net(net)
        document.getElementById("net").innerHTML = "Network: " + index
    }
}

function process_net(net) {
    var neurons = []
    for (let i = 0; i < net.neurons.length; i++) {
        neurons.push(make_neuron(net.neurons[i]))
    }
    for (let i = 0; i < neurons.length; i++) {
        set_neuron_conections(neurons[i], neurons, net.neurons[i].connections)
    }
    return {
        input_size:     loaded.input_size,
        output_size:    loaded.output_size,
        neurons:        neurons,
    }
}

function set_neuron_conections(neuron, neurons, connections) {
    for (let i = 0; i < connections.length; i++) {
        neuron.connections[i] = neurons[connections[i]]
    }
}

function make_neuron(net) {
    return {
        x:          0,
        y:          0,
        name:       net.name,
        bias:       net.bias,
        weights:     net.weights,
        connections: [],
    }
}

function set_neuron_x_position(neuron, x, dx) {
    if (neuron.x < x) {
        neuron.x = x
    }
    for (let i = 0; i < neuron.connections.length; i ++) {
        set_neuron_x_position(neuron.connections[i], x + dx, dx)
    }
}

function prerender_net(net) {
    var max_bias = Math.abs(net.neurons[net.input_size].bias)
    var max_w = false
    for (let i = net.input_size + 1; i < net.neurons.length; i ++) {
        if (max_bias < Math.abs(net.neurons[i].bias)) {
            max_bias = Math.abs(net.neurons[i].bias)
        }
        for (let j = 0; j < net.neurons[i].weights.length; j ++) {
            if (!max_w || max_w < Math.abs(net.neurons[i].weights[j])) {
                max_w = Math.abs(net.neurons[i].weights[j])
            }
        }
    }
    for (let i = net.input_size; i < net.neurons.length; i ++) {
        net.neurons[i].bias /= max_bias
        for (let j = 0; j < net.neurons[i].weights.length; j ++) {
            net.neurons[i].weights[j] /= max_w
        }
    }
    var scale = parseFloat(document.getElementById("scale").value) / 100

    var margin = parseFloat(document.getElementById("margin").value) * scale
    var margin_left = margin * scale
    var margin_right = margin * scale
    var margin_top = margin * scale
    var maring_bottom = margin * scale

    var dx = parseFloat(document.getElementById("dx").value) * scale
    var dy = parseFloat(document.getElementById("dy").value) * scale

    var width = 0
    var height = 0

    // set neurons x pos
    for (let i = net.input_size; i < net.input_size + net.output_size; i ++) {
        set_neuron_x_position(net.neurons[i], 0, dx)
    }

    // align x pos of input neurons
    max_x = 0
    for (let i = 0; i < net.input_size; i ++) {
        if (max_x < net.neurons[i].x) {
            max_x = net.neurons[i].x
        }
    }
    for (let i = 0; i < net.input_size; i ++) {
        net.neurons[i].x = max_x
    }
    // Invert x pos and add margin
    for (let i = 0; i < net.neurons.length; i ++) {
        net.neurons[i].x = margin_left + max_x - net.neurons[i].x
    }

    // set width
    width = margin_left + margin_right + max_x

    // set neurons y pos
    var ypos = {}
    var max_y = 0
    var neuron
    for (let i = 0; i < net.neurons.length; i ++) {
        neuron = net.neurons[i]
        if (!ypos.hasOwnProperty(neuron.x)) {
            ypos[neuron.x] = 0
        } else {
            ypos[neuron.x] += dy
        }
        neuron.y = ypos[neuron.x]
        if (max_y < neuron.y) {
            max_y = neuron.y
        }
    }
    // center neurons y pos and add margin top
    for (let i in ypos) {
        ypos[i] = (max_y - ypos[i]) / 2 + margin_top
    }
    for (let i = 0; i < net.neurons.length; i ++) {
        neuron = net.neurons[i]
        neuron.y += ypos[neuron.x]
    }

    // set height
    height = margin_top + maring_bottom + max_y

    // render
    render_net(net, width, height, scale)
}

function hex(w) {
    var color = ""
    var map = "0123456789ABCDEF"
    var int_w = parseInt(w * 255)
    while (int_w > 0) {
        color = map.charAt(int_w % 16) + color
        int_w = parseInt(int_w / 16)
    }
    while (color.length < 2) {
        color = "0" + color
    }
    return color
}

function get_color(w) {
    var b = false
    if (w < 0) {
        b = true
        w = -w
    }
    var color = hex(w)
    if (b) {
        return "#0000" + color
    }
    return "#" + color + "0000"
}

function render_net(net, width, height, scale) {
    var neurons = net.neurons
    var inputs = net.input_size
    var outputs = net.output_size

    var useBezier = document.getElementById("useBezier").checked
    var neuronRadio = parseInt(document.getElementById("neuronRadio").value) * scale
    var outputRadio = parseInt(document.getElementById("outputRadio").value) * scale

    var canvas = document.getElementById("output")
    canvas.width = width
    canvas.height = height

    var c = canvas.getContext("2d")

    c.clearRect(0, 0, width, height)
    c.textBaseline = "middle"
    c.font = "bold " + parseInt(parseFloat(document.getElementById("fontSize").value) * scale) + "px verdana, sans-serif";
    c.lineWidth = document.getElementById("lineWidth").value * scale

    var n, sx, sy, con, ex, xy
    for (let i = 0; i < neurons.length; i ++) {
        n = neurons[i]
        sx = n.x - neuronRadio
        sy = n.y
        for (let j = 0; j < n.connections.length; j ++) {
            con = n.connections[j]
            ex = con.x + neuronRadio + outputRadio
            ey = con.y

            c.beginPath()
            c.strokeStyle = get_color(n.weights[j])
            c.moveTo(sx, sy)
            if (useBezier) c.bezierCurveTo(sx - Math.max((ex - sx) / 2, 100 * scale), sy, ex + Math.max((ex - sx) / 2, 100 * scale), ey, ex, ey)
            else c.lineTo(ex, ey)
            c.stroke()
            c.closePath()
        }
    }
    var outcontrol
    c.strokeStyle = "#111111"
    for (let i = 0; i < neurons.length; i ++) {
        n = neurons[i]
        sx = n.x
        sy = n.y

        c.beginPath()
        if (i < net.input_size) {
            c.fillStyle = "#FF00FF"
            outcontrol = true
        } else if (i < net.input_size + net.output_size) {
            c.fillStyle = "#FFAA00"
            outcontrol = false
        } else {
            outcontrol = true
            c.fillStyle = "#00FF00"
        }

        c.moveTo(sx, sy)
        c.arc(sx, sy, neuronRadio, 0, Math.PI * 2)
        c.fill()
        c.closePath()

        c.beginPath()
        if (i >= net.input_size) {
            c.fillStyle = get_color(n.bias)
        }
        if (outcontrol) {
            c.moveTo(sx + neuronRadio + outputRadio, sy)
            c.arc(sx + neuronRadio + outputRadio, sy, outputRadio, 0, Math.PI * 2)
        } else {
            c.moveTo(sx, sy)
            c.arc(sx, sy, outputRadio, 0, Math.PI * 2)
        }
        
        c.fill()
        c.fillStyle = "#000000"
        if (n.name != "") {
            if (n.name.charAt(0) == "I") {
                c.textAlign = "right"
                c.fillText(n.name, sx - neuronRadio - outputRadio, sy)
            } else {
                c.textAlign = "left"
                c.fillText(n.name, sx + neuronRadio + outputRadio, sy)
            }
        } else {
            c.textAlign = "center"
            c.fillText(i, sx, sy)
        }
        c.closePath()
    }
}
