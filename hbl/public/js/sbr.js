'use strict';
// Identifies starting with Uppercase letter are in the global scope.
// let => block scope
// var => function scope
let UpdateInterval=1000; //ms
let IDOffice = "EgA.IDOffice";

window.onload = function() {
  let svg = d3.select("svg"),
    width = +svg.attr("width"),
    height = +svg.attr("height");

  let brdT = d3.select("#boards").selectAll(".brdT");
  let brdG = svg // svg group containing boards
    .append("g")
      .attr("transform","translate(220,40)");
  brdG.append("text")
    .classed("brdHdr",true)
      .attr("x",0).attr("y",-20).attr("opacity",0.3)
      .attr("text-anchor","end")
      .text("Boards");
  let brd = brdG
    .selectAll(".brd");

  let svcT = d3.select("#services").selectAll(".svcT");
  let svcG = svg // svg group containing senders
    .append("g")
    .attr("transform","translate(350,40)")
  svcG
    .append("text")
    .classed("svcHdr",true)
      .attr("x",0).attr("y",-20).attr("opacity",0.3)
        .text("Services");
  let svc = svcG
    .selectAll(".svc");

  let deadSvcG = svg // svg group containing dead services
    .append("g")
      .attr("transform","translate(650,40)");
  deadSvcG.append("text")
    .classed("dsvcHdr",true)
      .attr("x",0).attr("y",-20).attr("opacity",0.3)
      .text("R.I.P");
  let deadSvc = deadSvcG
      .selectAll(".dsvc");

  let sLink = svg // svg group containing svc send to board links
    .append("g")
      .attr("transform","translate(345,33)").selectAll(".sLink");
  let rLink = svg // svg group containing svc receive from board links
    .append("g")
      .attr("transform","translate(345,37)").selectAll(".rLink");

  // update get data from the API endpoint
  function update() {
    d3.json("/heartbeat", function(error, dat) {
      if (error) throw error;

      var now=new Date();
      d3.select("#updtTim").html(now);
      function svcs(){
        let n=[];
        for (let key in dat) {
          let s=dat[key];
          let live = now.getTime() - new Date(s.T) < 2000;
          if(s.Tx != null && s.Tx.length>0)var tx=true;
          if(s.Rx != null && s.Rx.length>0)var rx=true;
          n.push({"svc":s,"live":live,"tx":tx,"rx":rx});
        }
        return n;
      }

      function boards(){
        let n=[];
        let b={};
        for (let key in dat){
          for (let tx in dat[key].Tx){
            b[dat[key].Tx[tx]]=1;
          }
          for (let rx in dat[key].Rx){
            b[dat[key].Rx[rx]]=1;
          }
        }
        for (let k in b){
          n.push(k);
        }
        return n;
      }

      // links returns a map of services to boards.
      // Eg. l[0]={"s":[0,2],"r":[]}
      // means service 0 sends to boards 0 and 2
      //  and receives from none.
      function links() {
        let l=[];
        let s=svcs().filter(function(s){return s["live"]});
        let b=boards();
        for (let i in s){
          let name=s[i]["svc"].Name;
          let tx = s[i]["svc"].Tx;
          let  snd = d3.set();
          let rx = s[i]["svc"].Rx;
          let  rcv = d3.set();
          for (var j in tx) {
            snd.add(d3.scan(b,function(x){return -x.indexOf(tx[j])}));
            // console.log(name,i,"->",tx[j],d3.scan(b,function(x){return -x.indexOf(tx[j])}))
          }
          for (var k in rx) {
            if (rx[k] != IDOffice){ // ignore IDOffice receive link
              rcv.add(d3.scan(b,function(x){return -x.indexOf(rx[k])}));
            }
          }
          l.push({"s":snd,"r":rcv});
          // console.log(snd,rcv);
        }
        return l;
      }

      function svcView(s){
        return s["svc"].Name + " "+ s["svc"].T+":"+s["live"]+":"+s["tx"]+":"+s["rx"];
      }
      function dispSvcs(){
        let s=svcs();

        function dsHdr(){
          if(dead.length>0){
              deadSvcG.selectAll(".dsvcHdr").attr("opacity",1);
          } else {
            deadSvcG.selectAll(".dsvcHdr").attr("opacity",0.3);
          }
        }
        function sHdr(){
          if(lv.length>0){
            svcG.selectAll(".svcHdr").attr("opacity",1);
            d3.select("#svcHdr").style("opacity",1);
          } else {
            svcG.selectAll(".svcHdr").attr("opacity",0.3);
            d3.select("#svcHdr").style("opacity",0.3);
          }
        }

        let lv=s.filter(function(d){return d["live"]});
        svcT = svcT.data(lv);
        sHdr();
        svcT = svcT
          .classed("new",false)
          .classed("updated",true)
          .classed("live",function(d){return d["live"]}) // if not "live" set class "dead"
          .classed("manager",function(d){return d["svc"].Rank.indexOf(".M.")>0})
          .classed("dead",function(d){return !d["live"]});
        svcT.exit().html("");
        svcT = svcT.enter()
        .append("div").classed("svcT new",true)
        .merge(svcT)
          .html(function(d){return d["svc"].Name+"#"+d["svc"].Rank});

        // svg live services
        svc = svc.data(lv);
        svc = svc
          .classed("new",false)
          .classed("updated",true)
          .classed("live",function(d){return d["live"]}) // if not "live" set class "dead"
          .classed("manager",function(d){return d["svc"].Rank.indexOf(".M.")>0})
          .classed("dead",function(d){return !d["live"]});
        svc.exit().remove();
        svc = svc.enter()
          .append("text")
          .classed("svc new",true)
          .attr("x",0)
          .attr("dy",function(d,i){return i*1.2+"em"})
        .merge(svc)
          .text(function(d,i){return i+": "+d["svc"].Name+
            "#"+d["svc"].Rank});

        // svg dead services
        let dead=s.filter(function(d){return !d["live"]});

        dsHdr();

        deadSvc = deadSvc.data(dead);
        deadSvc = deadSvc
          .classed("new",false)
          .classed("updated",true)
          .classed("live",function(d){return d["live"]}) // if not "live" set class "dead"
          .classed("dead",function(d){return !d["live"]});
        deadSvc.exit().remove();
        deadSvc = deadSvc.enter()
          .append("text")
          .classed("dsvc new",true)
          .attr("x",0)
          .attr("dy",function(d,i){return i*1.2+"em"})
        .merge(deadSvc)
          .text(function(d,i){return i+": "+d["svc"].Name+
            ": " + d3.timeFormat("%a, %b %d, %H:%M:%S %Z")(new Date(d["svc"].T))
          });
      }

      function dispBoards(){

        function brdHdr(){
          // opacity attr for svg and style for h2 css
          if(b.length>0){
            d3.selectAll(".brdHdr").attr("opacity",1); // svg
            d3.select("#brdHdr").style("opacity",1); // html
          } else {
            d3.selectAll(".brdHdr").attr("opacity",0.3);
            d3.select("#brdHdr").style("opacity",0.3);
          }
        }

        let b=boards();
        brdHdr(); // display board header

        // html
        brdT = brdT.data(b);
        brdT = brdT
          .classed("new",false)
          .classed("updated",true);
        brdT = brdT.enter()
          .append("div").classed("brdT new",true)
          .merge(brdT)
            .html(function(d){return d});


        // svg
        // bind data b to selection brd (global variable) and update it
        brd = brd.data(b);
        // update existing elements by changing "new" class to "updated"
        brd = brd
          .classed("new",false) // existing elements no longer "new"
          .classed("updated",true); // instead they are now "updated"
        // remove exiting (outgoing) boards
        brd.exit().remove();
        // mark new entries, merge with old entries and perform operation on both
        brd = brd.enter()
          .append("text") // indent 2 if new selection, 4 if not.
            .classed("brd new",true) // no new sel. just an attr.
            .attr("x","0") // start at group origin x coord.
            .attr("dy",function(d,i){return i*1.2+"em"}) // move down one line space for new entries
            .attr("text-anchor","end")
          .merge(brd) // merge returns new selection: merge of enter and existing.
            .text(function(d,i){return i+": "+d}); // write text to both.
      }

      function sLinks(){
        let l=links();
        let s=[];
        for (let i in l){
          let bis=l[i]["s"].values(); // board indexes
          for (let j in bis){
            s.push([i,bis[j]]);
          }
        }
        return s
      }
      function rLinks(){
        let l=links();
        let r=[];
        for (let i in l){
          let bis=l[i]["r"].values(); // board indexes
          for (let j in bis){
            r.push([i,bis[j]]);
          }
        }
        return r
      }
      // dispSLinks displays Send links.
      function dispSLinks() {
        // console.log(l[2]["s"].values(),l[2]["r"].values());
        let l=sLinks(); // array of [svcIdx,boardIdx]
        sLink = sLink.data(l);
        sLink = sLink
            .classed("new",false)
            .classed("updated",true);

        sLink.exit().remove();

        sLink = sLink.enter()
          .append("line") // add SVG line
            .classed("sLink new",true)
            .merge(sLink)
            .attr("x1", 0)
            .attr("y1",function(d){//console.log(d);
              return d[0]*19})

            .attr("x2",-120)
            .attr("y2",function(d){return d[1]*19});
      }
      // dispRLinks displays Receive links.
      function dispRLinks() {
        let l=rLinks(); // array of [svcIdx,boardIdx]
        rLink = rLink.data(l);
        rLink = rLink
            .classed("new",false)
            .classed("updated",true);

        rLink.exit().remove();

        rLink = rLink.enter()
          .append("line") // add SVG line
            .classed("rLink new",true)
            .merge(rLink)
            .attr("x1", 0)
            .attr("y1",function(d){return d[0]*19})

            .attr("x2",-120)
            .attr("y2",function(d){return d[1]*19});
      }

      // ------------------------------------------------------------------

      dispSvcs();
      dispBoards();
      dispSLinks();
      dispRLinks();
    });
  }

  // d3.interval loads data from the API endpoint
  // and calls the update function.
  update();
  d3.interval(function() {
    update();
  },UpdateInterval);
}
