(function () {
	var $frm = $('form'),
		$cmp = $('#cmp'),
		$cls = $('#cls'),
		$url = $('#url'),
		copyKey = navigator.userAgent.indexOf('Macintosh') >= 0 ? '⌘-C' : 'Ctrl-C',
		lastUrl;

	var $e = function (name) {
		return $(document.createElement(name));
	};

	var resize = function () {
		var rect = $frm.get(0).getBoundingClientRect();
		$frm.css('margin-top', window.innerHeight / 3 - rect.height / 2);
	};

	var nameFrom = function (uri) {
		var parts = uri.substring(1).split('/');
		return parts[1];
	};

	var load = function () {
		var name = nameFrom(location.pathname);
		$.ajax({
			url: '/api/url/' + name,
			dataType: 'json',
		}).always(function (data) {
			if (!data.ok) {
				showError(data.error);
				return;
			}

			var route = data.route,
				url = (route && route.url) || '';
			$url.val(url).focus();
			urlDidChange();
		});
	};

	var showLink = function (name, src) {
        var lnk = '/' + name;

        if (src != '') {
            lnk = src + lnk;
        } else {
            lnk = location.origin + lnk;
        }

		$cmp.text('').removeClass('fuck').addClass('link');

		var a = $(document.createElement('a'))
			.attr('href', lnk)
			.text(lnk)
			.appendTo($cmp);

		var h = $(document.createElement('span'))
			.addClass('hnt')
			.text(copyKey)
			.appendTo($cmp);

		$cmp.css('transform', 'scaleY(1)');

		getSelection().setBaseAndExtent(a.get(0), 0, a.get(0), 1);
	};

	var showError = function (message) {
		$cmp.text('').removeClass('link').addClass('fuck');

		$e('span')
			.text('ERROR: ' + message)
			.appendTo($cmp);

		$cmp.css('transform', 'scaleY(1)');
	};

	var hideDrawer = function () {
		$cmp.css('transform', 'scaleY(0)');
	};

	var urlDidChange = function () {
		var url = $url.val().trim();
		if (url == lastUrl) {
			return;
		}

		lastUrl = url;

		hideDrawer();
		if (url) {
			$cls.fadeIn(200);
		} else {
			$cls.fadeOut(200);
		}
	};

	$frm.on('submit', function (e) {
		e.preventDefault();
		var name = nameFrom(location.pathname),
			url = $url.val().trim();

		$.ajax({
			type: 'POST',
			url: '/api/url/' + name,
			data: JSON.stringify({ url: url }),
			dataType: 'json',
		}).always(function (data) {
			if (!data.ok) {
				showError(data.error);
				return;
			}

			var route = data.route;
			if (!route) {
				hideDrawer();
				return;
			}

			var url = route.url || '',
				name = route.name || '',
			    src = route.source_host || '';
			if (url) {
				history.replaceState({}, null, '/edit/' + name);
				showLink(name, src);
			}
		});
	});

	$url
		.on('keydown', urlDidChange)
		.on('paste', urlDidChange)
		.on('change', urlDidChange);

	$cls.on('click', function (e) {
		$url.val('');
		$frm.submit();
		urlDidChange();
	});

	window.addEventListener('resize', resize);
	resize();
	urlDidChange();
	load();
})();
