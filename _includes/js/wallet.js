$(document).ready(function(){
	
	// Wallet presets
	$(".wall-ctnt").hide();
	// Wallet Title, showing / hiding content
	$(".wall-head").click(function(event){
		var idEvt = "#"+event.target.id;

		// hide others
		$(".wall-ctnt:visible").not($(idEvt).next()).toggle("slow");

		// show the right one & scroll to it when it is fully visible
		$(idEvt).next(".wall-ctnt").toggle("slow").promise().done(function(){
			$('html, body').animate({ scrollTop: $(idEvt).offset().top }, 500);
		});

	});

});

// @DEPRECATED 
// // create the wait function...ish
// $.wait = function(ms) {
// 	var defer = $.Deferred();
// 	setTimeout(function() { defer.resolve(); }, ms);
// 	return defer;
// };