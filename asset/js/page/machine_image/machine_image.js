$(document).ready(function () {
	$(".copyLink").on("click", function (event) {
		var str = $(this).prev("a").attr("href");

		var $textArea = $("<textarea class=\"copyClipboard\">" + str + "</textarea>");
		$textArea.insertAfter("body");
		$textArea.select();
		document.execCommand("copy");
		$(this).hide().fadeIn();
		$(".copyClipboard").remove();
	});

	$rotateLink = $(".rotateLink");

	$rotateLink.on("click", function (event) {
		var self = this;

		$.ajax({
			type: "GET",
			contentType: "application/json; charset=utf-8",
			url: $(this).attr("href"),
			data: {},
			dataType: "json",
			success: function (data) {
				if (data.status === "error") {
					showAjaxMsg(data.message);
					return;
				}

				showAjaxMsg(data.message);

				var url = $(self).attr("href");
				var id = GetUrlParameter(url, "imgId");
				var $image = $("img#" + id);
				var src = $image.attr("src");
				$image.removeAttr("src").attr("src", src);
			},

			error: function (result) {
				console.log(result.status, result.statusText);
			}
		});

		return false;
	});


	function GetUrlParameter(url, param) {
		var uRLVariables = url.split('&');

		for (var i = 0; i < uRLVariables.length; i++) {
			var parameterName = uRLVariables[i].split('=');

			if (parameterName[0] == param) {
				return parameterName[1];
			}
		}
	}
});