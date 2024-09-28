function ValidateInput(e) {
  const errors = [];
  console.log("ValidatingInput", e.target.type, e.target.value);
  if (e.target.type === 'file') {
    return errors;
  }

  if (e.target.type === 'radio') {
    return errors;
  }

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

  document.getElementById('type').addEventListener('change', function(e) {
    const type = e.target.value;
    if (type === 'link') {
      document.getElementById('link').classList.remove('hidden');
      document.getElementById('files').classList.add('hidden');
      document.getElementById('content').classList.add('hidden');
    } else {
      document.getElementById('link').classList.add('hidden');
      document.getElementById('files').classList.remove('hidden');
    }

    if (type === 'text') {
      document.getElementById('files').classList.add('hidden');
      document.getElementById('link').classList.add('hidden');
      document.getElementById('content').classList.remove('hidden');
    }

  });

  var inputs = document.querySelectorAll('input');
  for (var i = 0; i < inputs.length; i++) {
    inputs[i].addEventListener('change', ValidateInput);
  }

  document.querySelector('form').addEventListener('submit', function(e) {
    e.preventDefault();
    const errors = [];
    for (var i = 0; i < inputs.length; i++) {
      const input = inputs[i];
      errors.push(...ValidateInput({ target: input }));
    }
    if (errors.length > 0) {
      console.error('Please fix the errors before submitting the posts', errors);
    } else {
      document.querySelector('form').reset();
    }
  });
  
  document.querySelector('input[name=files]').addEventListener('change', function(e) {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = function(e) {
      document.querySelector('img').src = e.target.result;
    };
    reader.readAsDataURL(file);
  });

  document.getElementById('search').addEventListener('change', function(e) {
    console.log("Search", e); 
    console.log("Search", e.target.value);

    if (e.target.value.length <= 0) {
      document.getElementById('search-results').innerHTML = '';
    }
    // if looses  focus
    if (e.target.focused === false) {
      document.getElementById('search-results').innerHTML = '';
    }

  });

  function update() {
    document.querySelector('textarea[name=content]').value = quill.root.innerHTML;
  }

}

document.addEventListener('DOMContentLoaded', function() {
  load();
});
