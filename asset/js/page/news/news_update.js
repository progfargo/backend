$(document).ready(function () {
	$lang = $("html").attr("lang");

	$("input[name='recordDate']").datetimepicker({
		timepicker: false,
		format: "d-m-Y",
		mask: false,
		lang: $lang
	});

	smallEditor("#summary");
	bigEditor("#body");
});