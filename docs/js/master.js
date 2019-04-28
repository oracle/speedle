function showMenuSel(isShow) {
	$('#langselFloatDiv').css('visibility', isShow ? 'visible' : 'hidden');
}

var timer;
$(function() {
	$('#top-nav-lang li').mouseenter(function(event) {
		showMenuSel(true);
		clearTimeout(timer);
	});

	$('#top-nav-lang li').mouseleave(function(event) {
		timer = setTimeout(function() {
			showMenuSel(false);
		}, 200);
	});

	$('#langselFloatDiv').mouseenter(function(event) {
		clearTimeout(timer);
	});

	$('#langselFloatDiv').mouseleave(function(event) {
		showMenuSel(false);
	});

	setMarginBottom();
	$(window).resize(function() {
		setMarginBottom();
	});
});

function setMarginBottom() {
	var height = $('footer').outerHeight(true);
	var smargin = $('.features').css('margin-top');
	var margin = smargin ? parseInt($('.features').css('margin-top')) : 0;
	$('body').css('margin-bottom', height + margin + 'px');
}
