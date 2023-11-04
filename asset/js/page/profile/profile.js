$(document).ready(function() {
	var $deleteAvatar = $(".deleteAvatar");
	var $confirm = $(".deleteAvatarConfirm");
		
	$(".deleteAvatar").click(function(event) {
		$deleteAvatar.hide();
		$confirm.toggleClass("hidden");
		
		var self = this;
		timeOut = setTimeout(function() {
			$confirm.toggleClass("hidden");
			$deleteAvatar.show();
		}, 4000);
	});
});