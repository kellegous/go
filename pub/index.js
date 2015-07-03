(function() {

var resize = function() {
  var rect = form.get(0).getBoundingClientRect();
  form.css('margin-top', window.innerHeight/3 - rect.height/2);
};

var nameFrom = function(uri) {
  var parts = uri.substring(1).split('/');
  return parts[1];
};

var load = function() {
  var name = nameFrom(location.pathname);
  $.ajax({
    url: '/api/url/' + name,
    dataType: 'json'
  }).always(function(data) {
    if (!data.ok) {
      // TODO(knorton): Error
      return;
    }

    var route = data.route,
        url = route.url || '';
    $('#url').val(url).focus();
  });
}

var showLink = function(name) {
  var cmp = $('#cmp'),
      lnk = location.origin + '/' + name;

  var a = $(document.createElement('a'))
    .attr('href', lnk)
    .text(lnk)
    .appendTo(cmp.text(''));

  cmp.css('transform', 'scaleY(1)');

  getSelection().setBaseAndExtent(a.get(0), 0, a.get(0), 1);
};

var form = $('form').on('submit', function(e) {
  e.preventDefault();
  var name = nameFrom(location.pathname),
      url = $('#url').val().trim();

  $.ajax({
    type: 'POST',
    url : '/api/url/' + name,
    data : JSON.stringify({ url : url }),
    dataType : 'json'
  }).always(function(data, txt, xhr) {
    if (!data.ok) {
      return;
    }

    var route = data.route;
    if (!route) {
      // deleted
    }

    var url = route.url || '',
        name = route.name || '';
    if (url) {
      history.replaceState({}, null, '/edit/' + name);
      showLink(name);
    }
  });
});

$('#cls').on('click', function(e) {
  $('#url').val('');
  $('form').submit();
});

window.addEventListener('resize', resize);
resize();
load();

})();
