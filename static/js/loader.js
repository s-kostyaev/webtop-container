$(document).ready(function(){

    updatePage = function(data) {
        jPut.text.data = [data];
        jPut.table.data = data.Procs;
        $("#content").removeClass("invisible");
    };

    requestUrl = function(url) {
        var answer = $.getJSON(url,function(data) {
        updatePage(data);});
    }; 
        
    requestUrl("/command/top");
})
