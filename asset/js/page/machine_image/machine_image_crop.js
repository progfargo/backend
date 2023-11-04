$(document).ready(function () {
	var $image = $("#image");

	var $dimention = $("p.imgInfo span.dimention");
	var originalW, originalH, templateW, templateH;
	var resultW, resultH, resultX, resultY;

	var cropInner, cropHandler;
	var $cropButton = $("button:submit");

	originalW = $image.data("width");
	originalH = $image.data("height");

	if (originalW > originalH) {
		$image.rcrop({
			minSize: [160, 100],
			preserveAspectRatio: true,
			grid: true,
		});

		$image.on("rcrop-changed", function () {
			var dim = $(this).rcrop("getValues");

			resultW = Math.ceil(dim.width / templateW * originalW);
			resultH = Math.ceil(dim.height / templateH * originalH);
			resultX = Math.ceil(dim.x / templateW * originalW);
			resultY = Math.ceil(dim.y / templateH * originalH);

			resultX = (resultX < 0) ? 0 : resultX;
			resultY = (resultY < 0) ? 0 : resultY;
			resultW = (resultW > originalW ? originalW : resultW);
			resultH = (resultH > originalH ? originalH : resultH);

			var str = resultW + "px * " + resultH + "px";
			$dimention.html(str);

			if (resultW < 1280 || resultH < 800) {
				$dimention.addClass("invalidDim");
				$cropInner.addClass("invalidDimBorder");
				$cropHandler.addClass("invalidDimBackground");
				$cropButton.attr("disabled", "disabled");
				$cropButton.removeClass("buttonPrimary");
				$cropButton.addClass("buttonDefault");
			} else {
				$dimention.removeClass("invalidDim");
				$cropInner.removeClass("invalidDimBorder");
				$cropHandler.removeClass("invalidDimBackground");
				$cropButton.removeAttr("disabled");
				$cropButton.removeClass("buttonDefault");
				$cropButton.addClass("buttonPrimary");
			}
		});
	} else {
		$image.rcrop({
			minSize: [100, 160],
			preserveAspectRatio: true,
			grid: true,
		});

		$image.on("rcrop-changed", function () {
			var dim = $(this).rcrop("getValues");

			resultW = Math.ceil(dim.width / templateW * originalW);
			resultH = Math.ceil(dim.height / templateH * originalH);
			resultX = Math.ceil(dim.x / templateW * originalW);
			resultY = Math.ceil(dim.y / templateH * originalH);

			resultX = (resultX < 0) ? 0 : resultX;
			resultY = (resultY < 0) ? 0 : resultY;
			resultW = (resultW > originalW ? originalW : resultW);
			resultH = (resultH > originalH ? originalH : resultH);

			var str = resultW + "px * " + resultH + "px";
			$dimention.html(str);

			if (resultW < 800 || resultH < 1280) {
				$dimention.addClass("invalidDim");
				$cropInner.addClass("invalidDimBorder");
				$cropHandler.addClass("invalidDimBackground");
				$cropButton.attr("disabled", "disabled");
				$cropButton.removeClass("buttonPrimary");
				$cropButton.addClass("buttonDefault");
			} else {
				$dimention.removeClass("invalidDim");
				$cropInner.removeClass("invalidDimBorder");
				$cropHandler.removeClass("invalidDimBackground");
				$cropButton.removeAttr("disabled");
				$cropButton.removeClass("buttonDefault");
				$cropButton.addClass("buttonPrimary");
			}
		});
	}

	$image.on("rcrop-ready", function () {
		templateW = $image.width();
		templateH = $image.height();

		$cropInner = $(".rcrop-wrapper .rcrop-croparea .rcrop-croparea-inner");
		$cropHandler = $(".clayfy-handler.clayfy-touch-device");

		$(this).trigger("rcrop-changed");
	});

	$("form#cropForm").on("submit", function () {
		$("input[name='x[]']").val(resultX);
		$("input[name='y[]']").val(resultY);
		$("input[name='width[]']").val(resultW);
		$("input[name='height[]']").val(resultH);
	});
});