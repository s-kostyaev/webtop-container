$(document).ready(function(){

    updatePage = function(data) {
    jPut.text.data = data;
    var table = data.Procs;
    jPut.table.data = table;
    $("#content").removeClass("invisible");
    };

    requestUrl = function(url) {
    var answer = $.getJSON(url);
    updatePage(answer);
    }; 
        
    requestUrl("/command/top");
})
