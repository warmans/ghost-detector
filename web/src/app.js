import Chartist from "chartist";
import format from "date-fns/format";
import "./style.scss";

document.addEventListener('DOMContentLoaded', function () {

    init();

    let data = {
        // Our series array that contains series objects or in this case series data arrays
        series: [
            {
                name: 'series-1',
                data: []
            }
        ],
    };
    let config = {
        low: 0,
        showArea: true,
        axisX: {
            type: Chartist.FixedScaleAxis,
            divisor: 10,
            labelInterpolationFnc: function (value) {
                return format(value, "HH:mm:ss")
            }
        }
    };

    let chart = new Chartist.Line('.ct-chart', data, config);

    // start updating
    setInterval(() => {
        while (data.series[0].data.length > 100) {
            data.series[0].data.shift();
        }
        data.series[0].data.push({x: new Date(Date.now()), y: Math.random()});
        chart.update(data, config);
    }, 100)
});

function init() {
    let div = document.createElement('div');
    div.innerHTML = require("./app.tpl.html")({});
    document.body.appendChild(div);
}
