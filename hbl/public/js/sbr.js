'use strict';
// Identifies starting with Uppercase letter are in the global scope.
// let => block scope
// var => function scope
let UpdateInterval=1000; //ms

window.onload = function() {
  let svg = d3.select("svg"),
    width = +svg.attr("width"),
    height = +svg.attr("height");

  let brdT = d3.select("#boards").selectAll(".brdT");
  let brdG = svg // svg group containing boards
    .append("g")
      .attr("transform","translate(200,40)");
  brdG.append("text")
    .classed("brdHdr",true)
      .attr("x",0).attr("y",-20).attr("opacity",0.3)
      .attr("text-anchor","end")
      .text("Boards");
  let brd = brdG
    .selectAll(".brd");

  let svcT = d3.select("#services").selectAll(".svcT");
  let sndrG = svg // svg group containing senders
    .append("g")
    .attr("transform","translate(350,40)")
  sndrG
    .append("text")
    .classed("svcHdr",true)
      .attr("x",0).attr("y",-20).attr("opacity",0.3)
        .text("Services");
  let sndr = sndrG
    .selectAll(".sndr");

  let deadSvcG = svg // svg group containing dead services
    .append("g")
      .attr("transform","translate(650,40)");
  deadSvcG.append("text")
    .classed("dsvcHdr",true)
      .attr("x",0).attr("y",-20).attr("opacity",0.3)
      .text("R.I.P");
  // let dHdrSet=false;
  let deadSvc = deadSvcG
      .selectAll(".dsvc");
  // svg
  //   .append("text").classed("dsvcHdr",true)
  //     .attr("x",300).attr("y",80)
  //     .text("R.I.P");

  // update get data from the API endpoint
  function update() {
    d3.json("/heartbeat", function(error, dat) {
      if (error) throw error;

      function svcs(){
        let now=new Date();
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

      function svcView(s){
        return s["svc"].Name + " "+ s["svc"].T+":"+s["live"]+":"+s["tx"]+":"+s["rx"];
      }
      function dispSvcs(){
        // txt
        d3.select("#services-live").selectAll("div").remove();
        d3.select("#services-dead").selectAll("div").remove();
        let s=svcs();
        for (let i in s){
          if(s[i]["live"]) {
            d3.select("#services-live").append("div").html(svcView(s[i]));
          } else {
            d3.select("#services-dead").append("div").html(svcView(s[i]));
          }
        }

        function dsHdr(){
          if(dead.length>0){
              deadSvcG.selectAll(".dsvcHdr").attr("opacity",1);
          } else {
            deadSvcG.selectAll(".dsvcHdr").attr("opacity",0.3);
          }
        }
        function sHdr(){
          if(lv.length>0){
            sndrG.selectAll(".svcHdr").attr("opacity",1);
            d3.select("#svcHdr").style("opacity",1);
          } else {
            sndrG.selectAll(".svcHdr").attr("opacity",0.3);
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
        // svcT.exit().attr("ger",function(d,i){console.log(i,d)});
        svcT.exit().html("");
        svcT = svcT.enter()
        .append("div").classed("svcT new",true)
        .merge(svcT)
          .html(function(d){return d["svc"].Name+"#"+d["svc"].Rank});

        // svg live services
        sndr = sndr.data(lv);
        sndr = sndr
          .classed("new",false)
          .classed("updated",true)
          .classed("live",function(d){return d["live"]}) // if not "live" set class "dead"
          .classed("manager",function(d){return d["svc"].Rank.indexOf(".M.")>0})
          .classed("dead",function(d){return !d["live"]});
        sndr.exit().remove();
        sndr = sndr.enter()
          .append("text")
          .classed("sndr new",true)
          .attr("x",0)
          .attr("dy",function(d,i){return i*1.2+"em"})
        .merge(sndr)
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

      // ------------------------------------------------------------------

      dispSvcs();
      dispBoards();

    });
  }

  // d3.interval loads data from the API endpoint
  // and calls the update function.
  update();
  d3.interval(function() {
    update();
  },UpdateInterval);
}
