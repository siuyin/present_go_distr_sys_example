var updateInterval=1000; //ms

var svg = d3.select("svg"),
  width = +svg.attr("width"),
  height = +svg.attr("height");



// update get data from the API endpoint
function update() {
  d3.json("/heartbeat", function(error, dat) {
    if (error) throw error;

    function svcs(){
      var now=new Date();
      var n=[];
      for (var key in dat) {
        s=dat[key];
        var live = now.getTime() - new Date(s.T) < 2000;
        if(s.Tx != null && s.Tx.length>0)tx=true;
        if(s.Rx != null && s.Rx.length>0)rx=true;
        n.push({"svc":s,"live":live,"tx":tx,"rx":rx});
      }
      return n;
    }

    function boards(){
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

    function svcView(s){
      return s["svc"].Name + " "+ s["svc"].T+":"+s["live"]+":"+s["tx"]+":"+s["rx"];
    }
    function dispSvcs(){
      d3.select("#services-live").selectAll("div").remove();
      d3.select("#services-dead").selectAll("div").remove();
      s=svcs();
      for (var i in s){
        if(s[i]["live"]) {
          d3.select("#services-live").append("div").html(svcView(s[i]));
        } else {
          d3.select("#services-dead").append("div").html(svcView(s[i]));
        }
      }
    }

    function boardView(b){
      return b;
    }
    function dispBoards(){
      d3.select("#boards").selectAll("div").remove();
      b=boards();
      for (var i in b){
        d3.select("#boards").append("div").html(boardView(b[i]));
      }
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
},updateInterval);
