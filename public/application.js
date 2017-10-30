var DEFAULT_POLL_TIMEOUT = 15;

function poll(id, cb, cbTimeout, timeout) {
  $.get("/requests/"+id).done(function(data) {
    if (timeout === 0) {
      return cbTimeout();
    }

    if (data.payload) {
      return cb(data);
    }

    setTimeout(function() {
      poll(id, cb, cbTimeout, (timeout === undefined ? DEFAULT_POLL_TIMEOUT : timeout) - 1);
    }, 1000);
  });
}

function pollRequest(id) {
  poll(id, function(data) {
    for (var field in data.payload) {
      if (data.payload.hasOwnProperty(field)) {
        $('#' + data.type + "_" + field)
          .closest('.mdl-textfield').get(0)
          .MaterialTextfield.change(data.payload[field]);
      }
    }

    $('form input').attr('disabled', false);
    $('#progress').addClass('is-hidden');
    $('#result').removeClass('is-hidden');
    $('main').scrollTo('#result', 'fast');
  }, function() {
    $('form input').attr('disabled', false);
    $('#progress').addClass('is-hidden');
    $('#timeout').removeClass('is-hidden');
  });
}

function prepareForm(fn) {
  $('#form input[type="hidden"]').val(fn.id);
  $('#form .mdl-card__title-text span').html(fn.name);
  $('#inputs [rel="field"], #outputs [rel="field"]').remove();
  $.each(fn.inputs, function(i, input) {
    var template = $('#inputs > div:first-child').clone();
    template.attr('rel', 'field').find('input').attr({
      id: fn.id + "_" + input.id,
      name: "payload[" + input.id + "]",
      type: input.type,
      placeholder: input.hint,
      required: true,
      // autofocus: i == 0,
    });
    template.find('label').attr('for', fn.id + "_" + input.id).html(input.label);
    if (i > 0) { $('#inputs').append($('<br />').attr('rel', 'field')); }
    $('#inputs').append(template.removeClass('is-hidden'));
    template.get(0).MaterialTextfield = new MaterialTextfield(template.get(0));
  });
  $.each(fn.outputs, function(i, output) {
    var template = $('#outputs > div:first-child').clone();
    template.attr('rel', 'field').find('input').attr({
      id: fn.id + "_" + output.id,
      type: output.type,
    });
    template.find('label').attr('for', fn.id + "_" + output.id).html(output.label);
    if (i > 0) { $('#outputs').append($('<br />').attr('rel', 'field')); }
    $('#outputs').append(template.removeClass('is-hidden'));
    template.get(0).MaterialTextfield = new MaterialTextfield(template.get(0));
  });
  $('#result').addClass('is-hidden');

  $('#fns').slideUp();
  $('#form').data('fn', fn.id).slideDown();
}

$(document).on('click', '#back', function(e) {
  e.preventDefault();

  $('#fns').slideDown();
  $('#form').slideUp(function() {
    $('main').scrollTo('#' + $(this).data('fn'), 'fast');
  });
});

$(document).on('submit', 'form', function(e) {
  e.preventDefault();

  var formData = $(this).serializeJSON();

  $('form input').attr('disabled', true);
  $('#progress').removeClass('is-hidden');
  $('#timeout').addClass('is-hidden');
  $('#result').addClass('is-hidden');

  $.post('/requests', JSON.stringify(formData)).done(function(data) {
    pollRequest(data.id);
  });
});

$(function() {
  $.get('/discover').done(function(data) {
    poll(data.id, function(data) {
      $.each(data.payload, function(i, fn) {
        var template = $('#fns > .mdl-cell:first-child').clone();
        template.find('.mdl-card__title-text').html(fn.name)
        template.find('.mdl-card__supporting-text').html(fn.description);
        template.find('a').click(function(e) {
          e.preventDefault();
          prepareForm(fn);
        });
        $('#fns').append(template.attr('id', fn.id).removeClass('is-hidden'));
      });

      $('#loading').addClass('is-hidden');
      $('#content').removeClass('is-hidden');
    }, function() {
      $('#loading').addClass('is-hidden');
      $('#loading-failed').removeClass('is-hidden');
    });
  });
});
