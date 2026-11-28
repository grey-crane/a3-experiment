function generate_data_point(xMax, yMax, zMax){
    return {
        x: Math.floor(Math.random() * xMax),
        y: Math.floor(Math.random() * yMax),
        z: Math.floor(Math.random() * zMax),
    }
}

function generate_data_points(numPoints, xMax, yMax, zMax){
    const points = [];
    for (let i = 0; i < numPoints; i++){
        points.push(generate_data_point(xMax, yMax, zMax));
    }
    return points
}

console.log("d3 version", d3.version); // test that d3 loaded

// set the dimensions and margins of the graph
var margin = {top: 10, right: 30, bottom: 50, left: 60},
    width = 500 - margin.left - margin.right,
    height = 500 - margin.top - margin.bottom;

// Add an SVG - container
const svg = d3.select("#container").append("svg")
    .attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
    .append("g").attr("transform", "translate(" + margin.left + "," + margin.top + ")");

const data = generate_data_points(100, 400, 400, 100);

svg.selectAll("circle")
    .data(data).enter()
    .append("circle")
    .attr("cx", function(d) {return d.x+25})
    .attr("cy", function(d) {return d.y+25})
    .attr("r", function(d) {return 5 + d.z/10})
    .style("opacity", 0.8)