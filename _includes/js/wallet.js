$(document).ready(function(){
	
	// Wallet presets
	$(".wall-ctnt").hide();
	// Wallet Title, showing / hiding content
	$(".wall-head").click(function(event){
		// hide others
		$(".wall-ctnt:visible").not($("#"+event.target.id).next()).toggle("slow");

		// show the right one
		$("#"+event.target.id).next().toggle("slow");
		
		// scroll to it
		$.wait(1200).then(
			$('html, body').animate({
				scrollTop: $("#"+event.target.id).offset().top
			}, 1000)
		);
	});

});

// create the wait function...ish
$.wait = function(ms) {
	var defer = $.Deferred();
	setTimeout(function() { defer.resolve(); }, ms);
	return defer;
};