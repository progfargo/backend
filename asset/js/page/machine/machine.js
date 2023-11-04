$(document).ready(function () {
	$("select[name='catId']").change(function () {
		$("#categoryForm").submit();
	});

	$("select[name='manId']").change(function () {
		$("#manufacturerForm").submit();
	});
});