var nodeRadius=12;

var svg = d3.select("svg"),
width = +svg.attr("width"),
height = +svg.attr("height");
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
};

svg.append('svg:defs').append('svg:marker')
  .attr('id', 'end-arrow')
  .attr('viewBox', '0 -5 10 10')
  .attr('refX', 6)
  .attr('markerWidth', 3)
  .attr('markerHeight', 3)
  .attr('orient', 'auto')
  .append('svg:path')
  .attr('d', 'M0,-5L10,0L0,5')
  // .attr('fill', '#000')
;

var g = svg.append("g").attr("transform", "translate(" + width / 2 + "," + height / 2 + ")"),
link = g.append("g").attr("class","link").selectAll(".link"),
path = g.append("g").attr("class","link").selectAll(".link"),
node = g.append("g").attr("class", "node").selectAll("circle"),
nodet = g.append("g").attr("class","node").selectAll("text");

var nodes=[],links=[];

function linkDistance(l){
  switch(l.value){
    case "live":
    case "liveManager":
      return 60;
    case "IDOffice":
      return 100;
    case "dead":
      return 200;
    default:
      return 200;
  }
}
var simulation = d3.forceSimulation(nodes)
.force("charge", d3.forceManyBody().strength(-60))
.force("link", d3.forceLink(links)
  .id(function(d) { return d.id })
  .distance(function(d){return linkDistance(d)}))
// .force("x", d3.forceX())
// .force("y", d3.forceY())
.force("center", d3.forceCenter(0,0))
// .alphaTarget(1)
.on("tick", ticked);

function ticked() {
  node.attr("cx", function(d) { return d.x; })
      .attr("cy", function(d) { return d.y; });

  nodet.attr("x", function(d) { return d.x; })
      .attr("y", function(d) { return d.y; });

  link.attr("x1", function(d) { return d.source.x; })
      .attr("y1", function(d) { return d.source.y; })
      .attr("x2", function(d) { return d.target.x; })
      .attr("y2", function(d) { return d.target.y; });

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
}



var update = function() {
  d3.json("/heartbeat", function(error, dat) {
    if (error) throw error;

    nodes = getNodes(dat);
    links = getLinks(dat);

    var now = new Date();
    d3.select("#updateTime").html(now);

    restart();

    function liveNess(n){
      if (new Date().getTime() - new Date(n.T).getTime() > 2000){
        return "dead";
      } else {
        if (n.Rank.indexOf(".M.")>0){
          return "liveManager";
        }
        return "live";
      }
    }

    function linkType(source,status){
      if (source == "EgA.IDOffice" &&
        (status == "live" || status == "liveManager")) {
        return "IDOffice";
      }
      return status;
    }

    function strokeWidth(l) {
      switch(l.value){
        case "dead":
          return 0.3;
        case "IDOffice":
          return 0.3;
        default:
          return 2;
      }
    }

    function dispName(d){
      if (d.name !== undefined) {
        if (d.group == "dead") {
          return d.name+"#died:"+formatDateTime(d.t);
        } else {
          return d.name;
        }
      } else {
        return d.id;
      }
    }

    function formatDateTime(t) {
      return t.getFullYear() + "-" +
        d3.format("02d")(t.getMonth()) + "-" +
        d3.format("02d")(t.getDate()) + "_" +
        d3.format("02d")(t.getHours()) + ":" +
        d3.format("02d")(t.getMinutes()) +":" +
        d3.format("02d")(t.getSeconds());
    }

    function svcs(dat){
      var n=[];
      for (var key in dat) {
        n.push(dat[key]);
      }
      return n;
    }

    function boards(dat){
      var n=[];
      var b={};
      for (var key in dat){
        for (var tx in dat[key].Tx){
          b[dat[key].Tx[tx]]=1;
        }
        for (var rx in dat[key].Rx){
          b[dat[key].Rx[rx]]=1;
        }
      }
      for (var k in b){
        n.push(k);
      }
      return n;
    }

    function getNodes(dat){
      var s=svcs(dat);
      var b=boards(dat);
      var n=[];
      for (var i in b){
        // n.push({"id":b[i],"type":"board"});
        n.push({"id":b[i],"group":"board"});
      }
      for (var i in s){
        // n.push({"id":s[i].ID,"type":"svc","svc":s[i]})
        n.push({"id":s[i].ID,"name":s[i].Name,"group":liveNess(s[i]),
          "t":new Date(s[i].T),"rank":s[i].Rank});
      }
      return n;
    }

    function getLinks(dat){
      var l=[];
      var s=svcs(dat);
      for (var i in s){
        var n=s[i].ID;
        for (var tx in s[i].Tx){
          l.push({"source":n,"target":s[i].Tx[tx],
            "value":liveNess(s[i])});
        }
        for (var rx in s[i].Rx){
          l.push({"source":s[i].Rx[rx],"target":n,
            "value":linkType(s[i].Rx[rx],liveNess(s[i]))});
        }
      }
      return l;
    }


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

    function restart() {
      // Apply the general update pattern to the nodes.
      node = node.data(nodes, function(d) { return d.id;});
      node.exit().remove();
      node = node.enter()
        .append("circle").attr("class","node")
        .attr("fill", function(d) { return color(d.group); })
        .attr("r", nodeRadius)
        .call(d3.drag()
            .on("start", dragstarted)
            .on("drag", dragged)
            .on("end", dragended)).merge(node);

      nodet = nodet.data(nodes, function(d) { return d.id;});
      nodet.exit().remove();
      nodet = nodet.enter().append("text")
      // .attr("class","node")
      .text(function(d) { return dispName(d); })
      .call(d3.drag()
          .on("start", dragstarted)
          .on("drag", dragged)
          .on("end", dragended)).merge(nodet);

      // Apply the general update pattern to the links.
      // link = link.data(links, function(d) { return d.source.id + "-" + d.target.id; });
      // link.exit().remove();
      // link = link.enter().append("line")
      //   .attr("stroke-width", function(d){return strokeWidth(d)})
      //   .merge(link);
      // Apply the general update pattern to the links.
      path = path.data(links, function(d) { return d.source.id + "-" + d.target.id; });
      path.exit().remove();
      path = path.enter().append("path")
        .attr("stroke-width", function(d){return strokeWidth(d)})
        .style("marker-end","url(#end-arrow)")
        .merge(path);

      // Update and restart the simulation.
      simulation.nodes(nodes);
      simulation.force("link").links(links);
      simulation.alpha(1).restart();
    }

  });
};

// d3.interval loads data from the API endpoint
// and calls the update function.
update();
d3.interval(function() {
  update();
},30000);
