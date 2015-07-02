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
  }).success(function(data) {
    var url = data.url || '';
    $('#url').val(url);
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

  if (!url) {
    return;
  }

  $.ajax({
    type: 'POST',
    url : '/api/url/' + name,
    data : JSON.stringify({ url : url }),
    dataType : 'json'
  }).success(function(data) {
    var url = data.url || '';
    if (url) {
      showLink(name);
    }
  });
});

window.addEventListener('resize', resize);
resize();
load();

})();