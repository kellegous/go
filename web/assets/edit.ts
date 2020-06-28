/// <reference path="lib/dom.ts" />
/// <reference path="lib/types.ts" />
/// <reference path="lib/xhr.ts" />

namespace go {
	// Get the OS-specific shortcut key for copying.
	var copyKey = () =>
		navigator.userAgent.indexOf('Macintosh') >= 0 ? '⌘C (copy)' : 'Ctrl+C';

	// Extract the name from the page location.
	var nameFrom = (uri: string) => {
		var parts = uri.substring(1).split('/');
		return parts[1];
	};

	// Called with the window resizes.
	var windowDidResize = () => {
		var rect = $frm.getBoundingClientRect();
		dom.css(
			$frm,
			'margin-top',
			window.innerHeight / 3 - rect.height / 2 + 'px'
		);
	};

	// Called when the URL changes.
	var urlDidChange = () => {
		var url = ($url.value || '').trim();
		if (url == lastUrl) {
			return;
		}

		lastUrl = url;

		hideDrawer();
		if (url) {
			$cls.classList.add('vis');
		} else {
			$cls.classList.remove('vis');
		}
	};

	var formDidSubmit = (e: Event) => {
		e.preventDefault();

		var name = nameFrom(location.pathname),
			url = ($url.value || '').trim();

		xhr
			.post('/api/url/' + name)
			.sendJSON({ url: url })
			.onDone((data: string, status: number) => {
				var msg = <MsgRoute>JSON.parse(data);
				if (!msg.ok) {
					showError(msg.error);
					return;
				}

				var route = msg.route;
				if (!route) {
					hideDrawer();
					return;
				}

				var url = route.url || '',
					name = route.name || '',
					host = route.source_host || '';
				if (url) {
					history.replaceState({}, null, '/edit/' + name);
					showLink(name, host);
				}
			});
	};

	var formDidClear = () => {
		var name = nameFrom(location.pathname),
			url = ($url.value || '').trim();

		$url.value = '';
		urlDidChange();

		if (!name) {
			return;
		}

		xhr
			.create('DELETE', '/api/url/' + name)
			.send()
			.onDone((data: string, status: number) => {
				var msg = <Msg>JSON.parse(data);
				if (!msg.ok) {
					showError(msg.error);
				}
			});
	};

	var hideDrawer = () => {
		dom.css($cmp, 'transform', 'scaleY(0)');
	};

	var showError = (msg: string) => {
		$cmp.textContent = '';
		$cmp.classList.remove('link');
		$cmp.classList.add('fuck');

		var $s = dom.c('span');
		$s.textContent = 'ERROR: ' + msg;
		$cmp.appendChild($s);

		dom.css($cmp, 'transform', 'scaleY(1)');
	};

	// This function shows the keyword link in quick-copy dropdown
	var showLink = (name: string, src: string) => {
		var lnk = '/' + name;

		if (src != '') {
			lnk = src + lnk;
		} else {
			lnk = location.origin + lnk;
		}

		// Create a node text element and add class="link"
		$cmp.textContent = '';
		$cmp.classList.remove('fuck');
		$cmp.classList.add('link');

		// Create an anchor link element
		var $a = dom.c('a');
		$a.setAttribute('href', lnk);
		$a.textContent = lnk;
		$cmp.appendChild($a);

		// Add copy hint to	quick-copy drop down
		var $h = dom.c('span');
		$h.classList.add('hnt');
		$h.textContent = copyKey();
		$cmp.appendChild($h);

		// Open the quick-copy drawer
		dom.css($cmp, 'transform', 'scaleY(1)');

		// Select the text in the dropdown
		getSelection().setBaseAndExtent($a, 0, $a, 1);
	};

	// Called when the app loads initially.
	var appDidLoad = () => {
		windowDidResize();
		window.addEventListener('resize', windowDidResize, false);
		$frm.addEventListener('submit', formDidSubmit, false);

		$url.addEventListener('keyup', urlDidChange, false);
		$url.addEventListener('paste', urlDidChange, false);
		$url.addEventListener('change', urlDidChange, false);

		$cls.addEventListener('click', formDidClear, false);

		var name = nameFrom(location.pathname);
		if (!name) {
			$url.focus();
			return;
		}

		xhr
			.get('/api/url/' + name)
			.send()
			.onDone((data: string, status: number) => {
				var msg = <MsgRoute>JSON.parse(data);

				if (status != 200) {
					return;
				}

				// TODO(knorton): Hanlde things.
				var url = msg.route.url || '';
				$url.value = url;
				$url.focus();
				urlDidChange();
			});
	};

	var $frm = <HTMLFormElement>dom.q('form'),
		$cmp = dom.q('#cmp'),
		$cls = dom.q('#cls'),
		$url = <HTMLInputElement>dom.q('#url'),
		lastUrl: string,
		$key = nameFrom(location.pathname);

	appDidLoad();
}
