

function formatDateTime(t) {
  return t.getFullYear() + "-" +
    d3.format("02d")(t.getMonth()) + "-" +
    d3.format("02d")(t.getDate()) + "_" +
    d3.format("02d")(t.getHours()) + ":" +
    d3.format("02d")(t.getMinutes()) +":" +
    d3.format("02d")(t.getSeconds());
}

var svg = d3.select("svg"),
    width = +svg.attr("width"),
    height = +svg.attr("height");

// define arrow markers for graph links
svg.append('svg:defs').append('svg:marker')
    .attr('id', 'end-arrow')
    .attr('viewBox', '0 -5 10 10')
    .attr('refX', 6)
    .attr('markerWidth', 3)
    .attr('markerHeight', 3)
    .attr('orient', 'auto')
  .append('svg:path')
    .attr('d', 'M0,-5L10,0L0,5');
        // .attr('fill', '#000');

// var color = d3.scaleOrdinal(d3.schemeCategory10);
var color = function(group){
  switch(group){
    case "live":
      return d3.schemeCategory10[2];
      break;
    case "liveManager":
      return d3.schemeCategory10[1];
      break;
    case "board":
      return d3.schemeCategory10[9];
      break;
    default:
      return d3.schemeCategory10[3];
  }
}

let linkDistance = function(l){
  switch(l.value){
    case "live":
    case "liveManager":
      return 60;
    case "IDOffice":
      return 120;
    case "dead":
      return 300;
    default:
      return 300;
  }
}

var simulation = d3.forceSimulation()
    .force("link", d3.forceLink().id(function(d) { return d.id; })
      .distance(function(d){return linkDistance(d)}))
    .force("charge", d3.forceManyBody().strength(-75))
    .force("yPos", d3.forceY(height/2).strength(0.03)) // useful for landscape
    //.force("xPos", d3.forceX(width/2).strength(0.03)) // useful for portrait
    .force("center", d3.forceCenter(width / 2, height / 2));

