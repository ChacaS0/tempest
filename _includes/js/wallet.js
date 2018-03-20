$(document).ready(function(){
	
	// Wallet presets
	$(".wall-ctnt").hide();
	// Wallet Title, showing / hiding content
	$(".wall-head").click(function(event){
		// hide others
		$(".wall-ctnt:visible").not($("#"+event.target.id).next()).toggle("slow");

		// show the right one
		$("#"+event.target.id).next().toggle("slow");

	});

});