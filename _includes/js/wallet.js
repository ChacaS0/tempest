$(document).ready(function(){
	
	// Wallet presets
	$(".wall-ctnt").hide();
	// Wallet Title, showing / hiding content
	$(".wall-head").click(function(event){
		var idEvt = "#"+event.target.id;

		// hide others
		$(".wall-ctnt:visible").not($(idEvt).next()).toggle("slow");

		// show the right one
		$(idEvt).next(".wall-ctnt").toggle("slow");
		
		// $(getid).load(getsrc);
		// // get the script URL
		
		// var iframe = document.getElementById(getid);
		// iframe.contentDocument.location=iframe.src;
		// $('#'+getid).attr('src', $('#'+getid).attr('src'));
		
		// $(getid).attr("src", "");
		// $(getid).attr("src", getsrc);
		
		// scroll to it
		$.wait(410).then(function(){
			$('html').animate({
				scrollTop: $(idEvt).offset().top
			}, 500);
	}
);
		
		// load frame
		// $(idEvt).next("script").load(this);

	});

});

// create the wait function...ish
$.wait = function(ms) {
	var defer = $.Deferred();
	setTimeout(function() { defer.resolve(); }, ms);
	return defer;
};