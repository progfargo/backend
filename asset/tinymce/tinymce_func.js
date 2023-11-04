var smallEditor, bigEditor;

$(document).ready(function () {
	var lang = $("html").attr("lang");

	smallEditor = function (id, opt) {
		var height = (opt && opt.height) || "200px";

		tinymce.init({
			selector: id,
			entity_encoding: "raw",
			theme: "modern",
			language: lang,
			content_css: "/asset/tinymce/content.css",
			height: height,
			plugins: [
				"advlist autolink lists link charmap print preview anchor",
				"searchreplace visualblocks code fullscreen",
				"insertdatetime table contextmenu paste"
			],
			toolbar: "undo redo | styleselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent",
			setup: function (editor) {
				editor.on('init', function () {
					var editorId = editor.id;
					var textarea = $("#" + editorId);
					var tabindex = textarea.attr("tabindex");
					$("#" + editorId + "_ifr").attr("tabindex", textarea.attr("tabindex"));
					textarea.attr("tabindex", null);
				});
			},
		});
	};

	bigEditor = function (id, opt) {
		var height = (opt && opt.height) || "400px";

		tinymce.init({
			selector: id,
			entity_encoding: "raw",
			theme: "modern",
			language: lang,
			content_css: "/asset/tinymce/content.css",
			height: height,
			plugins: [
				"advlist autolink lists link image charmap print preview hr anchor pagebreak",
				"searchreplace wordcount visualblocks visualchars code fullscreen",
				"insertdatetime media nonbreaking save table contextmenu directionality",
				"emoticons template paste textcolor colorpicker textpattern"
			],
			toolbar1: "insertfile undo redo | styleselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | link image",
			image_advtab: true,
			setup: function (editor) {
				editor.on('init', function () {
					var editorId = editor.id;
					var textarea = $("#" + editorId);
					var tabindex = textarea.attr("tabindex");
					$("#" + editorId + "_ifr").attr("tabindex", textarea.attr("tabindex"));
					textarea.attr("tabindex", null);
				});
			},

			image_class_list: [
				{ title: 'Left', value: 'floatLeft' },
				{ title: 'Right', value: 'floatRight' }
			],
		});
	}
});
