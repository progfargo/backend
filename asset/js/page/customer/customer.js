$(document).ready(function () {
	$("select[name='sid']").change(function () {
		$("#sidForm").submit();
	});

	$("select[name='stat']").change(function () {
		$("#statForm").submit();
	});

	$("select[name='domain']").change(function () {
		$("#domainForm").submit();
	});
});