var nodeRadius=12;
d3.json("/heartbeat", function(error, dat) {
  var nodes=[], links=[], boards={};
  let liveNess = function(n){
    let now = new Date();
    if (now.getTime() - (new Date(n.T).getTime()) > 2000){
      return "dead";
    } else {
      if (n.Rank.indexOf(".M.")>0){
        return "liveManager";
      }
      return "live";
    }
  }
  let linkType = function(source,status){
    if (source == "EgA.IDOffice" &&
      (status == "live" || status == "liveManager")) {
      return "IDOffice";
    }
    return status;
  }

  if (error) throw error;
  for (var key in dat) {
    let now = new Date();
    let status = liveNess(dat[key]);

    nodes.push({"id":dat[key].ID,"name":dat[key].Name,
      "group":status,"t":new Date(dat[key].T),
      "rank":dat[key].Rank});

    for (var tx in dat[key].Tx) {
      boards[dat[key].Tx[tx]]=1;
      links.push({"source": dat[key].ID,
        "target":dat[key].Tx[tx],"value":status});
    }

    for (var rx in dat[key].Rx) {
      boards[dat[key].Rx[rx]]=1;
      links.push({"target": dat[key].ID,
        "source":dat[key].Rx[rx],"value":linkType(dat[key].Rx[rx],status)});
    }
  }
  for (var key in boards){
    // var gp = Math.floor(Math.random()*19)+2
    // var gp = Math.abs(key.hashCode())%18+2
    var gp="board";
    nodes.push({"id":key,"group":gp});
    // console.log(key, gp)
  }
  var graph={"nodes":nodes,"links":links};

  let strokeWidth = function(l) {
    switch(l.value){
      case "dead":
        return 0.3;
      case "IDOffice":
        return 0.3;
      default:
        return 2;
    }
  }


  var nodec = svg.append("g")
      .attr("class", "nodes")
    .selectAll("circle")
    .data(graph.nodes)
    .enter().append("circle")
      .attr("r",nodeRadius)
      .attr("fill", function(d) { return color(d.group); })
      // .attr("fill", function(d) { return color(nodeColor(d)); })
      .call(d3.drag()
          .on("start", dragstarted)
          .on("drag", dragged)
          .on("end", dragended));


  /*
   Let Vector V = P2 - P1, where P1 is the start position Vector
   and P2 the end position vector.
   And let r be the node radius of the start and end nodes.
   Then P'1, the position vector of the start of vector P'1 P'2 is:
   P'1 = P1+r.v and P'2 = P2-r.v
   */
   var path = svg.append("g")
   .attr("class","links")
   .selectAll(".links")
   .data(graph.links)
   .enter().append("path")
   .attr("stroke-width", function(d){return strokeWidth(d)})
   .style("marker-end","url(#end-arrow)")
  //  .attr("d","M10,10L100,100") // d will be computed in the tick function.
   .attr("stroke","#000");

  // var link = svg.append("g")
  //     .attr("class", "links")
  //   .selectAll("line")
  //   .data(graph.links)
  //   .enter().append("line")
  //   .attr("stroke-width",0);
      // .attr("stroke-width", function(d){return strokeWidth(d)})
      // .style('marker-end', 'url(#end-arrow)' );


  var dispName = function(d){
    if (d.name !== undefined) {
      if (d.group == "dead") {
        return d.name+"#died:"+formatDateTime(d.t);
      } else {
        return d.name;
      }
    } else {
      return d.id;
    }
  };

  var node = svg.append("g")
      .attr("class", "nodes")
    .selectAll("text")
    .data(graph.nodes)
    .enter().append("text")
      .text(function(d){return dispName(d)})
      .attr("font-family","sans-serif")
      .attr("font-size","0.9em")
      .attr("fill","#444455")
      .call(d3.drag()
          .on("start", dragstarted)
          .on("drag", dragged)
          .on("end", dragended));

  node.exit().remove();
  nodec.exit().remove();
  // link.exit().remove();
  path.exit().remove();

  // node.append("title")
  //     .text(function(d) { return d.id; });

  simulation
      .nodes(graph.nodes)
      .on("tick", ticked);

  simulation.force("link")
      .links(graph.links);

  function ticked() {
    // link
    //     .attr("x1", function(d) { return d.source.x; })
    //     .attr("y1", function(d) { return d.source.y; })
    //     .attr("x2", function(d) { return d.target.x; })
    //     .attr("y2", function(d) { return d.target.y; });
    /*
     Let Vector V = P2 - P1, where P1 is the start position Vector
     and P2 the end position vector.
     And let r be the node radius of the start and end nodes.
     Then P'1, the position vector of the start of vector P'1 P'2 is:
     P'1 = P1+r.v and P'2 = P2-r.v
     */
    path
        .attr("d",function(d){
          Vx = (d.target.x-d.source.x);
          Vy = (d.target.y-d.source.y);
          Vmag = Math.sqrt(Vx*Vx + Vy*Vy);
          VunitX = Vx / Vmag;
          VunitY = Vy / Vmag;
          PsX = d.source.x+nodeRadius*VunitX;
          PsY = d.source.y+nodeRadius*VunitY;
          PeX = d.target.x-nodeRadius*VunitX;
          PeY = d.target.y-nodeRadius*VunitY;
          return "M"+PsX+","+PsY+"L"+PeX+","+PeY;
        });

    nodec
        .attr("cx", function(d) { return d.x; })
        .attr("cy", function(d) { return d.y; });

    node
        .attr("x", function(d) { return d.x; })
        .attr("y", function(d) { return d.y; });
  }
});

function dragstarted(d) {
  if (!d3.event.active) simulation.alphaTarget(0.3).restart();
  d.fx = d.x;
  d.fy = d.y;
}

function dragged(d) {
  d.fx = d3.event.x;
  d.fy = d3.event.y;
}

function dragended(d) {
  if (!d3.event.active) simulation.alphaTarget(0);
  d.fx = null;
  d.fy = null;
}
