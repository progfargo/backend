$(document).ready(function() {
	var $urlList = $(".backgroundIdList span");

	var urlList = [];
	$urlList.each(function(item) {
		urlList.push($(this).text());
	});

	$loginPage = $("body#loginPage");
	var item = 0;
	$loginPage.css("background-image", "url(" + urlList[item] + ")")

	setInterval(function() {
		item++;
		if (item > urlList.length - 1) {
			item = 0;
		}

		$loginPage.css("background-image", "url(" + urlList[item] + ")")
	}, 10000);
});
