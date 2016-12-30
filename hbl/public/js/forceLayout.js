var svg = d3.select("svg"),
width = +svg.attr("width"),
height = +svg.attr("height"),
color = d3.scaleOrdinal(d3.schemeCategory20);

var g = svg.append("g").attr("transform", "translate(" + width / 2 + "," + height / 2 + ")"),
link = g.append("g").attr("stroke", "#000").attr("stroke-width", 1.5).selectAll(".link"),
node = g.append("g").attr("stroke", "#fff").attr("stroke-width", 1.5).selectAll(".node");

var nodes=[],links=[];

var simulation = d3.forceSimulation(nodes)
.force("charge", d3.forceManyBody().strength(-1000))
// .force("link", d3.forceLink(links).distance(200))
.force("link", d3.forceLink(links).id(function(d) { return d.id }).distance(200))
.force("x", d3.forceX())
.force("y", d3.forceY())
.alphaTarget(1)
.on("tick", ticked);

function ticked() {
  node.attr("cx", function(d) { return d.x; })
      .attr("cy", function(d) { return d.y; });

  link.attr("x1", function(d) { return d.source.x; })
      .attr("y1", function(d) { return d.source.y; })
      .attr("x2", function(d) { return d.target.x; })
      .attr("y2", function(d) { return d.target.y; });
}

// d3.interval loads data from the API endpoint
// and calls the update function.
d3.interval(function() {
  d3.json("/heartbeat", function(error, dat) {
    if (error) throw error;

    nodes = getNodes(dat);
    links = getLinks(dat);

    d3.select("#dat").html(JSON.stringify(links) + "<p>" +
      JSON.stringify(nodes)
    );

    restart();

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
        n.push({"id":b[i]});
      }
      for (var i in s){
        // n.push({"id":s[i].ID,"type":"svc","svc":s[i]})
        n.push({"id":s[i].ID})
      }
      return n;
    }

    function getLinks(dat){
      var l=[];
      var s=svcs(dat);
      for (var i in s){
        var n=s[i].ID;
        for (var tx in s[i].Tx){
          l.push({"source":n,"target":s[i].Tx[tx]});
        }
        for (var rx in s[i].Rx){
          l.push({"source":s[i].Rx[rx],"target":n});
        }
      }
      return l;
    }



    function restart() {
      // Apply the general update pattern to the nodes.
      node = node.data(nodes, function(d) { return d.id;});
      node.exit().remove();
      node = node.enter().append("circle").attr("class","node").attr("fill", function(d) { return color(d.id); }).attr("r", 8).merge(node);

      // Apply the general update pattern to the links.
      link = link.data(links, function(d) { return d.source.id + "-" + d.target.id; });
      link.exit().remove();
      link = link.enter().append("line").attr("class","link").merge(link);

      // Update and restart the simulation.
      simulation.nodes(nodes);
      simulation.force("link").links(links);
      simulation.alpha(1).restart();
    }

  });
},30);
