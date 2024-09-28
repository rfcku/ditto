function ValidateInput(e) {
  const errors = [];
  if (e.target.value.length <= 0) {
    errors.push({ target: e.target, message: 'This field is required' });
    e.target.classList.add('border-red-500');
    e.target.insertAdjacentHTML(
      'afterend',
      '<p class="text-red-500 text-xs italic error">This field is required</p>'
    );
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

  document.querySelector('form').addEventListener('submit', function (e) {
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
}

document.addEventListener('DOMContentLoaded', function () {
  load();
});
