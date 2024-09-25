function ValidateInput(e) {
  const errors = [];
  if (e.target.value.length <= 0) {
    errors.push({ target: e.target, message: 'This field is required' });
    e.target.classList.add('border-red-500');
    e.target.insertAdjacentHTML(
      'afterend',
      '<p class="text-red-500 text-xs italic error">This field is required</p>'
    );
  } else {
    if (e.target.name === 'link') {
      const url = e.target.value;
      const regex = /^(http|https):\/\/[^ "]+$/;
      if (!regex.test(url)) {
        errors.push({ target: e.target, message: 'Invalid URL' });
        e.target.classList.add('border-red-500');
        e.target.insertAdjacentHTML(
          'afterend',
          '<p class="text-red-500 text-xs italic error">Invalid URL</p>'
        );
      }
    }

    if (e.target.name === 'tagify') {
      const tags = e.target.value.split(',');
      if (tags.length > 3) {
        errors.push({ target: e.target, message: 'Only 3 tags are allowed' });
        e.target.classList.add('border-red-500');
        e.target.insertAdjacentHTML(
          'afterend',
          '<p class="text-red-500 text-xs italic error">Only 3 tags are allowed</p>'
        );
      }
    }
  }

  if (errors.length <= 0) {
    e.target.classList.remove('border-red-500');

    const error = document.querySelector('.error');
    if (error) {
      error.remove();
    }
  }

  return errors;
}

function load() {
  var inputs = document.querySelectorAll('input');
  for (var i = 0; i < inputs.length; i++) {
    inputs[i].addEventListener('change', ValidateInput);
  }
  new Tagify(document.querySelector('input[name=tagify]'), {
    whitelist: [],
    maxTags: 5,
    originalInputValueFormat: (valuesArr) =>
      valuesArr.map((item) => item.value).join(','),
  });

  new Tagify(document.querySelector('input[name=disabled-user-input]'), {
    whitelist: ['foo', 'bar', 'baz'],
    userInput: false,
  });

  document
    .querySelector('input[name=tagify]')
    .addEventListener('change', function(e) {
      document.querySelector('input[name=tags]').value = e.target.tagifyValue;
    });

  document.querySelector('form').addEventListener('submit', function(e) {
    e.preventDefault();
    const errors = [];
    for (var i = 0; i < inputs.length; i++) {
      const input = inputs[i];
      errors.push(...ValidateInput({ target: input }));
    }
    if (errors.length > 0) {
      alert('Please fix the errors before submitting the posts');
    } else {
      document.querySelector('form').reset();
    }
  });
  const quill = new Quill('#editor', {
    theme: 'snow',
  });

  function update() {
    document.querySelector('textarea[name=content]').value = quill.root.innerHTML;
  }


  quill.on(Quill.events.TEXT_CHANGE, update);

}

document.addEventListener('DOMContentLoaded', function() {
  load();

});
