/// <reference path="lib/dom.ts" />
/// <reference path="lib/types.ts" />
/// <reference path="lib/xhr.ts" />

namespace go {
    // Get the OS-specific shortcut key for copying.
    var copyKey = () => navigator.userAgent.indexOf('Macintosh') >= 0
        ? '⌘-C'
        : 'Ctrl-C';

    // Extract the name from the page location.
    var nameFrom = (uri: string) => {
        var parts = uri.substring(1).split('/');
        return parts[1];
    };

    // Called with the window resizes.
    var windowDidResize = () => {
        var rect = $frm.getBoundingClientRect();
        // Change the top margin of the form to put the middle of the form 
        // at the 1/3rd point in the window.
        dom.css($frm, 'margin-top', Math.max(50, (window.innerHeight/3 - rect.height/2)) + 'px');
    };

    // Called when the URL changes.
    var urlDidChange = () => {
        var url = ($url.value || '').trim();
        if (url == $route.url) {
            return;
        }

        $route.url = url;

        hideDrawer();
        if (url) {
            $cls.classList.add('vis');
        } else {
            $cls.classList.remove('vis');
        }

        postShortUrl(url, $shorturl.value);
    };

    var shortUrlDidChange = () => {
        const shorturl = ($shorturl.value || '').trim();
        if (shorturl == $route.name) {
            return;
        }

        $route.name = shorturl;

        $route.generated = false;

        postShortUrl($url.value || '', shorturl);
    };

    var formDidSubmit = (e: Event) => {
        e.preventDefault();

        var name = nameFrom(location.pathname),
            url = ($url.value || '').trim();

        postShortUrl(url, name);
    };

    var postShortUrl = (url: string, shorturl: string) => {
        $route.name = shorturl;
        $route.url = url;
        $route.modified_count += 1;

        xhr.post('/api/url/' + $route.name)
            .sendJSON($route)
            .onDone((data: string, status: number) => {
                var msg = <MsgRoute>JSON.parse(data);
                if (!msg.ok) {
                    showError(msg.error);
                    return;
                }
                console.log("Received response: ", msg);

                var route = msg.route;
                if (!route) {
                    // hideDrawer();
                    return;
                }

                if (route.modified_count >= $route.modified_count){
                    $route = route;

                    if ($route.url != $url.value && document.activeElement != $url){
                        $url.value = $route.url;
                    }

                    showLink($route.name)
                } else {
                    console.log("Route rejected, perhaps network delay?")
                }
            });
    }

    var formDidClear = () => {
        var name = nameFrom(location.pathname),
            url = ($url.value || '').trim();

        $url.value = '';
        urlDidChange();

        if (!name) {
            return;
        }

        xhr.create('DELETE', '/api/url/' + name)
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

    var showLink = (name: string) => {
        console.log("Showing link:", name)
        if (name){
            dom.css($links, 'transform', 'scaleY(1)');
        } else {
            dom.css($links, 'transform', 'scaleY(0)');
        }

        if (name != $shorturl.value && document.activeElement != $shorturl){
            $shorturl.value = name;
        }

        for (var echo of $echos) {
            console.log(echo);
            echo.innerText = name;
        }

        history.replaceState({}, null, '/edit/' + name);
        return;
    };

    // Called when the app loads initially.
    var appDidLoad = () => {
        windowDidResize();
        window.addEventListener('resize', windowDidResize, false);
        $frm.addEventListener('submit', formDidSubmit, false);

        $url.addEventListener('keyup', urlDidChange, false);
        $url.addEventListener('paste', urlDidChange, false);
        $url.addEventListener('change', urlDidChange, false);

        $shorturl.addEventListener('keyup', shortUrlDidChange, false);
        $shorturl.addEventListener('paste', shortUrlDidChange, false);
        $shorturl.addEventListener('change', shortUrlDidChange, false);

        $cls.addEventListener('click', formDidClear, false);

        var name = nameFrom(location.pathname);
        showLink(name);
        if (!name) {
            // We are making a new link
            $url.focus();
            $route = <Route>{
                name: "",
                url: "",
                generated: false,
                modified_count: 0
            };
            $route.uid = Math.floor(Math.random() * (1<<31)).toString();
            $inited = true;
            return;
        }

        $shorturl.value = name;

        xhr.get('/api/url/' + name)
            .send()
            .onDone((data: string, status: number) => {
                var msg = <MsgRoute>JSON.parse(data);

                if (status == 200) {
                    $route = msg.route;
                } else {
                    $route = <Route>{
                        name: name,
                        url: "",
                        generated: false,
                        modified_count: 0
                    };
                    $route.uid = Math.floor(Math.random() * (1<<31)).toString();
                }

                $url.value = $route.url;
                $url.focus();
                urlDidChange();
                $inited = true;
            });
    };

    var $frm = <HTMLFormElement>dom.q('form'),
        $cmp = dom.q('#cmp'),
        $links = dom.q('#links'),
        $cls = dom.q('#cls'),
        $url = <HTMLInputElement>dom.q('#url'),
        $shorturl = <HTMLInputElement>dom.q('#shorturl'),
        $echos = <Array<HTMLElement>>Array.prototype.slice.call(document.getElementsByClassName("echo")),
        $uid: string,
        // This object stores the latest route available, with the data type as defined on the server.
        // It is null if the app has not yet initialized.
        $route: Route,
        $inited: boolean = false;

    appDidLoad();
}